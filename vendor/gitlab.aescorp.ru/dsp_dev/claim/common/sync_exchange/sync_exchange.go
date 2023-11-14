package sync_exchange

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/syndtr/goleveldb/leveldb"

	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/data_packer"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_confirm"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_global"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_types"
)

// PRIVATE

var (
	nc       *nats.Conn
	packer   *data_packer.DataPacker
	db       *leveldb.DB
	block    sync.RWMutex
	block1   sync.Mutex
	isInited bool
)

func requestTopic(topic string) string {
	return sync_global.SyncRoot + topic + "/"
}

func responseTopic(pack *sync_types.SyncPackage) string {
	return sync_global.SyncRoot + pack.Head.Sender + "/" + pack.Head.NetID + "/"
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
	// block1.Lock()
	// defer block1.Unlock()

	// _json, _ := sync_types.SyncPackageToJSON(&pack)
	// log.Printf("[DEBUG] SendMessage, topic: %s, message:\n\t%s", topic, _json)

	// Новое сообщение
	err := sync_confirm.NewConfirmation(db, pack.Head.NetID, wait)
	if err != nil {
		// TODO: Лог
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}

	// Получаем JSON
	raw_data, err := sync_types.SyncPackageToJSON(&pack)
	if err != nil {
		// Создание сообщения неудачно
		err1 := sync_confirm.MakeConfirmation(db, pack.Head.NetID, false)
		if err1 != nil {
			return fmt.Errorf("doSendMessage, Error: %v, Error1: %v", err, err1)
		}
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}
	// Пакуем сообщение
	data, err := packer.Pack([]byte(raw_data))
	if err != nil {
		// Упаковка сообщения неудачна
		err1 := sync_confirm.MakeConfirmation(db, pack.Head.NetID, false)
		if err1 != nil {
			return fmt.Errorf("doSendMessage, Error: %v, Error1: %v", err, err1)
		}
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}
	// Создание сообщения удачно
	err = sync_confirm.MakeConfirmation(db, pack.Head.NetID, true)
	if err != nil {
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}

	// Отправка сообщения
	err = nc.Publish(topic, data)
	if err != nil {
		// Отправка сообщения неудачна
		err1 := sync_confirm.SentConfirmation(db, pack.Head.NetID, false)
		if err1 != nil {
			return fmt.Errorf("doSendMessage, Error: %v, Error1: %v", err, err1)
		}
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}
	// Отправка сообщения удачна
	err = sync_confirm.SentConfirmation(db, pack.Head.NetID, true)
	if err != nil {
		return fmt.Errorf("doSendMessage, Error: %v", err)
	}

	return nil
}

