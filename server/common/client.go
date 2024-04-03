package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/vistart/project20240227/server/models"
	"gorm.io/gorm"
)

const (
	ClientTypeNone = iota
	ClientType1
	ClientType2
	ClientType3
	ClientType4
	ClientType5
	ClientType6
	ClientType7
)

var ClientTypeNames = map[int]string{
	ClientTypeNone: "无",
	ClientType1:    "开关",
}

// ClientBaseInterface 表示一个客户端应该具备的方法。
type ClientBaseInterface interface {
	// ID 表示设备唯一编号。若设备注册时已存在ID，则无法注册。
	ID() string

	// Type 表示设备类型。
	Type() models.ClientType

	InsertNewClient() (int64, error)
}

type ClientSessionInterface interface {
	// SendToSessionChannel 向通道发送内容。需要发送的内容需要实现 encoding/json 的 Marshaller 接口，即 MarshalJSON() 方法。
	SendToSessionChannel(v *EventBase[any])
	CreateSessionChannel()
	CloseSessionChannel()
	GetSessionChannel() SessionChannel
}

type ClientActivityInterface interface {
	ReceiveActivity(int8)
}

type ClientConsumptionInterface interface {
	ReceiveReportConsumption(float32, time.Time) (int64, error)
}

// ClientBase 保存了ID和Type，但不能直接访问。
type ClientBase struct {
	id         string
	clientType models.ClientType
	ClientBaseInterface

	sessionChannel SessionChannel
	ClientSessionInterface

	ClientActivityInterface
	ClientConsumptionInterface
}

func (c *ClientBase) ID() string {
	return c.id
}

func (c *ClientBase) Type() models.ClientType {
	return c.clientType
}

func (c *ClientBase) InsertNewClient() (int64, error) {
	_, err := models.GetClient(DB, c.ID())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		model := models.NewClient(c.ID(), c.ID(), c.Type())
		return models.RegisterNewClient(DB, model)
	}
	return 0, nil
}

// SendToSessionChannel 向客户端通道发送内容。
// 发送的内容目前是任意类型，但目前仅支持
func (c *ClientBase) SendToSessionChannel(v any) {
	c.sessionChannel <- v
}

func (c *ClientBase) CreateSessionChannel() {
	c.sessionChannel = make(SessionChannel)
}

func (c *ClientBase) CloseSessionChannel() {
	close(c.sessionChannel)
}

func (c *ClientBase) GetSessionChannel() SessionChannel {
	return c.sessionChannel
}

func (c *ClientBase) ReceiveActivity(content int8) {
	client, err := models.GetClient(DB, c.ID())
	if err != nil {
		return
	}
	client.InsertNewActivity(DB, content)
}

func (c *ClientBase) ReceiveReportConsumption(consumption float32, recordedAt time.Time) (int64, error) {
	client, err := models.GetClient(DB, c.ID())
	if err != nil {
		return 0, nil
	}
	return client.InsertNewConsumption(DB, consumption, recordedAt)
}

// Client 客户端。
type Client struct {
	ClientBase
}

// NewClient 实例化一个新的客户端。
func NewClient(id string, clientType models.ClientType) *Client {
	client := &Client{
		ClientBase{id: id, clientType: clientType},
	}
	return client
}

type ErrClientNotFound struct {
	id string
	error
}

func (e ErrClientNotFound) Error() string {
	return fmt.Sprintf("client `%s` not found", e.id)
}

func NewErrClientNotFound(id string) ErrClientNotFound {
	return ErrClientNotFound{id: id}
}
