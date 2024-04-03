package models

import (
	"time"

	"gorm.io/gorm"
)

type ClientPreparedCommand struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	ClientID    string     `gorm:"column:client_id;not null;size:255"`
	PowerModeID uint64     `gorm:"column:power_mode_id;not null"`
	Code        int        `gorm:"column:code;type:int"`
	Data        string     `gorm:"column:data;not null"`
	CreatedAt   *time.Time `gorm:"column:created_at;autoCreateTime:milli;not null;default:current_timestamp(3)"`
	UpdatedAt   *time.Time `gorm:"column:updated_at;autoUpdateTime:milli;not null;default:current_timestamp(3);onUpdate:default:current_timestamp(3)"`
}

func (ClientPreparedCommand) TableName() string {
	return "client_prepared_command"
}

func NewClientPreparedCommand(client *Client, mode *PowerMode, code int, data string) *ClientPreparedCommand {
	if client == nil || mode == nil {
		return nil
	}
	return &ClientPreparedCommand{
		ClientID:    client.ID,
		PowerModeID: mode.ID,
		Code:        code,
		Data:        data,
	}
}

func CreateNewClientPreparedCommand(db *gorm.DB, command *ClientPreparedCommand) (int64, error) {
	tx := db.Save(command)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func GetClientPreparedCommand(db *gorm.DB, id uint64) *ClientPreparedCommand {
	var command ClientPreparedCommand
	tx := db.Take(&command, id)
	if tx.Error != nil {
		return nil
	}
	return &command
}

func GetClientPreparedCommands(db *gorm.DB, page, pageSize int, clientID *string) ([]ClientPreparedCommand, int64, error) {
	tx := db.Model(&ClientPreparedCommand{})

	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if clientID != nil {
		tx = tx.Where("client_id = ?", *clientID)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 0
	}
	if pageSize > 0 {
		// 分页
		offset := (page - 1) * pageSize
		tx = tx.Limit(pageSize).Offset(offset)
	}

	// 查询结果
	var records []ClientPreparedCommand
	err = tx.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
