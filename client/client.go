package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/vistart/project20240227/server/common"
)

type ClientInterface interface {
	ID() string
	Type() int
	Register()
	Report()
	SetHeader()
	ProcessEvent()
}

type Client struct {
	id         string
	clientType int
	PowerMode  *PowerMode
	ClientInterface
}

func NewClient(id string, clientType int, powerFactor int) *Client {
	return &Client{
		id:         id,
		clientType: clientType,
		PowerMode:  NewPowerMode(powerFactor),
	}
}

func (c *Client) ID() string {
	return c.id
}

func (c *Client) Type() int {
	return c.clientType
}
func (c *Client) Register(exitFunc func()) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/client/register", apiSocket), nil)
	if err != nil {
		log.Println(err)
		return
	}

	c.SetHeader(req)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		// 处理错误
		log.Println(err)
		return
	}

	// 读取响应正文
	body := resp.Body
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Println(err)
		}
	}(body)

	defer exitFunc()

	scanner := bufio.NewScanner(resp.Body)
	var event, data string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event:") {
			// 提取事件类型并去掉 "event:" 前缀
			event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			// 提取数据行并去掉 "data:" 前缀
			data = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		} else if line == "" {
			// 空行表示事件结束，进行下一步处理
			if event != "" || data != "" {
				//fmt.Println("Event:", event)
				//fmt.Println("Data:", data)
				e := c.ProcessEvent(event, data)
				fmt.Println(event, e.MarshalData())
				if _, ok := e.(*common.EventDisconnect); ok {
					return
				}
			}
			// 重置事件和数据，准备接收下一个事件
			event, data = "", ""
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading SSE:", err)
	}
}

func (c *Client) Report() {
	client := &http.Client{}

	postData := url.Values{}
	postData.Set("consumption", strconv.Itoa(c.PowerMode.GetConsumption()))
	postData.Set("recorded_at", fmt.Sprintf("%d", time.Now().Unix()))

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/client/report", apiSocket), strings.NewReader(postData.Encode()))
	if err != nil {
		log.Println(err)
		return
	}

	c.SetHeader(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		// 处理错误
		log.Println(err)
		return
	}
	defer resp.Body.Close()
}

func (c *Client) SetHeader(req *http.Request) {
	clientID := c.ID()
	clientType := fmt.Sprintf("%d", c.Type())
	req.Header.Set(common.RequestClientID, clientID)
	req.Header.Set(common.RequestClientType, clientType)
	clientIDBytes := []byte(clientID)
	sum := md5.Sum(append(clientIDBytes, byte(c.Type())))
	req.Header.Set(common.RequestAuthorization, hex.EncodeToString(sum[:]))
}

func (c *Client) ProcessEvent(event, data string) common.EventInterface {
	switch event {
	case common.EventNameRegistration:
		e := &common.EventRegistration{}
		err := e.EventBase.UnmarshalData(data)
		if err != nil {
			return nil
		}
		return e
	case common.EventNameCommandPower:
		e := &common.EventCommandPower{}
		err := e.EventBase.UnmarshalData(data)
		if err != nil {
			return nil
		}
		c.PowerMode.Change(e.EventBase.Data.Power)
		return e
	case common.EventNameMessage:
		e := &common.EventMessage{}
		err := e.EventBase.UnmarshalData(data)
		if err != nil {
			return nil
		}
		return e
	case common.EventNameDisconnect:
		e := &common.EventDisconnect{}
		return e
	default:
		return &common.EventNone{}
	}
}
