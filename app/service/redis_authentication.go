package service

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type RedisAuthentication struct {
	client *redis.Client
}

func (i *RedisAuthentication) Insert(auth *Authentication) error {
	if err := i.Update(auth); err != nil {
		return err
	}
	cmd := i.client.LPush("AUTHLIST", auth.Username)
	return cmd.Err()
}

func (i *RedisAuthentication) Update(auth *Authentication) error {
	authKey := i.constructAuthKey(auth.Username)
	marshalled, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	cmd := i.client.Set(authKey, string(marshalled), 0)
	if cmd.Err() != nil {
		return err
	}
	sessionKey := i.constructSessionKey(auth.SessionID)
	cmd = i.client.Set(sessionKey, auth.Username, 0)
	if cmd.Err() != nil {
		return err
	}
	return nil
}

func (i *RedisAuthentication) FindBySessionID(sessionID string) (*Authentication, error) {
	sessionKey := i.constructSessionKey(sessionID)
	cmd := i.client.Get(sessionKey)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return nil, nil
		}
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return i.FindByUsername(res)
}

func (i *RedisAuthentication) FindByUsername(username string) (*Authentication, error) {
	authKey := i.constructAuthKey(username)
	cmd := i.client.Get(authKey)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return nil, nil
		}
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	auth := new(Authentication)
	if err := json.Unmarshal([]byte(res), &auth); err != nil {
		return nil, err
	}
	return auth, nil
}

func (i *RedisAuthentication) List(limit uint, offset uint) ([]*Authentication, error) {
	cmd := i.client.LRange("AUTHLIST", int64(offset), int64(offset+limit))
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	authentications := make([]*Authentication, 0, len(res))
	for _, val := range res {
		auth, err := i.FindByUsername(val)
		if err != nil {
			return nil, err
		}
		authentications = append(authentications, auth)
	}
	return authentications, nil
}

func (i *RedisAuthentication) constructAuthKey(username string) string {
	return fmt.Sprintf("AUTH:%s", username)
}

func (i *RedisAuthentication) constructSessionKey(sessionID string) string {
	return fmt.Sprintf("SESSION:%s", sessionID)
}

func NewRedisAuthentication(option *redis.Options) *RedisAuthentication {
	instance := new(RedisAuthentication)
	instance.client = redis.NewClient(option)
	return instance
}
