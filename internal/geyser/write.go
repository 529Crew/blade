package geyser

import (
	"sync"

	pb "github.com/529Crew/blade/internal/geyser/proto"
)

var currentSubReq = &pb.SubscribeRequest{
	Accounts: make(map[string]*pb.SubscribeRequestFilterAccounts),
	// Slots: make(map[string]*pb.SubscribeRequestFilterSlots),
	Transactions: make(map[string]*pb.SubscribeRequestFilterTransactions),
	// TransactionsStatus: make(map[string]*pb.SubscribeRequestFilterTransactions),
	// Blocks: make(map[string]*pb.SubscribeRequestFilterBlocks),
	// BlocksMeta: make(map[string]*pb.SubscribeRequestFilterBlocksMeta),
	// Entry: make(map[string]*pb.SubscribeRequestFilterEntry),
}
var subReqMutex = &sync.Mutex{}

/* ACCOUNTS */
func SubscribeAccounts(label string, filter *pb.SubscribeRequestFilterAccounts) (string, chan *pb.SubscribeUpdate, error) {
	subReqMutex.Lock()
	currentSubReq.Accounts[label] = filter
	subReqMutex.Unlock()

	err := geyserStream.Send(currentSubReq)
	if err != nil {
		return "", nil, err
	}

	listenerId, listener := Listen(label)

	return listenerId, listener, nil
}

func UnsubscribeAccounts(label string, listenerId string) error {
	subReqMutex.Lock()
	delete(currentSubReq.Accounts, label)
	subReqMutex.Unlock()

	err := geyserStream.Send(currentSubReq)
	if err != nil {
		return err
	}

	Unlisten(label, listenerId)

	return nil
}

/* TRANSACTIONS */
func SubscribeTransactions(label string, filter *pb.SubscribeRequestFilterTransactions) (string, chan *pb.SubscribeUpdate, error) {
	subReqMutex.Lock()
	currentSubReq.Transactions[label] = filter
	subReqMutex.Unlock()

	err := geyserStream.Send(currentSubReq)
	if err != nil {
		return "", nil, err
	}

	listenerId, listener := Listen(label)

	return listenerId, listener, nil
}

func UnsubscribeTransactions(label string, listenerId string) error {
	subReqMutex.Lock()
	delete(currentSubReq.Transactions, label)
	subReqMutex.Unlock()

	err := geyserStream.Send(currentSubReq)
	if err != nil {
		return err
	}

	Unlisten(label, listenerId)

	return nil
}
