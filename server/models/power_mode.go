package models

import (
	"time"

	"gorm.io/gorm"
)

type PowerMode struct {
	ID        uint64     `gorm:"column:id;primaryKey"`
	Name      string     `gorm:"column:name;not null;size:255"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime:milli;not null;default:current_timestamp(3)"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime:milli;not null;default:current_timestamp(3);onUpdate:default:current_timestamp(3)"`

	ClientPreparedCommands []*ClientPreparedCommand `gorm:"many2many:power_mode_client_prepared_command"`
}

func (*PowerMode) TableName() string {
	return "power_mode"
}

func NewPowerMode(name string) *PowerMode {
	return &PowerMode{
		Name: name,
	}
}

// CreateNewPowerMode 创建新的能耗模式，并保存。
func CreateNewPowerMode(db *gorm.DB, name string) (int64, error) {
	mode := NewPowerMode(name)
	tx := db.Save(mode)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func GetPowerMode(db *gorm.DB, id uint64) *PowerMode {
	var mode PowerMode
	tx := db.Take(&mode, id)
	if tx.Error != nil {
		return nil
	}
	return &mode
}

// GetPowerModeByName 根据名称查找能耗模式。
func GetPowerModeByName(db *gorm.DB, name string) *PowerMode {
	var mode PowerMode
	tx := db.Take(&mode, map[string]interface{}{"name": name})
	if tx.Error != nil {
		return nil
	}
	return &mode
}

func GetPowerModes(db *gorm.DB, page, pageSize int) ([]PowerMode, int64, error) {
	tx := db.Model(&PowerMode{})

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
	var records []PowerMode
	err = tx.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (m *PowerMode) UpdateName(db *gorm.DB, name string) (int64, error) {
	m.Name = name
	tx := db.Model(m).Update("name", name)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

func RemovePowerMode(db *gorm.DB, mode *PowerMode) (int64, error) {
	if mode == nil {
		return 0, gorm.ErrRecordNotFound
	}
	tx := db.Delete(mode, mode.ID)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// Execute 执行能耗模式。TODO:
// 1. 获取当前能耗模式的预制命令列表。
// 2. 分别执行每个命令；如果客户端不存在，则跳过。
// 3. 记录成功执行的命令。并返回执行的总数。
func (m *PowerMode) Execute(db *gorm.DB) (int64, error) {

	return 0, nil
}

//
//func (m *PowerMode) AddCommand(db *gorm.DB, command *ClientPreparedCommand) (int64, error) {
//	record := NewPowerModeClientPreparedCommand(db, m, command)
//	if record == nil {
//		return 0, gorm.ErrRecordNotFound
//	}
//
//	tx := db.Save(record)
//	if tx.Error != nil {
//		return 0, tx.Error
//	}
//	return tx.RowsAffected, nil
//}
//
//func (m *PowerMode) RemoveCommand(db *gorm.DB, command *ClientPreparedCommand) (int64, error) {
//	if command == nil {
//		return 0, gorm.ErrRecordNotFound
//	}
//	tx := db.Where("command_id = ?", command.ID).Where("power_mode_id = ?", m.ID).Delete(&PowerModeClientPreparedCommand{})
//	if tx.Error != nil {
//		return 0, tx.Error
//	}
//	return tx.RowsAffected, nil
//}
//
//func (m *PowerMode) GetCommands(db *gorm.DB) ([]ClientPreparedCommand, int64, error) {
//	var commands []ClientPreparedCommand
//	var relations []PowerModeClientPreparedCommand
//	tx := db.Find(&relations, "power_mode_id = ?", m.ID)
//	if tx.Error != nil {
//		return commands, 0, tx.Error
//	}
//	return commands, tx.RowsAffected, nil
//}
