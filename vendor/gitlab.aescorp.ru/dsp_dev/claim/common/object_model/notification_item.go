package object_model

import "time"

type NotificationItem struct {
	CommonStruct
	NotifyAt           time.Time `json:"notify_at"            gorm:"column:notify_at;default:null"`
	NotifyAttachmentID int64     `json:"notify_attachment_id" gorm:"column:notify_attachment_id;default:null"`
	NotifyChannel      int       `json:"notify_channel"       gorm:"column:notify_channel;default:null"`
	NotifyCode         int       `json:"notify_code"          gorm:"column:notify_code;default:null"`
	NotifyDone         bool      `json:"notify_done"          gorm:"column:notify_done;default:null"`
	NotifyExtCode      string    `json:"notify_ext_code"      gorm:"column:notify_ext_code;default:null"`
	NotifyMailingCode  string    `json:"notify_mailing_code"  gorm:"column:notify_mailing_code;default:null"`
	NotifyReportID     int64     `json:"notify_report_id"     gorm:"column:notify_report_id;default:null"`
	NotifyTypeID       int64     `json:"notify_type_id"       gorm:"column:notify_type_id;default:null"`
}
