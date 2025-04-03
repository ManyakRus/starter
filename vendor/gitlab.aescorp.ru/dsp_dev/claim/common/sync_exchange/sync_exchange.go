package sync_exchange

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/data_packer"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/liveness"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_confirm"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_global"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_types"
)

// PRIVATE

var (
	nc        *nats.Conn
	packer    *data_packer.DataPacker
	block     sync.RWMutex
	block1    sync.Mutex
	isInited  bool
	confirmer sync_confirm.Confirmer
)

func GetUseConfirmerEnv() bool {
	val, ok := os.LookupEnv("SYNC_EXCHANGE_USE_CONFIRMER")
	if !ok {
		return false
	} else {
		return strings.EqualFold("true", val)
	}
}

func fullTopic(topic string) string {
	return sync_global.SyncRoot + topic + "/"
}

func setIsInited(b bool) {
	block.Lock()
	defer block.Unlock()
	isInited = b
}

func getIsInited() bool {
	block.RLock()
	defer block.RUnlock()
	return isInited
}

// doSendMessage Непосредственно отправка сообщения
func doSendMessage(topic string, pack sync_types.SyncPackage, wait bool) error {
	// Новое сообщение
	err := confirmer.NewConfirmation(pack.Head.NetID, wait)
	if err != nil {
		// TODO: Лог
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}

	// Получаем JSON
	rawData, err := sync_types.SyncPackageToJSON(&pack)
	if err != nil {
		// Создание сообщения неудачно
		err1 := confirmer.MakeConfirmation(pack.Head.NetID, false)
		if err1 != nil {
			return fmt.Errorf("doSendMessage, Error: %v, Error1: %v", err, err1)
		}
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}
	// Пакуем сообщение
	data, err := packer.Pack([]byte(rawData))
	if err != nil {
		// Упаковка сообщения неудачна
		err1 := confirmer.MakeConfirmation(pack.Head.NetID, false)
		if err1 != nil {
			return fmt.Errorf("doSendMessage, Error: %v, Error1: %v", err, err1)
		}
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}
	// Создание сообщения удачно
	err = confirmer.MakeConfirmation(pack.Head.NetID, true)
	if err != nil {
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}

	// Отправка сообщения
	err = nc.Publish(topic, data)
	if err != nil {
		// Отправка сообщения неудачна
		err1 := confirmer.SentConfirmation(pack.Head.NetID, false)
		if err1 != nil {
			return fmt.Errorf("doSendMessage, Error: %v, Error1: %v", err, err1)
		}
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}
	// Отправка сообщения удачна
	err = confirmer.SentConfirmation(pack.Head.NetID, true)
	if err != nil {
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}

	return nil
}

// doWaitMessage Ожидание сообщения
func doWaitMessage(topic string, queue string, callback Callback) error {
	// log.Printf("[INFO] WaitMessage, topic: %s\n", topic)

	_, err := nc.QueueSubscribe(topic, queue, func(msg *nats.Msg) {
		_data, _ := packer.Unpack(msg.Data)
		pack, err := sync_types.SyncPackageFromJSON(string(_data))
		pack.Msg = msg
		if err != nil {
			log.Println(err)
			return
		}

		go callback(&pack)
	})

	if err != nil {
		return fmt.Errorf("doWaitMessage, Error: %v", err)
	}

	return nil
}

// PUBLIC

// Callback Функция возврата для подписки на события шины
type Callback func(pack *sync_types.SyncPackage)

// InitSyncExchange Функция инициализации подключения к шине
func InitSyncExchange(url string, serviceName string, version string) error {
	block1.Lock()
	defer block1.Unlock()

	log.Printf("[INFO] sync_exchange, InitSyncExchange, url: %v, service: %v, version: %v", url, serviceName, version)

	if getIsInited() {
		log.Println("[INFO] sync_exchange, InitSyncExchange, already inited")
		return nil
	}

	err := sync_global.SetSyncService(serviceName)
	if err != nil {
		return fmt.Errorf("InitSyncExchange, SetSyncService, error: %v", err)
	}

	_packer := data_packer.NewDataPacker()
	packer = _packer

	_nc, err := nats.Connect(url, nats.Name(serviceName))
	if err != nil {
		return fmt.Errorf("InitSyncExchange, Connect, error: %v", err)
	}
	nc = _nc

	status := nc.Status()
	switch status {
	case nats.DISCONNECTED, nats.CLOSED:
		return fmt.Errorf("InitSyncExchange, NATS connection status: %v", status.String())
	default:
		log.Printf("[INFO] sync_exchange, InitSyncExchange, NATS connection status: %v\n", status.String())
	}

	// TODO Вынести путь в параметр функции
	storePath := "./store"
	// TODO Тут обработать не подтверждённые пакеты

	if GetUseConfirmerEnv() {
		confirmer, err = sync_confirm.NewSyncConfirmer(storePath)
	} else {
		confirmer, err = sync_confirm.NewNoConfirmer(storePath)
	}
	if err != nil {
		return fmt.Errorf("InitSyncExchange, NewConfirmer, path: %q, error: %v", storePath, err)
	}

	setIsInited(true)

	go liveness.RunLiveness(nc, serviceName, version)

	return nil
}

