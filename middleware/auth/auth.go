package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type AuthInterface interface {
	CreateAuth(uint, *Tokens) error
	FetchUuid(string) (string, error)
	DeleteRefresh(string) error
	DeleteTokens(*TokenMetadata) error
}

type RedisAuthService struct {
	client *redis.Client
}

var _ AuthInterface = &RedisAuthService{}

func NewAuthService(client *redis.Client) *RedisAuthService {
	return &RedisAuthService{client: client}
}

//Save token metadata to Redis
func (tk *RedisAuthService) CreateAuth(userId uint, td *Tokens) error {
	at := time.Unix(td.AtExpires, 0)
	wst := time.Unix(td.WsExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(td.AccessUuid, userId, at.Sub(now)).Result()
	if err != nil {
		return err
	}

	wstCreated, err := tk.client.Set(td.WsUuid, userId, wst.Sub(now)).Result()
	if err != nil {
		return err
	}

	rtCreated, err := tk.client.Set(td.RefreshUuid, userId, rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" || wstCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

//Check the metadata saved
func (tk *RedisAuthService) FetchUuid(uuid string) (string, error) {
	userid, err := tk.client.Get(uuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}

// TODO list current access for the account
// func (tk *RedisAuthService) FetchAuthById(userId uint) (string, error) {
// 	accessUuid, err := tk.client.Get("++"+fmt.Sprint(userId)).Result() // check in the redis api
// 	if err != nil {
// 		return "", err
// 	}
// 	return accessUuid, nil
// }

func (tk *RedisAuthService) DeleteTokens(authD *TokenMetadata) error {
	//get uuids
	refreshUuid := fmt.Sprintf("%s++%d", authD.Uuid, authD.UserId)
	wstUuid := fmt.Sprintf("%s++WS", authD.Uuid)
	//delete access token
	deletedAt, err := tk.client.Del(authD.Uuid).Result()
	if err != nil {
		return err
	}
	//delete websocket token
	deletedWst, err := tk.client.Del(wstUuid).Result()
	if err != nil {
		return err
	}
	//delete refresh token
	deletedRt, err := tk.client.Del(refreshUuid).Result()
	if err != nil {
		return err
	}
	//When the record is deleted, the return value is 1
	if deletedAt != 1 || deletedRt != 1 || deletedWst != 1 {
		return errors.New("something went wrong")
	}
	return nil
}

func (tk *RedisAuthService) DeleteRefresh(refreshUuid string) error {
	//delete refresh token
	deleted, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}
	return nil
}
