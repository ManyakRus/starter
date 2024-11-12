package sync_confirm

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_global"
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

type Confirmer interface {
	DeInitConfirm() error
	NewConfirmation(netID string, wait bool) error
	GetConfirmation(netID string) (*Confirmation, error)
	MakeConfirmation(netID string, b bool) error
	SentConfirmation(netID string, b bool) error
	RecvConfirmation(netID string, b bool) error
}

type SyncConfirmer struct {
	db *leveldb.DB
}

var confirmer Confirmer
var block sync.RWMutex

func NewSyncConfirmer(path string) (h Confirmer, err error) {
	block.Lock()
	defer block.Unlock()
	if confirmer == nil {
		var sc SyncConfirmer
		_db, err := leveldb.OpenFile(fmt.Sprintf("%s/%s.db", path, sync_global.SyncService()), nil)
		if err != nil {
			return nil, err
		}
		sc.db = _db
		confirmer = &sc
	}

	return confirmer, err
}

func (s *SyncConfirmer) getIsInited() bool {
	block.RLock()
	defer block.RUnlock()
	return confirmer != nil
}

func (s *SyncConfirmer) DeInitConfirm() error {
	block.Lock()
	defer block.Unlock()

	if confirmer == nil {
		return nil
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("DeInitConfirm, Close, error: %v", err)
	}

	s.db = nil
	confirmer = nil

	return nil
}

func (s *SyncConfirmer) NewConfirmation(netID string, wait bool) error {
	if !s.getIsInited() {
		return errors.New("NewConfirmation, not inited")
	}

	conf := Confirmation{
		CreateAt: time.Now(),
		Wait:     wait,
	}

	value, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = s.db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncConfirmer) GetConfirmation(netID string) (*Confirmation, error) {
	if !s.getIsInited() {
		return nil, fmt.Errorf("GetConfirmation, not inited")
	}

	value, err := s.db.Get([]byte(netID), nil)
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

func (s *SyncConfirmer) MakeConfirmation(netID string, b bool) error {
	if !s.getIsInited() {
		return fmt.Errorf("MakeConfirmation, not inited")
	}

	c, err := s.GetConfirmation(netID)
	if err != nil {
		return err
	}

	c.Make = b
	c.MakeAt = time.Now()

	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = s.db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncConfirmer) SentConfirmation(netID string, b bool) error {
	if !s.getIsInited() {
		return fmt.Errorf("SentConfirmation, not inited")
	}

	if s.db == nil {
		return fmt.Errorf("db is not inited")
	}

	c, err := s.GetConfirmation(netID)
	if err != nil {
		return err
	}

	c.Sent = b
	c.SentAt = time.Now()

	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = s.db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}

func (s *SyncConfirmer) RecvConfirmation(netID string, b bool) error {
	if !s.getIsInited() {
		return fmt.Errorf("RecvConfirmation, not inited")
	}

	c, err := s.GetConfirmation(netID)
	if err != nil {
		return err
	}

	c.Recv = b
	c.RecvAt = time.Now()

	value, err := json.Marshal(c)
	if err != nil {
		return err
	}

	err = s.db.Put([]byte(netID), value, nil)
	if err != nil {
		return err
	}

	return nil
}
