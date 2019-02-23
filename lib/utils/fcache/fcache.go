package fcache

import (
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/logs"
)

var cm cache.Cache//, err := cache.NewCache("memory", `{"interval":60}`)
var inited bool

func Prepare() (err error) {
	if !inited {
		cm, err = cache.NewCache("file", `{"CachePath":"./cache","FileSuffix":".cache","DirectoryLevel":"2","EmbedExpiry":"120"}`)
		if err != nil {
			logs.Error("cm")
			return
		}
		inited = true
	}

	return
}

func Get(key string) interface{} {
	Prepare()
	if cm == nil {
		return nil
	}

	return cm.Get(key)
}

func Set(key string, val interface{}) {
	if cm != nil {
		cm.Put(key, val, 0)
	}
}

func GetInt(key string) int64 {
	if cm != nil {
		v := cm.Get(key)
		if vi, ok := v.(int64); ok {
			return vi
		}
	}
	return 0
}

func SetInt(key string, val int64) {
	if cm != nil {
		cm.Put(key, val, 0)
	}
}

func GetString(key string) string {
	if cm != nil {
		v := cm.Get(key)
		if vi, ok := v.(string); ok {
			return vi
		}
	}
	return ""
}

func SetString(key string, val string) {
	if cm != nil {
		cm.Put(key, val, 0)
	}
}
