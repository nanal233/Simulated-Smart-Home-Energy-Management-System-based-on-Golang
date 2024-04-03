package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPowerMode(t *testing.T) {
	setUpAll(t)
	defer tearDown(t)

	// 测试 CreateNewPowerMode
	result, err := CreateNewPowerMode(db, "test-power-mode-get")
	assert.Equal(t, int64(1), result)
	assert.Nil(t, err)

	// 测试 GetPowerModes
	modes, count, err := GetPowerModes(db, 1, 1)
	assert.Len(t, modes, 1)
	assert.Equal(t, int64(1), count)
	assert.Nil(t, err)

	// 测试 GetPowerMode
	mode := GetPowerMode(db, modes[0].ID)
	assert.NotNil(t, mode)
	assert.Equal(t, modes[0].ID, mode.ID)
	assert.Equal(t, modes[0].Name, mode.Name)
	assert.Equal(t, modes[0].CreatedAt, mode.CreatedAt)
	assert.Equal(t, modes[0].UpdatedAt, mode.UpdatedAt)

	// 测试 UpdateName
	result, err = mode.UpdateName(db, "test-power-mode-get-updated")
	assert.Equal(t, int64(1), result)
	assert.Nil(t, err)

	modeUpdated := GetPowerMode(db, modes[0].ID)
	assert.NotNil(t, modeUpdated)
	assert.Equal(t, mode.ID, modeUpdated.ID)
	assert.Equal(t, mode.Name, modeUpdated.Name)
	assert.NotEqual(t, modes[0].Name, modeUpdated.Name)
	assert.Equal(t, mode.UpdatedAt, modeUpdated.UpdatedAt)
	assert.NotEqual(t, modes[0].UpdatedAt, modeUpdated.UpdatedAt)

	// 测试 RemovePowerMode
	result, err = RemovePowerMode(db, mode)
	assert.Equal(t, int64(1), result)
	assert.Nil(t, err)

	modeDeleted := GetPowerMode(db, modes[0].ID)
	assert.Nil(t, modeDeleted)
}

// TestGetPowerModeByName 测试通过名称获取能耗模式。
func TestGetPowerModeByName(t *testing.T) {
	setUpAll(t)
	defer tearDown(t)

	// 测试 CreateNewPowerMode
	result, err := CreateNewPowerMode(db, "test-power-mode-get-by-name")
	assert.Equal(t, int64(1), result)
	assert.Nil(t, err)

	// 测试 GetPowerModeByName
	mode := GetPowerModeByName(db, "test-power-mode-get-by-name")
	assert.NotNil(t, mode)
	assert.Equal(t, "test-power-mode-get-by-name", mode.Name)

	result, err = RemovePowerMode(db, mode)
	assert.Equal(t, int64(1), result)
	assert.Nil(t, err)

	modeDeleted := GetPowerModeByName(db, "test-power-mode-get-by-name")
	assert.Nil(t, modeDeleted)
}

// TestPowerModeAndClientPreparedCommand 测试能耗模式，及添加、删除预制命令。
func TestPowerModeAndClientPreparedCommand(t *testing.T) {
	setUpAll(t)
	defer tearDown(t)

	result, err := CreateNewPowerMode(db, "test-power-mode-1")
	assert.Equal(t, int64(1), result)
	assert.Nil(t, err)

	mode := GetPowerModeByName(db, "test-power-mode-1")
	assert.NotNil(t, mode)
	assert.Equal(t, "test-power-mode-1", mode.Name)

	// 注册新客户端
	client := NewClient("id-client-power-mode-1", "name-client-power-mode-1", 1)
	result, err = RegisterNewClient(db, client)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), result)

	// 注册后添加预制命令。
	prepared := NewClientPreparedCommand(client, mode, 1, "data")
	result, err = CreateNewClientPreparedCommand(db, prepared)
	assert.Equal(t, int64(1), result)
	assert.Nil(t, err)

	//// 加入前，获取命令数，应当为0。
	//commands, total, err := mode.GetCommands(db)
	//assert.Nil(t, err)
	//assert.Equal(t, int64(0), total)
	//assert.Len(t, commands, 0)
	//
	//// 准备好预制命令后，加入到能耗模式中。
	//result, err = mode.AddCommand(db, prepared)
	//assert.Equal(t, int64(1), result)
	//assert.Nil(t, err)
	//
	//// 再次加入会报错
	//result, err = mode.AddCommand(db, prepared)
	//assert.Equal(t, int64(0), result)
	//assert.NotNil(t, err)
	//log.Print(err) // 报 duplicate entry
	//
	//// 获取命令
	//_, total, err = mode.GetCommands(db)
	//assert.Nil(t, err)
	//assert.Equal(t, int64(1), total)
}
