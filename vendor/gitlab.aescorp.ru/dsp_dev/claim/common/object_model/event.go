package object_model

import (
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

// Event Событие для календаря.
type Event struct {
	CommonStruct
	NameStruct
	EventType        EventType `json:"event_type"         gorm:"-:all"`
	EventTypeID      int64     `json:"event_type_id"      gorm:"column:event_type_id;default:null"`
	CreatedBy        Employee  `json:"created_by"        gorm:"-:all"`
	CreatedByID      int64     `json:"created_by_id"      gorm:"column:created_by_id;default:null"`
	ModifiedBy       Employee  `json:"modified_by"        gorm:"-:all"`
	ModifiedByID     int64     `json:"modified_by_id"     gorm:"column:modified_by_id;default:null"`
	DeletedBy        Employee  `json:"deleted_by"         gorm:"-:all"`
	DeletedByID      int64     `json:"deleted_by_id"      gorm:"column:deleted_by_id;default:null"`
	StartAt          time.Time `json:"start_at"           gorm:"column:start_at;default:null"`
	FinishAt         time.Time `json:"finish_at"          gorm:"column:finish_at;default:null"`
	IsAllDay         bool      `json:"is_all_day"         gorm:"column:is_all_day;default:false"`
	RepeatPeriod     int64     `json:"repeat_period"      gorm:"column:repeat_period;default:null"`
	RepeatNumber     int64     `json:"repeat_number"      gorm:"column:repeat_number;default:null"`
	Performer        Employee  `json:"performer"          gorm:"-:all"`
	PerformerID      int64     `json:"performer_id"       gorm:"column:performer_id;default:null"`
	RelativeNoticeAt time.Time `json:"relative_notice_at" gorm:"column:relative_notice_at;default:null"`
	Color            string    `json:"color"              gorm:"column:color;default:null"`
	Priority         int64     `json:"priority"           gorm:"column:priority;default:null"`
}

// TableName -- возвращает имя таблицы в БД, нужен для gorm
func (e Event) TableNameDB() string {
	return "events"
}

// GetID -- возвращает ID объекта
func (e Event) GetID() int64 {
	return e.ID
}

// NewEvent -- Новый объект события
func NewEvent() Event {
	sf := Event{}
	return sf
}

func AsEvent(b []byte) (Event, error) {
	e := Event{}
	err := msgpack.Unmarshal(b, &e)
	if err != nil {
		return Event{}, err
	}
	return e, nil
}

func EventAsBytes(e *Event) ([]byte, error) {
	b, err := msgpack.Marshal(e)
	if err != nil {
		return nil, err
	}
	return b, nil
}
