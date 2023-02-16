package cache

import (
	"github.com/topjohncian/cloudreve-pro-epay/internal/appconf"
	"go.uber.org/fx"
)

func Cache() fx.Option {
	return fx.Module("cache", fx.Provide(func(conf *appconf.Config) Driver {
		if conf.RedisEnabled {
			return NewRedisStore(10, "tcp", conf.RedisServer, conf.RedisPassword, conf.RedisDB)
		} else {
			return NewMemoStore()
		}
	}))
}

// Driver 键值缓存存储容器
type Driver interface {
	// 设置值，ttl为过期时间，单位为秒
	Set(key string, value interface{}, ttl int) error

	// 取值，并返回是否成功
	Get(key string) (interface{}, bool)

	// 批量取值，返回成功取值的map即不存在的值
	Gets(keys []string, prefix string) (map[string]interface{}, []string)

	// 批量设置值，所有的key都会加上prefix前缀
	Sets(values map[string]interface{}, prefix string) error

	// 删除值
	Delete(keys []string, prefix string) error
}

// // Set 设置缓存值
// func Set(key string, value interface{}, ttl int) error {
// 	return Store.Set(key, value, ttl)
// }

// // Get 获取缓存值
// func Get(key string) (interface{}, bool) {
// 	return Store.Get(key)
// }

// // Deletes 删除值
// func Deletes(keys []string, prefix string) error {
// 	return Store.Delete(keys, prefix)
// }
