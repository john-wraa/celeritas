package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Has(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Forget(string) error
	EmptyByMatch(string) error
	Empty() error
}

type RedisCache struct {
	Conn   *redis.Pool
	Prefix string
}

type Entry map[string]interface{}

func (c *RedisCache) Has(key string) (bool, error) {
	innerKey := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	ok, err := redis.Bool(conn.Do("EXISTS", innerKey))
	if err != nil {
		return false, err
	}
	return ok, nil
}

func encode(item Entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decode(str string) (Entry, error) {
	item := Entry{}
	b := bytes.Buffer{}
	b.Write([]byte(str))
	d := gob.NewDecoder(&b)
	err := d.Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (c *RedisCache) Get(key string) (interface{}, error) {
	innerKey := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	cacheEntry, err := redis.Bytes(conn.Do("GET", innerKey))
	if err != nil {
		return nil, err
	}

	decoded, err := decode(string(cacheEntry))
	if err != nil {
		return nil, err
	}

	item := decoded[innerKey]

	return item, nil
}

func (c *RedisCache) Set(key string, value interface{}, expires ...int) error {
	innerKey := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	entry := Entry{}
	entry[innerKey] = value

	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		//goland:noinspection SpellCheckingInspection
		_, err := conn.Do("SETEX", innerKey, expires[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("SET", innerKey, string(encoded))
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Forget(key string) error {
	innerKey := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	_, err := conn.Do("DEL", innerKey)
	if err != nil {
		return err
	}

	return nil
}

func (c *RedisCache) EmptyByMatch(key string) error {
	innerKey := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	keys, err := c.getKeys(innerKey)
	if err != nil {
		return err
	}

	for _, k := range keys {
		_, err := conn.Do("DEL", k)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) Empty() error {
	innerKey := fmt.Sprintf("%s:", c.Prefix)
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	keys, err := c.getKeys(innerKey)
	if err != nil {
		return err
	}

	for _, k := range keys {
		_, err := conn.Do("DEL", k)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *RedisCache) getKeys(pattern string) ([]string, error) {
	conn := c.Conn.Get()
	defer func(conn redis.Conn) {
		_ = conn.Close()
	}(conn)

	iter := 0
	//keys := []string{}
	var keys []string

	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", fmt.Sprintf("%s*", pattern)))
		if err != nil {
			return keys, err
		}

		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}
	}

	return keys, nil
}
