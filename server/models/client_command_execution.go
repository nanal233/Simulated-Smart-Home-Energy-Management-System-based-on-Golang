package models

import (
	"time"
)

// ClientCommandExecution 表示客户端命令执行历史。
type ClientCommandExecution struct {
	ID        uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	ClientID  string     `gorm:"column:client_id;size:255;not null"`
	Code      int        `gorm:"column:code;not null"`
	Data      string     `gorm:"column:data;type:text;not null"`
	SentAt    time.Time  `gorm:"column:sent_at;type:timestamp;not null"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime:milli;not null;default:current_timestamp(3)"`

	// 多对多关系
	ClientActivities []ClientActivity `gorm:"many2many:Client"`
}

// TableName 表名
func (ClientCommandExecution) TableName() string {
	return "client_command_execution"
}
