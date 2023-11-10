package sync_types

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gitlab.aescorp.ru/dsp_dev/claim/common/sync_exchange/sync_global"
	"time"
)

// SyncPackage Пакет. Содержит заголовок и тело.
type SyncPackage struct {
	Head SyncHead `json:"head"`
	Body SyncBody `json:"body"`
}

// SyncHead Заголовок пакета. Содержит данные для идентификации.
type SyncHead struct {
	DestVer string `json:"SyncDestVer"`
	Sender  string `json:"sender"`
	NetID   string `json:"netID"`
	Created string `json:"created"`
}

// SyncBody Тело пакета. Содержит подобъекты согласно назначению.
type SyncBody struct {
	Command string     `json:"command,omitempty"`
	Params  SyncParams `json:"params,omitempty"`
	Result  SyncResult `json:"result,omitempty"`
	Error   SyncError  `json:"error,omitempty"`
	Object  SyncObject `json:"object,omitempty"`
}

// SyncParams Параметры команды.
type SyncParams map[string]interface{}

// SyncResult Результат выполнения команды.
type SyncResult map[string]interface{}

// SyncError Структура содержащая ошибку.
type SyncError struct {
	Place   string `json:"place"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// SyncObject Объектная модель.
type SyncObject []byte

func makeSyncHead(sender string) SyncHead {
	dt := time.Now().Local().Format("2006-01-02 15:04:05.000")

	id := uuid.New().String()

	result := SyncHead{sync_global.SyncDestVer, sender, id, dt}

	return result
}

// IsValid Validate SyncPackage header
func (p *SyncPackage) IsValid() bool {
	_ver := sync_global.SyncDestVer
	return p.Head.DestVer == _ver
}

// IsCommand Check SyncPackage is command package
func (p *SyncPackage) IsCommand() bool {
	return p.Body.Command != ""
}

// IsResult Check SyncPackage is result package
func (p *SyncPackage) IsResult() bool {
	return len(p.Body.Result) != 0
}

// IsError Check SyncPackage is error package
func (p *SyncPackage) IsError() bool {
	return (p.Body.Error.Code != 0) || (p.Body.Error.Message != "")
}

// SyncPackageToJSON View SyncPackage as JSON string
func SyncPackageToJSON(p *SyncPackage) (string, error) {
	if p == nil {
		return "", fmt.Errorf("SyncPackage is nil")
	}

	result, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	// DEBUG
	// log.Println(string(result))

	return string(result), nil
}

// SyncPackageFromJSON Make SyncPackage from JSON string
func SyncPackageFromJSON(msg string) (SyncPackage, error) {
	// DEBUG
	// log.Println(string(msg))

	result := SyncPackage{}
	err := json.Unmarshal([]byte(msg), &result)
	if err != nil {
		return result, fmt.Errorf("SyncPackageFromJSON(): in unmarshall JSON, err=%w", err)
	}

	return result, nil
}

func NewSyncParams() SyncParams {
	return make(SyncParams)
}

func NewSyncResult() SyncResult {
	return make(SyncResult)
}

// MakeSyncCommand Create SyncPackage as command package
func MakeSyncCommand(command string, params SyncParams) SyncPackage {
	_head := makeSyncHead(sync_global.SyncService())
	_body := SyncBody{Command: command, Params: params}
	_result := SyncPackage{_head, _body}

	return _result
}

// MakeSyncResult Create SyncPackage as result package
func MakeSyncResult(result SyncResult) SyncPackage {
	_head := makeSyncHead(sync_global.SyncService())
	_body := SyncBody{Result: result}
	_result := SyncPackage{_head, _body}

	return _result
}

// MakeSyncError Create SyncPackage as error package
func MakeSyncError(place string, code int, message string) SyncPackage {
	_head := makeSyncHead(sync_global.SyncService())
	_body := SyncBody{Error: SyncError{place, code, message}}
	_result := SyncPackage{_head, _body}

	return _result
}

// MakeSyncObject Create SyncPackage as object package
func MakeSyncObject(object *SyncObject) SyncPackage {
	_head := makeSyncHead(sync_global.SyncService())
	_body := SyncBody{Object: *object}
	_result := SyncPackage{_head, _body}

	return _result
}
