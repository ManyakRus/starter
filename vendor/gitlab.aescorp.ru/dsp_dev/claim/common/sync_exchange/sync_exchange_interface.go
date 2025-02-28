package sync_exchange

import (
	"sync"

	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_types"
)

type ISyncExchange interface {
	InitSyncExchange(url string, serviceName string, version string) error
	DeInitSyncExchange() error
	SendMessage(topic string, pack sync_types.SyncPackage) error
	WaitMessage(topic string, callback Callback) error
	Subscribe(topic string, callback Callback) error
	SendRequest(receiver string, pack sync_types.SyncPackage, timeout int) (result sync_types.SyncPackage, err error)
	SendResponse(packIn *sync_types.SyncPackage, packOut sync_types.SyncPackage) error
}

type SSyncExhange struct{}

var New = sync.OnceValue(func() ISyncExchange {
	var h = SSyncExhange{}

	return &h
})

func (s *SSyncExhange) InitSyncExchange(url string, serviceName string, version string) error {
	return InitSyncExchange(url, serviceName, version)
}

func (s *SSyncExhange) DeInitSyncExchange() error {
	return DeInitSyncExchange()
}

func (s *SSyncExhange) SendMessage(topic string, pack sync_types.SyncPackage) error {
	return SendMessage(topic, pack)
}

func (s *SSyncExhange) WaitMessage(topic string, callback Callback) error {
	return WaitMessage(topic, callback)
}

func (s *SSyncExhange) Subscribe(topic string, callback Callback) error {
	return Subscribe(topic, callback)
}

func (s *SSyncExhange) SendRequest(receiver string, pack sync_types.SyncPackage, timeout int) (result sync_types.SyncPackage, err error) {
	return SendRequest(receiver, pack, timeout)
}

func (s *SSyncExhange) SendResponse(packIn *sync_types.SyncPackage, packOut sync_types.SyncPackage) error {
	return SendResponse(packIn, packOut)
}
