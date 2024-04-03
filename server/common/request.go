package common

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/vistart/project20240227/server/models"
	"gorm.io/gorm"
)

const RequestClientID = "x-request-client-id"
const RequestClientType = "x-request-client-type"
const RequestAuthorization = "x-request-authorization"

type RequestClientAuthorization struct {
	ClientID      string
	ClientType    models.ClientType
	Authorization string
}

func NewRequestClientAuthorization(id string, clientType models.ClientType, authorization string) *RequestClientAuthorization {
	return &RequestClientAuthorization{
		ClientID:      id,
		ClientType:    clientType,
		Authorization: authorization,
	}
}

type ErrRequestBadClientID struct {
	ID string
	error
}

func (e ErrRequestBadClientID) Error() string {
	return fmt.Sprintf("bad client id: %s", e.ID)
}

type ErrRequestDuplicateClientID struct {
	ID string
	error
}

func (e ErrRequestDuplicateClientID) Error() string {
	return fmt.Sprintf("duplicate client id: %s", e.ID)
}

type ErrRequestBadClientType struct {
	Type models.ClientType
	error
}

func (e ErrRequestBadClientType) Error() string {
	return fmt.Sprintf("bad client type: %d", e.Type)
}

type ErrRequestAuthFailed struct {
	Authorization string
	error
}

func (e ErrRequestAuthFailed) Error() string {
	return fmt.Sprintf("auth failed: %s", e.Authorization)
}

func (r *RequestClientAuthorization) Auth() error {
	// 定义字符串
	charRange := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 检查字符串中的字符是否在字符范围内
	for _, c := range r.ClientID {
		if !strings.ContainsRune(charRange, c) {
			return ErrRequestBadClientID{
				ID: r.ClientID,
			}
		}
	}
	client, err := models.GetClient(DB, r.ClientID)
	if !errors.Is(err, gorm.ErrRecordNotFound) && err != nil {
		return err
	}
	if client != nil && client.Type != r.ClientType {
		return ErrRequestBadClientType{Type: r.ClientType}
	}
	clientID := []byte(r.ClientID)
	sum := md5.Sum(append(clientID, byte(r.ClientType)))
	log.Println(hex.EncodeToString(sum[:]))
	if hex.EncodeToString(sum[:]) != r.Authorization {
		return ErrRequestAuthFailed{
			Authorization: r.Authorization,
		}
	}
	return nil
}
