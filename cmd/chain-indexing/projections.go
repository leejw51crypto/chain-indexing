package main

import (
	"github.com/crypto-com/chain-indexing/appinterface/projection/account_message"
	"github.com/crypto-com/chain-indexing/appinterface/projection/block"
	"github.com/crypto-com/chain-indexing/appinterface/projection/blockevent"
	transaction "github.com/crypto-com/chain-indexing/appinterface/projection/transaction"
	"github.com/crypto-com/chain-indexing/appinterface/projection/validator"

	// crossfire "github.com/crypto-com/chain-indexing/appinterface/projection/crossfire"

	"github.com/crypto-com/chain-indexing/appinterface/projection/validatorstats"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	projection_entity "github.com/crypto-com/chain-indexing/entity/projection"
	applogger "github.com/crypto-com/chain-indexing/internal/logger"
)

func initProjections(
	logger applogger.Logger,
	rdbConn rdb.Conn,
	consNodeAddressPrefix string,
) []projection_entity.Projection {
	return []projection_entity.Projection{
		block.NewBlock(logger, rdbConn),
		transaction.NewTransaction(logger, rdbConn),
		blockevent.NewBlockEvent(logger, rdbConn),
		validator.NewValidator(
			logger, rdbConn, consNodeAddressPrefix,
		),
		validatorstats.NewValidatorStats(logger, rdbConn),
		account_message.NewAccountMessage(logger, rdbConn),

		// TODO: crossfire projection to be registered here

		// register more projections here
	}
}
