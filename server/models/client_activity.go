package models

import "time"

type ClientActivity struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	ClientID  string     `gorm:"column:client_id;size:255;not null"`
	Status    int8       `gorm:"column:status;not null"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime:milli;not null;default:current_timestamp(3)"`

	// 多对多关系
	ClientCommandExecutions []ClientCommandExecution `gorm:"many2many:Client"`
}

// TableName 表名
func (ClientActivity) TableName() string {
	return "client_activity"
}