// DeInitSyncExchange Функция де-инициализации подключения к шине
func DeInitSyncExchange() error {
	block1.Lock()
	defer block1.Unlock()

	if !getIsInited() {
		return fmt.Errorf("DeInitSyncExchange, not inited")
	}
	defer setIsInited(false)

	nc.Close()

	err := confirmer.DeInitConfirm()
	if err != nil {
		return fmt.Errorf("DeInitSyncExchange, DeInitConfirm, error: %v", err)
	}

	nc = nil
	packer = nil

	return nil
}

// SendMessage Отправка сообщения в шину без ожидания ответа
func SendMessage(topic string, pack sync_types.SyncPackage) error {
	block1.Lock()
	defer block1.Unlock()

	if !getIsInited() {
		return fmt.Errorf("SendMessage, not inited")
	}

	err := doSendMessage(topic, pack, false)
	if err != nil {
		return fmt.Errorf("SendMessage, Error: %v", err)
	}

	return nil
}

// WaitMessage Ожидание сообщения из определённого топика
func WaitMessage(topic string, callback Callback) error {
	if !getIsInited() {
		return fmt.Errorf("WaitMessage, not inited")
	}

	_topic := topic
	if !strings.HasPrefix(_topic, sync_global.SyncRoot) {
		_topic = fullTopic(topic)
	}

	return doWaitMessage(_topic, sync_global.SyncQueue, callback)
}

// QueueSubscribe Ожидание сообщения из определённого топика
func QueueSubscribe(topic string, queue string, callback Callback) error {
	if !getIsInited() {
		return fmt.Errorf("WaitMessage, not inited")
	}

	return doWaitMessage(topic, queue, callback)
}

// Subscribe Ожидание сообщения из определённого топика
func Subscribe(topic string, callback Callback) error {
	if !getIsInited() {
		return fmt.Errorf("subscribe, not inited")
	}

	_, err := nc.Subscribe(topic, func(msg *nats.Msg) {
		_data, _ := packer.Unpack(msg.Data)
		pack, err := sync_types.SyncPackageFromJSON(string(_data))
		pack.Msg = msg
		if err != nil {
			log.Println(err)
			return
		}
		go callback(&pack)
	})

	if err != nil {
		return fmt.Errorf("subscribe, Error: %v", err)
	}

	return nil
}

// SendRequest Отправка запроса с ожиданием ответа
func SendRequest(receiver string, pack sync_types.SyncPackage, timeout int) (result sync_types.SyncPackage, err error) {
	result = sync_types.MakeSyncError("", 0, "")

	if !getIsInited() {
		return result, fmt.Errorf("SendRequest, not inited")
	}

	// Новое сообщение
	if err = confirmer.NewConfirmation(pack.Head.NetID, true); err != nil {
		log.Printf("[ERROR] SendRequest, NewConfirmation error: %s\n", err.Error())
	}

	_topic := fullTopic(receiver)

	rawData, err := sync_types.SyncPackageToJSON(&pack)
	if err != nil {
		// Создание сообщения неудачно
		if err1 := confirmer.MakeConfirmation(pack.Head.NetID, false); err1 != nil {
			log.Printf("[ERROR] SendRequest, SyncPackageToJSON error: %s, MakeConfirmation error: %s\n", err.Error(), err1.Error())
		}
		log.Printf("[ERROR] SendRequest, SyncPackageToJSON error: %s\n", err.Error())
		return result, err
	}

	// Пакуем сообщение
	data, err := packer.Pack([]byte(rawData))
	if err != nil {
		// Упаковка сообщения неудачна
		if err1 := confirmer.MakeConfirmation(pack.Head.NetID, false); err1 != nil {
			log.Printf("[ERROR] SendRequest, Pack error: %s, MakeConfirmation error: %s\n", err.Error(), err1.Error())
		}
		log.Printf("[ERROR] SendRequest, Pack error: %s\n", err.Error())
		return result, err
	}

	// Создание сообщения удачно
	if err = confirmer.MakeConfirmation(pack.Head.NetID, true); err != nil {
		log.Printf("[ERROR] SendRequest ok, MakeConfirmation error: %s\n", err.Error())
	}

	if timeout == -1 {
		timeout = 24 * 60 * 60 * 1000
	}

	msg, err := nc.Request(_topic, data, time.Duration(timeout)*time.Second)
	if err != nil {
		if err1 := confirmer.SentConfirmation(pack.Head.NetID, false); err1 != nil {
			log.Printf("[ERROR] SendRequest (%v), Request error: %s, SentConfirmation error: %s\n", _topic, err.Error(), err1.Error())
		}
		log.Printf("[ERROR] SendRequest (%v), Request error: %s\n", _topic, err.Error())
		return result, err
	}

	// Отправка сообщения удачна
	if err = confirmer.SentConfirmation(pack.Head.NetID, true); err != nil {
		log.Printf("[ERROR] SendRequest, Request ok, SentConfirmation error: %s\n", err.Error())
	}

	_data, err := packer.Unpack(msg.Data)
	if err != nil {
		log.Printf("[ERROR] SendRequest, Unpack, error: %s\n", err.Error())
		return result, err
	}
	result, err = sync_types.SyncPackageFromJSON(string(_data))
	if err != nil {
		log.Printf("[ERROR] SendRequest, SyncPackageFromJSON, error: %s\n", err.Error())
		result.Body.Error.Code = 3
		if err == nats.ErrTimeout {
			result.Body.Error.Code = 4
		}
		return result, err
	}

	if err = confirmer.RecvConfirmation(pack.Head.NetID, true); err != nil {
		log.Printf("[ERROR] SendRequest, Request ok, RecvConfirmation error: %s\n", err.Error())
	}

	return result, nil
}

