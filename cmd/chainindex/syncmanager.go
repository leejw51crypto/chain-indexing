package main

import (
	"fmt"
	"github.com/crypto-com/chainindex/appinterface/rdb"
	"github.com/crypto-com/chainindex/appinterface/rdbstatusstore"
	"time"

	chainfeed "github.com/crypto-com/chainindex/infrastructure/feed/chain"
	"github.com/crypto-com/chainindex/infrastructure/notification"
	"github.com/crypto-com/chainindex/infrastructure/tendermint"
	applogger "github.com/crypto-com/chainindex/internal/logger"
)

const DEFAULT_POLLING_INTERVAL = 5 * time.Second

type SyncManager struct {
	rdbConn rdb.Conn
	client  *tendermint.HTTPClient
	logger  applogger.Logger
	Subject *chainfeed.BlockSubject

	pollingInterval time.Duration
}

// NewSyncManager creates a new feed with polling for latest block starts at a specific height
func NewSyncManager(logger applogger.Logger, tendermintRPCUrl string, rdbConn rdb.Conn) *SyncManager {
	tendermintClient := tendermint.NewHTTPClient(tendermintRPCUrl)

	return &SyncManager{
		rdbConn: rdbConn,
		client:  tendermintClient,
		logger:  logger,

		pollingInterval: DEFAULT_POLLING_INTERVAL,
	}
}

func (manager *SyncManager) UpdateIndexedHeight(nextHeight int64, handle *rdb.Handle) error {
	statusStore := rdbstatusstore.NewRDbStatusStoreImpl(handle)
	err := statusStore.UpdateLastIndexedBlockHeight(nextHeight)
	if err != nil {
		return fmt.Errorf("error running UpdateLastIndexedBlockHeight %v", err)
	}
	return nil
}

// SyncBlocks makes request to tendermint, create and dispatch notifications
func (manager *SyncManager) SyncBlocks(latestHeight int64) error {
	statusStore := rdbstatusstore.NewRDbStatusStoreImpl(manager.rdbConn.ToHandle())
	lastIndexedHeight, err := statusStore.GetLastIndexedBlockHeight()
	if err != nil {
		return fmt.Errorf("error running GetLastIndexedBlockHeight %v", err)
	}

	// Sync next height to avoid duplication
	currentIndexingHeight := lastIndexedHeight + 1

	for currentIndexingHeight < latestHeight {
		// Request tendermint RPC
		block, rawBlock, err := manager.client.Block(currentIndexingHeight)
		if err != nil {
			return fmt.Errorf("error getting chain's block at %d: %v", currentIndexingHeight, err)
		}

		blockResults, err := manager.client.BlockResults(currentIndexingHeight)
		if err != nil {
			return fmt.Errorf("error getting chain's block_results at %d: %v", currentIndexingHeight, err)
		}

		// Create new block notification and notify subscribers
		notif := notification.NewBlockNotification(
			currentIndexingHeight, block, rawBlock, blockResults,
		)
		manager.Subject.Notify(notif, manager.rdbConn.ToHandle())

		// Current block indexing done, update db and sync next height
		manager.logger.Infof("block height %d synced and events produced", block.Height)
		err = manager.UpdateIndexedHeight(currentIndexingHeight, manager.rdbConn.ToHandle())
		if err != nil {
			return fmt.Errorf("error updating last indexed height for height %d: %v", currentIndexingHeight, err)
		}

		// If there is any error before, short-circuit return in the error handling
		// while the local currentIndexingHeight won't be incremented and will be retried later
		currentIndexingHeight += 1
	}
	return nil
}

// InitSubject creates subject and attach subscribers
func (manager *SyncManager) InitSubject() *chainfeed.BlockSubject {
	// Currently only the chain processor subscriber
	// add more subscriber base on the need
	chainProcessor := chainfeed.NewBlockSubscriber(0)

	blockSubject := chainfeed.NewBlockSubject()
	blockSubject.Attach(chainProcessor)

	return blockSubject
}

// Run starts the polling service for blocks
// new BlockFeedSubject and add listeners
func (manager *SyncManager) Run() error {
	manager.Subject = manager.InitSubject()

	tracker := chainfeed.NewBlockHeightTracker(manager.logger, manager.client)
	for {
		latestHeight := tracker.GetLatestBlockHeight()
		if latestHeight == nil {
			<-time.After(manager.pollingInterval)
			continue
		}
		if err := manager.SyncBlocks(*latestHeight); err != nil {
			manager.logger.Errorf("Error synchronizing blocks: %v", err)
		}

		<-time.After(manager.pollingInterval)
	}
}
