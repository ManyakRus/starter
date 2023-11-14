package sync_confirm

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_global"
	"log"
	"sync"
	"time"
)

type Confirmation struct {
	Request      string
	Response     string
	CreateAt     time.Time
	Wait         bool
	WaitDuration time.Time
	Make         bool
	MakeAt       time.Time
	Sent         bool
	SentAt       time.Time
	Recv         bool
	RecvAt       time.Time
}

var (
	db       *leveldb.DB
	block    sync.RWMutex
	block1   sync.Mutex
	isInited bool
)

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

func InitConfirm(path string) (*leveldb.DB, error) {
	block1.Lock()
	defer block1.Unlock()

	if getIsInited() {
		log.Println("[INFO] InitConfirm, already inited")
		return db, nil
	}

	_db, err := leveldb.OpenFile(fmt.Sprintf("%s/%s.db", path, sync_global.SyncService()), nil)
	if err != nil {
		return nil, fmt.Errorf("InitConfirm, OpenFile, error: %v", err)
	}
	db = _db

	setIsInited(true)

	return db, nil
}

func DeInitConfirm() error {
	block1.Lock()
	defer block1.Unlock()

	if !getIsInited() {
		return fmt.Errorf("DeInitConfirm, not inited")
	}
	defer setIsInited(false)

	err := db.Close()
	if err != nil {
		return fmt.Errorf("DeInitConfirm, Close, error: %v", err)
	}

	db = nil

	return nil
}

func NewConfirmation(db *leveldb.DB, netID string, wait bool) error {
	if !getIsInited() {
		return fmt.Errorf("NewConfirmation, not inited")
	}

	conf := Confirmation{
		CreateAt: time.Now(),
		Wait:     wait,
		// WaitDuration
	}

	value, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}

func GetConfirmation(db *leveldb.DB, netID string) (*Confirmation, error) {
	if !getIsInited() {
		return nil, fmt.Errorf("GetConfirmation, not inited")
	}

	value, err := db.Get([]byte(netID), nil)
	if err != nil {
		return nil, err
	}

	conf := Confirmation{}
	err = json.Unmarshal(value, &conf)
	if err != nil {
		return nil, fmt.Errorf("GetConfirmation(): in unmarshall JSON, err=%w", err)
	}

	return &conf, nil
}

func MakeConfirmation(db *leveldb.DB, netID string, b bool) error {
	if !getIsInited() {
		return fmt.Errorf("MakeConfirmation, not inited")
	}

	c, err := GetConfirmation(db, netID)
	if err != nil {
		return err
	}

	c.Make = b
	c.MakeAt = time.Now()

	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}

func SentConfirmation(db *leveldb.DB, netID string, b bool) error {
	if !getIsInited() {
		return fmt.Errorf("SentConfirmation, not inited")
	}

	if db == nil {
		return fmt.Errorf("db is not inited")
	}

	c, err := GetConfirmation(db, netID)
	if err != nil {
		return err
	}

	c.Sent = b
	c.SentAt = time.Now()

	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}

func RecvConfirmation(db *leveldb.DB, netID string, b bool) error {
	if !getIsInited() {
		return fmt.Errorf("RecvConfirmation, not inited")
	}

	c, err := GetConfirmation(db, netID)
	if err != nil {
		return err
	}

	c.Recv = b
	c.RecvAt = time.Now()

	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}