// SendResponse Отправка ответа на запрос
func SendResponse(packIn *sync_types.SyncPackage, packOut sync_types.SyncPackage) error {
	if !getIsInited() {
		return fmt.Errorf("SendResponse, not inited")
	}

	if packOut.Body.Result == nil {
		packOut.Body.Result = make(sync_types.SyncResult)
	}
	packOut.Body.Result["netID"] = packIn.Head.NetID

	// Новое сообщение
	if err := confirmer.NewConfirmation(packIn.Head.NetID, true); err != nil {
		log.Printf("[ERROR] SendResponse, NewConfirmation error: %s\n", err.Error())
	}

	msg := packIn.Msg
	if msg == nil {
		return fmt.Errorf("SendResponse, Error: packIn.Msg is nil")
	}

	rawData, err := sync_types.SyncPackageToJSON(&packOut)
	if err != nil {
		// Создание сообщения неудачно
		if err1 := confirmer.MakeConfirmation(packIn.Head.NetID, false); err1 != nil {
			log.Printf("[ERROR] SendResponse, SyncPackageToJSON error: %s, MakeConfirmation error: %s\n", err.Error(), err1.Error())
		}
		log.Printf("[ERROR] SendResponse SyncPackageToJSON error: %s\n", err.Error())
		return fmt.Errorf("SendResponse, SyncPackageToJSON error: %v", err)
	}
	// Пакуем сообщение
	data, err := packer.Pack([]byte(rawData))
	if err != nil {
		// Упаковка сообщения неудачна
		if err1 := confirmer.MakeConfirmation(packIn.Head.NetID, false); err1 != nil {
			log.Printf("[ERROR] SendResponse Pack error: %s, MakeConfirmation error: %s\n", err.Error(), err1.Error())
		}
		log.Printf("[ERROR] SendResponse, Pack error: %s\n", err.Error())
		return fmt.Errorf("SendResponse, Pack error: %v", err)
	}

	err = msg.Respond(data)
	if err != nil {
		if err1 := confirmer.SentConfirmation(packIn.Head.NetID, false); err1 != nil {
			log.Printf("[ERROR] SendResponse, Respond error: %s, SentConfirmation error: %s\n", err.Error(), err1.Error())
		}
		log.Printf("[ERROR] SendResponse, Respond error: %s\n", err.Error())
		return fmt.Errorf("SendResponse, Respond error: %v", err)
	}

	// Отправка сообщения удачна
	if err = confirmer.SentConfirmation(packIn.Head.NetID, true); err != nil {
		log.Printf("[ERROR] SendResponse ok, SentConfirmation error: %s\n", err.Error())
	}

	return nil
}

func SendRawMessage(topic string, data []byte) error {
	if !getIsInited() {
		return fmt.Errorf("SendRawMessage, not inited")
	}

	// Отправка сообщения
	err := nc.Publish(topic, data)
	if err != nil {
		return fmt.Errorf("SendRawMessage, Error: %v", err)
	}
	return nil
}
