package models

import (
	"time"

	"gorm.io/gorm"
)

type PowerModeExecution struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	PowerModeID uint64     `gorm:"column:power_mode_id;not null"`
	CreatedAt   *time.Time `gorm:"column:created_at;autoCreateTime:milli;not null;default:current_timestamp(3)"`
}

func (PowerModeExecution) TableName() string {
	return "power_mode_execution"
}

func NewPowerModeExecution(db *gorm.DB, mode *PowerMode) *PowerModeExecution {
	if mode == nil {
		return nil
	}
	return &PowerModeExecution{
		PowerModeID: mode.ID,
	}
}

// CreateNewPowerModeExecution 创建新的能耗模式执行记录，并插入数据库。
func CreateNewPowerModeExecution(db *gorm.DB, mode *PowerMode) (int64, error) {
	if mode == nil {
		return 0, gorm.ErrRecordNotFound
	}
	execution := NewPowerModeExecution(db, mode)
	tx := db.Save(execution)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// GetPowerModeExecutions 获取能耗模式执行记录列表。
func GetPowerModeExecutions(db *gorm.DB, page, pageSize int) ([]PowerModeExecution, int64, error) {
	tx := db.Model(&PowerModeExecution{})

	var total int64
	err := tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
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
	var records []PowerModeExecution
	err = tx.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