// doWaitMessage Ожидание сообщения
func doWaitMessage(topic string, callback Callback) error {
	// log.Printf("[INFO] WaitMessage, topic: %s\n", topic)

	// TODO: обернуть в свой обработчик запросов, чтобы обработать ошибку
	_, err := nc.Subscribe(topic, func(msg *nats.Msg) {
		_data, _ := packer.Unpack(msg.Data)
		pack, err := sync_types.SyncPackageFromJSON(string(_data))
		if err != nil {
			log.Println(err)
			return
		}

		// netID := pack.Body.Result["netID"]
		// if netID != "" {
		// 	err = sync_confirm.RecvConfirmation(db, netID, true)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// }

		callback(&pack)

		// if netID != "" {
		// 	c, err := sync_confirm.GetConfirmation(db, netID)
		// 	if err != nil {
		// 		log.Println(err)
		// 	}
		// 	log.Printf("[DEBUG] Message: %v, CreateAt: %v, MakeAt: %v, SentAt: %v, RecvAt: %v\n", netID, c.CreateAt, c.MakeAt, c.SentAt, c.RecvAt)
		// }
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
func InitSyncExchange(url string, serviceName string) error {
	block1.Lock()
	defer block1.Unlock()

	if getIsInited() {
		log.Println("[INFO] InitSyncExchange, already inited")
		return nil
	}

	err := sync_global.SetSyncService(serviceName)
	if err != nil {
		return fmt.Errorf("InitSyncExchange, SetSyncService, error: %v", err)
	}

	_packer := data_packer.NewDataPacker()
	packer = _packer

	_nc, err := nats.Connect(url)
	if err != nil {
		return fmt.Errorf("InitSyncExchange, Connect, error: %v", err)
	}
	nc = _nc

	status := nc.Status()
	switch status {
	case nats.DISCONNECTED, nats.CLOSED:
		return fmt.Errorf("InitSyncExchange, NATS connection status: %v\n", status.String())
	default:
		log.Printf("[INFO] InitSyncExchange, NATS connection status: %v\n", status.String())
	}

	// TODO Вынести путь в параметр функции
	storePath := "./store"
	// TODO Тут обработать не подтверждённые пакеты
	_db, err := sync_confirm.InitConfirm(storePath)
	if err != nil {
		return fmt.Errorf("InitSyncExchange, InitConfirm, path: %q, error: %v", storePath, err)
	}
	db = _db

	setIsInited(true)

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

	err := sync_confirm.DeInitConfirm()
	if err != nil {
		return fmt.Errorf("DeInitSyncExchange, DeInitConfirm, error: %v", err)
	}

	nc = nil
	packer = nil
	db = nil

	return nil
}

// SendMessage Отправка сообщения в шину без ожидания ответа
func SendMessage(topic string, pack sync_types.SyncPackage) error {
	block1.Lock()
	defer block1.Unlock()

	if !getIsInited() {
		return fmt.Errorf("SendMessage, not inited")
	}
	// log.Println("[DEBUG] SendMessage")
	err := doSendMessage(topic, pack, false)
	if err != nil {
		return fmt.Errorf("SendMessage, Error: %v", err)
	}

	return nil
}

// WaitMessage Ожидание сообщения из определённого топика
func WaitMessage(topic string, callback Callback) error {
	// block1.Lock()
	// defer block1.Unlock()

	if !getIsInited() {
		return fmt.Errorf("WaitMessage, not inited")
	}
	// log.Println("[DEBUG] WaitMessage")

	return doWaitMessage(topic, callback)
}

// SendRequest Отправка запроса с ожиданием ответа
func SendRequest(receiver string, pack sync_types.SyncPackage, timeout int) (result sync_types.SyncPackage, err error) {
	block1.Lock()
	defer block1.Unlock()

	result = sync_types.MakeSyncError("", 0, "")

	if !getIsInited() {
		return result, fmt.Errorf("SendRequest, not inited")
	}
	// log.Println("[DEBUG] SendRequest")
	// _time := time.Now()

	blockDone := sync.Mutex{}
	done := false
	topic1 := responseTopic(&pack)
	sub, err := nc.Subscribe(topic1,
		func(msg *nats.Msg) {
			_data, err := packer.Unpack(msg.Data)
			if err != nil {
				log.Printf("[ERROR] SendRequest, Subscribe, Unpack, error: %s\n", err.Error())
			}
			_pack, err := sync_types.SyncPackageFromJSON(string(_data))
			if err != nil {
				log.Printf("[ERROR] SendRequest, Subscribe, SyncPackageFromJSON, error: %s\n", err.Error())
			}

			netID := fmt.Sprintf("%v", _pack.Body.Result["netID"])
			if netID != "" {
				err = sync_confirm.RecvConfirmation(db, netID, true)
				if err != nil {
					log.Println(err)
				}
				_, err := sync_confirm.GetConfirmation(db, netID)
				if err != nil {
					log.Println(err)
				}
				// log.Printf("[DEBUG] Message: %v, CreateAt: %v, MakeAt: %v, SentAt: %v, RecvAt: %v\n", netID, c.CreateAt, c.MakeAt, c.SentAt, c.RecvAt)
			}

			blockDone.Lock()
			result = _pack
			done = true
			blockDone.Unlock()
		})
	if err != nil {
		result.Body.Error.Code = 1
		return result, err
	}

	topic2 := requestTopic(receiver)
	// log.Printf("[DEBUG] SendRequest, SendMessage, receiver: %s, sender: %s, command: %s\n", _topic, pack.Head.Sender, pack.Body.Command)

	// _json, _ := sync_types.SyncPackageToJSON(&pack)
	// log.Printf("[DEBUG] SendRequest, request, topic: %s\n\t%s", topic2, _json)
	// log.Printf("[DEBUG] SendRequest, request, topic: %s\n", topic2)

	err = doSendMessage(topic2, pack, true)
	if err != nil {
		result.Body.Error.Code = 2
		return result, err
	}

	// log.Printf("[DEBUG] SendRequest, wait timeout: %d\n", timeout)
	if timeout == -1 {
		for i := 0; i < 24*60*60*1000; i += 25 {
			blockDone.Lock()
			isDone := done
			blockDone.Unlock()
			if isDone {
				break
			}
			time.Sleep(25 * time.Millisecond)
		}
	} else {
		for i := 0; i < timeout; i += 25 {
			blockDone.Lock()
			isDone := done
			blockDone.Unlock()
			if isDone {
				break
			}
			time.Sleep(25 * time.Millisecond)
		}
	}

	err = sub.Unsubscribe()
	if err != nil {
		result.Body.Error.Code = 3
		return result, err
	}

	if !done {
		result.Body.Error.Code = 4
		return result, fmt.Errorf("timeout while waiting for a response")
	}

	// _json, _ = sync_types.SyncPackageToJSON(&result)
	// log.Printf("[DEBUG] SendRequest, response, duration: %v, topic: %s\n\t%s", time.Since(_time), topic1, _json)
	// log.Printf("[DEBUG] SendRequest, response, duration: %v, topic: %s\n", time.Since(_time), topic1)

	return result, nil
}

// SendResponse Отправка ответа на запрос
func SendResponse(packIn *sync_types.SyncPackage, packOut sync_types.SyncPackage) error {
	if !getIsInited() {
		return fmt.Errorf("SendResponse, not inited")
	}

	// log.Println("[DEBUG] SendResponse")
	_topic := responseTopic(packIn)

	if packOut.Body.Result == nil {
		packOut.Body.Result = make(sync_types.SyncResult)
	}
	packOut.Body.Result["netID"] = packIn.Head.NetID

	// _json, _ := sync_types.SyncPackageToJSON(&packOut)
	// log.Printf("[DEBUG] SendResponse, topic: %s\n\t%s", _topic, _json)

	err := doSendMessage(_topic, packOut, false)
	if err != nil {
		return fmt.Errorf("SendResponse, Error: %v", err)
	}

	return nil
}
