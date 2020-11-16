package chain

import (
	"fmt"
	"sync"

	event_interface "github.com/crypto-com/chainindex/appinterface/event"
	"github.com/crypto-com/chainindex/infrastructure/notification"
)

type BlockSubject struct {
	Observers sync.Map
}

func NewBlockSubject() *BlockSubject {
	return &BlockSubject{
		Observers: sync.Map{},
	}
}

func (s *BlockSubject) Attach(subj *BlockSubscriber) {
	s.Observers.Store(subj, struct{}{})
}

func (s *BlockSubject) Notify(n *notification.BlockNotification, eventStore *event_interface.RDbStore) {
	s.Observers.Range(func(key interface{}, value interface{}) bool {
		if key == nil || value == nil {
			return false
		}

		if err := key.(*BlockSubscriber).OnNotification(n, eventStore); err != nil {
			fmt.Println("error when subscriber run callback function", err)
		}
		return true
	})
}
