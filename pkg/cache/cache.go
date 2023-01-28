package cache

import (
	"errors"
	"fmt"
	"github.com/danymarita/redis-otp-sample/pkg/config"
	"github.com/gomodule/redigo/redis"
	"github.com/im7mortal/kmutex"
	"github.com/spf13/cast"
	"time"
)

type ICache interface {
	CheckCacheExists(key string) bool
	ReadCache(key string) (data []byte, err error)
	WriteCache(key string, data []byte, ttl time.Duration) (err error)
	WriteCacheIfEmpty(key string, data []byte, ttl time.Duration) (err error)
	DeleteCache(key string) (err error)
	IncrementWithTtlCache(key string, ttl time.Duration) (incr int64, err error)
	IncrementCache(key string) (incr int64, err error)
}

type cache struct {
	pool   *redis.Pool
	kmutex *kmutex.Kmutex
}

// NewCacheRepository initiate cache repo
func NewCache(cfg config.ConfigObject) ICache {
	dialConnectTimeoutOption := redis.DialConnectTimeout(cast.ToDuration(cfg.RedisDialConnectTimeout))
	readTimeoutOption := redis.DialReadTimeout(cast.ToDuration(cfg.RedisReadTimeout))
	writeTimeoutOption := redis.DialWriteTimeout(cast.ToDuration(cfg.RedisWriteTimeout))

	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.RedisHost, cast.ToInt(cfg.RedisPort)), dialConnectTimeoutOption, readTimeoutOption, writeTimeoutOption)
			if err != nil {
				return nil, fmt.Errorf("ERROR connect redis | %v", err)
			}
			if cfg.RedisPassword != "" {
				if _, err := c.Do("AUTH", cfg.RedisPassword); err != nil {
					return nil, fmt.Errorf("ERROR connect redis | %v", err)
				}
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				return err
			}
			return nil
		},
		MaxIdle:         cast.ToInt(cfg.RedisConnIdleMax),
		MaxActive:       cast.ToInt(cfg.RedisConnActiveMax),
		IdleTimeout:     cast.ToDuration(cfg.RedisIdleTimeout),
		Wait:            cast.ToBool(cfg.RedisIsWait),
		MaxConnLifetime: cast.ToDuration(cfg.RedisConnLifetimeMax),
	}

	return &cache{
		pool:   pool,
		kmutex: kmutex.New(),
	}
}

func (c *cache) IncrementWithTtlCache(key string, ttl time.Duration) (incr int64, err error) {
	conn := c.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("INCR", key)
	conn.Send("EXPIRE", key, ttl.Seconds())

	res, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return
	}
	incr = res[0].(int64)
	return
}

func (c *cache) IncrementCache(key string) (incr int64, err error) {
	c.kmutex.Lock(key)
	defer c.kmutex.Unlock(key)

	conn := c.pool.Get()
	defer conn.Close()

	res, err := conn.Do("INCR", key)
	if err != nil {
		return
	}
	incr = res.(int64)
	return
}

func (c *cache) CheckCacheExists(key string) bool {
	// check whether cache value is empty
	conn := c.pool.Get()
	defer conn.Close()

	exists, _ := redis.Bool(conn.Do("EXISTS", key))

	return exists
}

func (c *cache) ReadCache(key string) (data []byte, err error) {
	c.kmutex.Lock(key)
	defer c.kmutex.Unlock(key)

	// check whether cache value is empty
	conn := c.pool.Get()
	defer conn.Close()

	exists := c.CheckCacheExists(key)
	if exists {
		data, err = redis.Bytes(conn.Do("GET", key))
		return
	}
	return nil, errors.New(fmt.Sprintf("Cache key didn't exists. Key : %s", key))
}

// WriteCache this will and must write the data to cache with corresponding key using locking
func (c *cache) WriteCache(key string, data []byte, ttl time.Duration) (err error) {
	c.kmutex.Lock(key)
	defer c.kmutex.Unlock(key)

	// write data to cache
	conn := c.pool.Get()
	defer conn.Close()

	_, err = conn.Do("SETEX", key, ttl.Seconds(), data)

	return
}

// WriteCacheIfEmpty will try to write to cache, if the data still empty after locking
func (c *cache) WriteCacheIfEmpty(key string, data []byte, ttl time.Duration) (err error) {
	c.kmutex.Lock(key)
	defer c.kmutex.Unlock(key)

	// check whether cache value is empty
	conn := c.pool.Get()
	defer conn.Close()

	_, err = conn.Do("GET", key)
	if err != nil {
		if err == redis.ErrNil {
			return nil //return nil as the data already set, no need to overwrite
		}

		return err
	}

	// write data to cache
	_, err = conn.Do("SETEX", key, ttl.Seconds(), data)
	if err != nil {
		return err
	}

	return nil
}

func (c *cache) DeleteCache(key string) (err error) {
	c.kmutex.Lock(key)
	defer c.kmutex.Unlock(key)

	// write data to cache
	conn := c.pool.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", key)

	return
}
