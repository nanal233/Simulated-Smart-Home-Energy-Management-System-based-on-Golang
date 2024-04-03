package common

import (
	"encoding/json"
	"sync"
)

// EventInterface 表示一个事件应该实现的方法。
type EventInterface interface {
	// UnmarshalData 将data解析到自身的data字段。
	UnmarshalData(data string) error
	// MarshalData 将data字段序列化。
	MarshalData() string
}

// EventBase 表示一个事件应该具备的字段。
type EventBase[T any] struct {
	Code           int `json:"code"` // 事件代码。
	Data           T   `json:"data"` // 事件内容。
	mu             sync.RWMutex
	EventInterface `json:"-"` // 不参与序列化与反序列化
}

// MarshalData 表示单独将Data序列化。
func (e *EventBase[T]) MarshalData() string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	d, _ := json.Marshal(e.Data)
	return string(d)
}

func (e *EventBase[T]) UnmarshalData(data string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	return json.Unmarshal([]byte(data), &e.Data)
}

func NewEventBase[T any](code int, data T) *EventBase[T] {
	return &EventBase[T]{
		Code: code,
		Data: data,
	}
}

const (
	EventCodeNone         = iota // 无事件。此时data为nil。
	EventCodeRegistration        // 注册结果。
	EventCodeCommandPower        // 命令。
	EventCodeMessage             // 消息。
	EventCodeDisconnect
)

const (
	EventNameNone         = ""
	EventNameRegistration = "registration"
	EventNameCommandPower = "command-power"
	EventNameMessage      = "message"
	EventNameDisconnect   = "disconnect"
)

var EventCodeNameMap = map[int]string{
	EventCodeNone:         EventNameNone,
	EventCodeRegistration: EventNameRegistration,
	EventCodeCommandPower: EventNameCommandPower,
	EventCodeMessage:      EventNameMessage,
	EventCodeDisconnect:   EventNameDisconnect,
}

const (
	EventRegistrationSucceeded                    = iota         // 注册成功
	EventRegistrationFailedIdAuthNotMatch         = iota + 10000 // 注册失败：ID验证不匹配
	EventRegistrationFailedIdConnectionDuplicated                // 注册失败：ID连接重复
)

type EventNone struct {
	EventBase[any]
}

type EventRegistration struct {
	EventBase[int]
}

func NewEventRegistrationSucceeded() *EventBase[int] {
	return &EventBase[int]{
		Code: EventCodeRegistration,
		Data: EventRegistrationSucceeded,
	}
}

func NewEventRegistrationFailed(data int) *EventBase[int] {
	return &EventBase[int]{
		Code: EventCodeRegistration,
		Data: data,
	}
}

type EventCommandPowerData struct {
	Power int `json:"power"`
}

type EventCommandPower struct {
	EventBase[EventCommandPowerData]
}

type EventMessageData struct {
	Message string `json:"message"`
}

type EventMessage struct {
	EventBase[EventMessageData]
}

type EventDisconnect struct {
	EventBase[struct{}]
}

func NewEventCommand(data EventCommandPowerData) *EventBase[EventCommandPowerData] {
	return &EventBase[EventCommandPowerData]{
		Code: EventCodeCommandPower,
		Data: data,
	}
}

// NewEventCommandPower 实例化一个调整功率命令事件。
func NewEventCommandPower(power int) *EventBase[EventCommandPowerData] {
	return NewEventCommand(EventCommandPowerData{power})
}

const (
	// ClientActivityOff 表示设备不活跃。
	ClientActivityOff = iota
	// ClientActivityOn 表示设备活跃。
	ClientActivityOn
)

func NewEventMessage(data string) *EventBase[EventMessageData] {
	return &EventBase[EventMessageData]{
		Code: EventCodeMessage,
		Data: struct {
			Message string `json:"message"`
		}{Message: data},
	}
}

func NewEventDisconnect() *EventBase[struct{}] {
	return &EventBase[struct{}]{
		Code: EventCodeDisconnect,
		Data: struct{}{},
	}
}
