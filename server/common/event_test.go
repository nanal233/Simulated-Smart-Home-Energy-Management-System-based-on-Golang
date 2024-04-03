package common

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewEventBase 测试实例化基础命令事件。
func TestNewEventBase(t *testing.T) {
	event := NewEventBase[int](0, 0)
	content, err := json.Marshal(event)
	log.Println(string(content), err)
}

// TestNewEventRegistrationSucceeded 测试实例化命令注册成功。
func TestNewEventRegistrationSucceeded(t *testing.T) {
	event := NewEventRegistrationSucceeded()
	content, err := json.Marshal(event)
	assert.Equal(t, "{\"code\":1,\"data\":0}", string(content))
	assert.Nil(t, err)
}

// TestNewEventCommand 测试实例化命令事件。
func TestNewEventCommand(t *testing.T) {
	event := NewEventCommand(EventCommandPowerData{0})
	content, err := json.Marshal(event)
	assert.Equal(t, "{\"code\":2,\"data\":{\"power\":0}}", string(content))
	assert.Nil(t, err)
}

// TestEvent_Unmarshal 测试事件反序列化。
func TestEvent_Unmarshal(t *testing.T) {
	event := NewEventCommand(EventCommandPowerData{0})
	content, err := json.Marshal(event)
	assert.Equal(t, "{\"code\":2,\"data\":{\"power\":0}}", string(content))
	assert.Nil(t, err)

	event1 := NewEventCommand(EventCommandPowerData{1})
	assert.Equal(t, EventCommandPowerData(EventCommandPowerData{Power: 1}), event1.Data)
	err = json.Unmarshal(content, &event1)
	assert.Nil(t, err)
	assert.Equal(t, EventCommandPowerData(EventCommandPowerData{Power: 0}), event1.Data)
}
