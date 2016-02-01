package cache

import (
	rc "github.com/danjac/podbaby/cache/Godeps/_workspace/src/gopkg.in/go-redis/cache.v1"
	"github.com/danjac/podbaby/cache/Godeps/_workspace/src/gopkg.in/redis.v3"
	"github.com/danjac/podbaby/cache/Godeps/_workspace/src/gopkg.in/vmihailenco/msgpack.v2"
	"github.com/danjac/podbaby/config"
	"time"
)

// Setter should generate a cacheable value if cache empty
type Setter func() error

type Cache interface {
	Get(string, time.Duration, interface{}, Setter) error
	Delete(string) error
}

func New(cfg *config.Config) Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": ":6379",
		},

		DialTimeout:  3 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})

	codec := &rc.Codec{
		Ring: ring,

		Marshal: func(v interface{}) ([]byte, error) {
			return msgpack.Marshal(v)
		},
		Unmarshal: func(b []byte, v interface{}) error {
			return msgpack.Unmarshal(b, v)
		},
	}

	return &defaultCache{codec}
}

type defaultCache struct {
	codec *rc.Codec
}

func (c *defaultCache) Delete(key string) error {
	return c.codec.Delete(key)
}

func (c *defaultCache) Get(key string, timeout time.Duration, obj interface{}, fn Setter) error {

	err := c.codec.Get(key, &obj)
	if err == nil {
		return nil
	}

	if err != rc.ErrCacheMiss {
		return err
	}

	if err := fn(); err != nil {
		return err
	}

	return c.codec.Set(&rc.Item{
		Key:        key,
		Object:     obj,
		Expiration: timeout,
	})

}
