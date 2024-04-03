package common

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vistart/project20240227/server/models"
)

type SessionManager struct {
	Message       SessionChannel     // 向所有目前活跃的会话广播消息。
	NewClients    chan *Client       // 接收新加入的客户端。
	ClosedClients chan *Client       // 接收退出的客户端。
	TotalClients  map[string]*Client // 目前活跃客户端。键为客户端ID。
	mu            sync.RWMutex       // TotalClients 读写锁。
}

func NewSessionManager() (session *SessionManager) {
	session = &SessionManager{
		Message:       make(SessionChannel),
		NewClients:    make(chan *Client),
		ClosedClients: make(chan *Client),
		TotalClients:  make(map[string]*Client),
	}
	return session
}

const GinKeySessionChannel = "session_channel"

// NewSessionChannelHandler 为 gin 的请求准备的实例化会话通道的句柄。
// 1. 实例化一个新的 Client 结构体。该结构体会实例化新的会话通道。
// 2. 该新实例化的 Client 结构体送入 NewClients 通道。
// 3. 当该请求断开时，该新实例化的 Client 结构体送入 ClosedClients 通道。
func (s *SessionManager) NewSessionChannelHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, _ := c.Get("client-id")
		clientType, _ := c.Get("client-type")
		if s.GetClient(clientID.(string)) != nil { // 如果ClientID已存在，代表这个客户端已连接。
			c.AbortWithStatusJSON(http.StatusBadRequest, "client id duplicated")
			return
		}
		client := NewClient(clientID.(string), clientType.(models.ClientType))
		newly, err := client.InsertNewClient()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		if newly > 0 {
			log.Printf("New client[%s] added.", clientID)
		} else {
			log.Printf("The existed client[%s] connected.", clientID)
		}
		client.CreateSessionChannel()
		s.SetClient(clientID.(string), client)
		s.NewClients <- client
		defer func() {
			if s.GetClient(clientID.(string)) == nil { // 如果ClientID不存在，则忽略。
				c.AbortWithStatusJSON(http.StatusInternalServerError, "client not found")
				return
			} else {
				client.CloseSessionChannel()
				s.DeleteClient(clientID.(string))
				s.ClosedClients <- client
			}
		}()
		c.Set("client", client)
		c.Next()
	}
}

// GetClientHandler 为 gin 的控制器动作添加句柄。该句柄用于根据参数client-id拿到指定的Client实例地址。如果指定Client不存在，则报错。
func (s *SessionManager) GetClientHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID, existed := c.Get("client-id")
		if !existed {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "client id not specified")
			return
		}
		client := s.GetClient(clientID.(string))
		if client == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, "client not connected")
			return
		} else {
			c.Set("client", client)
		}
		c.Next()
	}
}

// SetHeadersHandler 为连接设置固定的请求头信息。
func (s *SessionManager) SetHeadersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

// GetClient 获得当前活动会话中某个客户端。如果某个客户端已断开，则返回nil。
func (s *SessionManager) GetClient(id string) *Client {
	if len(id) == 0 {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if client, existed := s.TotalClients[id]; existed {
		return client
	}
	return nil
}

func (s *SessionManager) SetClient(id string, client *Client) {
	if len(id) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalClients[id] = client
}

func (s *SessionManager) DeleteClient(id string) {
	if len(id) == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.TotalClients, id)
}

// GetIsActive 检查客户端是否存在。
func (s *SessionManager) GetIsActive(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, existed := s.TotalClients[id]; existed {
		return true
	}
	return false
}

// Count 检查客户端活跃客户数。
func (s *SessionManager) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.TotalClients)
}

// Serve 提供服务。
// 当有客户端连接或断开时，输出日志并记录到数据库中。
// 当有需要发广播消息时，为每个客户端广播消息。
func (s *SessionManager) Serve() {
	for {
		select {
		case client := <-s.NewClients:
			log.Printf("Client[%s] added. %d registered client(s)", client.ID(), s.Count())
			client.ReceiveActivity(ClientActivityOn)
			channel := client.GetSessionChannel()
			channel <- NewEventMessage("connected")
		case client := <-s.ClosedClients:
			log.Printf("Client[%s] removed. %d registered client(s).", client.ID(), s.Count())
			client.ReceiveActivity(ClientActivityOff)
		case message := <-s.Message:
			for _, client := range s.TotalClients {
				client.GetSessionChannel() <- message
			}
		}
	}
}

// BroadcastTimestamp 广播时间戳。同时也起到发现客户端已断开的作用。间隔时间单位为毫秒。
func (s *SessionManager) BroadcastTimestamp(interval int64) {
	if interval == 0 {
		panic(errors.New("broadcast timestamp interval is zero"))
	}
	ticker := time.NewTicker(time.Millisecond * time.Duration(interval))
	for {
		select {
		case <-ticker.C:
			now := time.Now().String()
			// log.Printf("Broadcast: %s\n", now)
			s.Message <- NewEventMessage(now)
		}
	}
}

var GlobalSessionManager *SessionManager
