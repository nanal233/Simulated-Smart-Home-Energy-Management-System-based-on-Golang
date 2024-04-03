package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID        string     `gorm:"primaryKey;size:255;column:id"`
	Name      string     `gorm:"column:name;size:255;not null"`
	Type      ClientType `gorm:"column:type;type:int:not null"`
	CreatedAt *time.Time `gorm:"column:created_at;autoCreateTime:milli;not null;default:current_timestamp(3)"`
	UpdatedAt *time.Time `gorm:"column:updated_at;autoUpdateTime:milli;not null;default:current_timestamp(3);onUpdate:default:current_timestamp(3)"`

	ClientActivities        []ClientActivity         `gorm:"foreignKey:client_id" json:"-"`
	ClientCommandExecutions []ClientCommandExecution `gorm:"foreignKey:client_id" json:"-"`

	IsActive bool `gorm:"-"`
}

func NewClient(id string, name string, clientType ClientType) *Client {
	return &Client{
		ID:   id,
		Name: name,
		Type: clientType,
	}
}

// RegisterNewClient registers a new Client in the database.
// It accepts a database connection and a Client struct as input.
// It returns the number of rows affected (should be 1 for a successful insertion)
// and any error that occurred during the operation.
func RegisterNewClient(db *gorm.DB, client *Client) (int64, error) {
	if client == nil {
		return 0, gorm.ErrRecordNotFound
	}
	tx := db.Save(client)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// RemoveClient removes an existing Client from the database.
// It accepts a database connection and a Client struct with a valid ID as input.
// It returns the number of rows affected (should be 1 for a successful deletion)
// and any error that occurred during the operation.
func RemoveClient(db *gorm.DB, client *Client) (int64, error) {
	if client == nil {
		return 0, gorm.ErrRecordNotFound
	}
	tx := db.Delete(client, client.ID)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// TableName 表名
func (*Client) TableName() string {
	return "client"
}

func GetClient(db *gorm.DB, id string) (*Client, error) {
	var client Client
	err := db.Take(&client, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetClients 获取客户端列表。
func GetClients(db *gorm.DB, page, pageSize int, clientType ClientType) ([]Client, int64, error) {
	tx := db.Model(&Client{})

	if clientType > 0 {
		tx = tx.Where("type = ?", clientType)
	}

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
	var records []Client
	err = tx.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil

}

// ClientExists 检查客户端是否存在。
func ClientExists(db *gorm.DB, id string) bool {
	var count int64
	err := db.Model(&Client{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// GetActivities 返回当前 Client 的所有 Activity 列表，支持翻页和根据 status 查询。
func (c *Client) GetActivities(db *gorm.DB, page, pageSize int, status *int8) ([]ClientActivity, int, error) {
	// 创建查询对象
	tx := db.Model(&ClientActivity{}).Where("client_id = ?", c.ID)

	if status != nil {
		tx = tx.Where("status = ?", *status)
	}

	// 查询总数
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
	var records []ClientActivity
	err = tx.Order("created_at desc").Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, int(total), nil
}

func (c *Client) InsertNewActivity(db *gorm.DB, status int8) (int64, error) {
	record := &ClientActivity{
		ClientID: c.ID,
		Status:   status,
	}
	tx := db.Save(record)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// GetCommandExecutions 返回当前 Client 的所有 Command 执行历史列表，支持翻页和根据 code 查询。
func (c *Client) GetCommandExecutions(db *gorm.DB, page, pageSize int, code *int) ([]ClientCommandExecution, int, error) {
	// 创建查询对象
	tx := db.Model(&ClientCommandExecution{}).Where("client_id = ?", c.ID)

	if code != nil {
		tx = tx.Where("code = ?", *code)
	}

	// 查询总数
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
	var records []ClientCommandExecution
	err = tx.Order("created_at desc").Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, int(total), nil
}

// InsertNewCommandExecution 插入一条命令执行历史。
func (c *Client) InsertNewCommandExecution(db *gorm.DB, code int, data string, sentAt *time.Time) (int64, error) {
	now := time.Now()
	if sentAt == nil {
		sentAt = &now
	}
	record := &ClientCommandExecution{
		ClientID: c.ID,
		Code:     code,
		Data:     data,
		SentAt:   *sentAt,
	}
	tx := db.Save(record)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

type ClientType int

func (c *Client) GetConsumptions(db *gorm.DB, page, pageSize int) ([]ClientConsumption, int64, error) {
	// 创建查询对象
	tx := db.Model(&ClientConsumption{}).Where("client_id = ?", c.ID)

	// 查询总数
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
	var records []ClientConsumption
	err = tx.Order("created_at desc").Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (c *Client) InsertNewConsumption(db *gorm.DB, consumption float32, recordedAt time.Time) (int64, error) {
	record := &ClientConsumption{
		ClientID:    c.ID,
		Consumption: consumption,
		RecordedAt:  recordedAt,
	}
	tx := db.Save(record)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// UpdateName 更新当前客户端的名称。
func (c *Client) UpdateName(db *gorm.DB, name string) (int64, error) {
	c.Name = name
	tx := db.Model(c).Update("name", name)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// Delete 删除自身。
func (c *Client) Delete(db *gorm.DB) (int64, error) {
	tx := db.Delete(c)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
