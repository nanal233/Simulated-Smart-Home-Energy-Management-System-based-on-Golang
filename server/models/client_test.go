package models

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var dbPrepared sync.Once
var client *Client
var clientPrepared sync.Once

func prepareDatabase() {
	// 初始化数据库连接
	// 需要指定 parseTime=True，否则无法解析时间。
	dsn := "root:123456@tcp(1.n.rho.im:13406)/project20240227?charset=utf8mb4&parseTime=True&loc=Local"
	_db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	db = _db
}

func prepareClient() {
	c := NewClient("test", "test", 1)

	err := db.Save(c).Error
	if err != nil {
		panic(err)
	}
	client = c
}

func TestMain(m *testing.M) {
	dbPrepared.Do(prepareDatabase)
	db.Begin()
	db.Exec("DELETE FROM `power_mode_client_prepared_command`")
	db.Exec("DELETE FROM `power_mode_execution`")
	db.Exec("DELETE FROM `power_mode`")
	db.Exec("DELETE FROM `client_prepared_command`")
	db.Exec("DELETE FROM `client_command_execution`")
	db.Exec("DELETE FROM `client_activity`")
	db.Exec("DELETE FROM `client`")
	clientPrepared.Do(prepareClient)
	m.Run()
	db.Rollback()
}

func setUpAll(t *testing.T) {
}

func tearDown(t *testing.T) {
}
func TestCreateClient(t *testing.T) {
	setUpAll(t)
	defer tearDown(t)

	client := Client{
		ID:   "12345",
		Name: "Test Client",
		Type: 1,
	}

	err := db.Create(&client).Error
	if err != nil {
		t.Errorf("Error creating client: %v", err)
	}
}
func TestUpdateClient(t *testing.T) {
	setUpAll(t)
	defer tearDown(t)

	client := Client{
		ID:   "67890",
		Name: "Updated Client Name",
		Type: 1,
	}

	err := db.Save(&client).Error
	if err != nil {
		t.Errorf("Error updating client: %v", err)
	}

	var client1 Client
	tx := db.First(&client1, "67890")
	if tx.Error != nil {
		t.Errorf("Failed to take the existed: %v", tx.Error)
	}
	assert.Equal(t, client.ID, client1.ID)
	assert.Equal(t, client.Name, client1.Name)

	err = db.Model(&client1).Update("name", "Newly Updated Client Name").Error
	if err != nil {
		t.Errorf("Error updating client: %v", err)
	}

	var client2 Client
	tx = db.First(&client2, "67890")
	if tx.Error != nil {
		t.Errorf("Failed to take the updated: %v", tx.Error)
	}

	assert.Equal(t, client.ID, client2.ID)
	assert.Equal(t, "Newly Updated Client Name", client2.Name)
}

func TestClient_GetActivities(t *testing.T) {
	setUpAll(t)
	defer tearDown(t)

	count, err := client.InsertNewActivity(db, 0)
	assert.Equal(t, int64(1), count)
	assert.Nil(t, err)
}

func TestClient_GetCommandExecutions(t *testing.T) {
	setUpAll(t)
	defer tearDown(t)
	count, err := client.InsertNewCommandExecution(db, 1, "0", nil)
	assert.Equal(t, int64(1), count)
	assert.Nil(t, err)
}
