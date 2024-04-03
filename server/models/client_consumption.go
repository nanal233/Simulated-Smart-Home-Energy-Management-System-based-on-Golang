package models

import "time"

type ClientConsumption struct {
	ID          uint64     `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	ClientID    string     `json:"-" gorm:"column:client_id;size:255;not null"`
	Consumption float32    `gorm:"column:consumption;not null"`
	RecordedAt  time.Time  `gorm:"column:recorded_at;not null"`
	CreatedAt   *time.Time `json:"-" gorm:"column:created_at;autoCreateTime:milli;not null;default:current_timestamp(3)"`
}

// TableName 表名
func (ClientConsumption) TableName() string {
	return "client_consumption"
}
