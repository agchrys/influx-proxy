// Copyright 2016 Eleme. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package backend

import (
	"errors"
	"log"
)

const (
	VERSION = "1.1"
)

var (
	ErrIllegalConfig = errors.New("illegal config")
)

//func LoadStructFromMap(data map[string]string, o interface{}) (err error) {
//	var x int
//	val := reflect.ValueOf(o).Elem()
//	for i := 0; i < val.NumField(); i++ {
//		valueField := val.Field(i)
//		typeField := val.Type().Field(i)
//
//		name := strings.ToLower(typeField.Name)
//		s, ok := data[name]
//		if !ok {
//			continue
//		}
//
//		switch typeField.Type.Kind() {
//		case reflect.String:
//			valueField.SetString(s)
//		case reflect.Int:
//			x, err = strconv.Atoi(s)
//			if err != nil {
//				log.Printf("%s: %s", err, name)
//				return
//			}
//			valueField.SetInt(int64(x))
//		}
//	}
//	return
//}

type NodeConfig struct {
	ListenAddr   string `json:"listenaddr"`
	DB           string `json:"db"`
	Zone         string `json:"zone"`
	Nexts        string `json:"nextts"`
	Interval     int    `json:"interval"`
	IdleTimeout  int    `json:"idletimeout"`
	WriteTracing int    `json:"writetracing"`
	QueryTracing int    `json:"querytracing"`
}

type BackendConfig struct {
	URL             string `json:url`
	DB              string `json:"db"`
	Zone            string `json:"zone"`
	Interval        int    `json:"interval"`
	Timeout         int    `json:"timeout"`
	TimeoutQuery    int    `json:"timeoutquery"`
	MaxRowLimit     int    `json:"maxrowlimit"`
	CheckInterval   int    `json:"checkinterval"`
	RewriteInterval int    `json:"rewriteinterval"`
	WriteOnly       int    `json:"writeonly"`
}

//agchrys
type GlobelConfig struct {
	Node      NodeConfig               `json:"node"`
	Backs     map[string]BackendConfig `json:"backends"`
	Keymaps   map[string][]string      `json:"keymaps"`
	KeyIgnore []string                 `json:keyignore`
}

type JsonConfigSource struct {
	c *GlobelConfig
}

func (jcs *JsonConfigSource) LoadBackends() (backends map[string]*BackendConfig, err error) {
	backends = make(map[string]*BackendConfig)

	var cfg *BackendConfig
	for name, back := range jcs.c.Backs {
		cfg, err = jcs.LoadConfigFromJson(&back)
		if err != nil {
			log.Printf("read json config error: %s", err)
			return
		}
		backends[name] = cfg
	}
	log.Printf("%d backends loaded from cfg.json", len(backends))
	return
}

func (jcs *JsonConfigSource) LoadConfigFromJson(back *BackendConfig) (cfg *BackendConfig, err error) {
	cfg = &BackendConfig{
		URL:back.URL,
		DB:back.DB,
		Zone:back.Zone,
		Interval:back.Interval,
		Timeout:back.Timeout,
		TimeoutQuery:back.TimeoutQuery,
		MaxRowLimit:back.MaxRowLimit,
		CheckInterval:back.CheckInterval,
		RewriteInterval:back.RewriteInterval,
		WriteOnly:back.WriteOnly,
	}


	if cfg.Interval == 0 {
		cfg.Interval = 1000
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10000
	}
	if cfg.TimeoutQuery == 0 {
		cfg.TimeoutQuery = 600000
	}
	if cfg.MaxRowLimit == 0 {
		cfg.MaxRowLimit = 10000
	}
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = 1000
	}
	if cfg.RewriteInterval == 0 {
		cfg.RewriteInterval = 10000
	}
	return
}

func (jcs *JsonConfigSource) LoadMeasurements() (m_map map[string][]string, err error) {
	m_map = make(map[string][]string, 0)

	for key, val := range jcs.c.Keymaps {
		m_map[key] = val
	}
	log.Printf("%d measurements loaded from cfg.json", len(m_map))
	return
}

//type RedisConfigSource struct {
//	client *redis.Client
//	node   string
//	zone   string
//}

func NewJsonConfigSource(c *GlobelConfig) (jcs *JsonConfigSource) {
	jcs = &JsonConfigSource{
		c: c,
	}
	return
}

//func NewRedisConfigSource(options *redis.Options, node string) (rcs *RedisConfigSource) {
//	rcs = &RedisConfigSource{
//		client: redis.NewClient(options),
//		node:   node,
//	}
//	return
//}

//func (rcs *RedisConfigSource) LoadNode() (nodecfg NodeConfig, err error) {
//	val, err := rcs.client.HGetAll("default_node").Result()
//	if err != nil {
//		log.Printf("redis load error: b:%s", rcs.node)
//		return
//	}
//
//	err = LoadStructFromMap(val, &nodecfg)
//	if err != nil {
//		log.Printf("redis load error: b:%s", rcs.node)
//		return
//	}
//
//	val, err = rcs.client.HGetAll("n:" + rcs.node).Result()
//	if err != nil {
//		log.Printf("redis load error: b:%s", rcs.node)
//		return
//	}
//
//	err = LoadStructFromMap(val, &nodecfg)
//	if err != nil {
//		log.Printf("redis load error: b:%s", rcs.node)
//		return
//	}
//	log.Printf("node config loaded.")
//	return
//}

//func (rcs *RedisConfigSource) LoadBackends() (backends map[string]*BackendConfig, err error) {
//	backends = make(map[string]*BackendConfig)
//
//	names, err := rcs.client.Keys("b:*").Result()
//	if err != nil {
//		log.Printf("read redis error: %s", err)
//		return
//	}
//
//	var cfg *BackendConfig
//	for _, name := range names {
//		name = name[2:len(name)]
//		cfg, err = rcs.LoadConfigFromRedis(name)
//		if err != nil {
//			log.Printf("read redis config error: %s", err)
//			return
//		}
//		backends[name] = cfg
//	}
//	log.Printf("%d backends loaded from redis.", len(backends))
//	return
//}
//
//func (rcs *RedisConfigSource) LoadConfigFromRedis(name string) (cfg *BackendConfig, err error) {
//	val, err := rcs.client.HGetAll("b:" + name).Result()
//	if err != nil {
//		log.Printf("redis load error: b:%s", name)
//		return
//	}
//
//	cfg = &BackendConfig{}
//	err = LoadStructFromMap(val, cfg)
//	if err != nil {
//		return
//	}
//
//	if cfg.Interval == 0 {
//		cfg.Interval = 1000
//	}
//	if cfg.Timeout == 0 {
//		cfg.Timeout = 10000
//	}
//	if cfg.TimeoutQuery == 0 {
//		cfg.TimeoutQuery = 600000
//	}
//	if cfg.MaxRowLimit == 0 {
//		cfg.MaxRowLimit = 10000
//	}
//	if cfg.CheckInterval == 0 {
//		cfg.CheckInterval = 1000
//	}
//	if cfg.RewriteInterval == 0 {
//		cfg.RewriteInterval = 10000
//	}
//	return
//}

//func (rcs *RedisConfigSource) LoadMeasurements() (m_map map[string][]string, err error) {
//	m_map = make(map[string][]string, 0)
//
//	names, err := rcs.client.Keys("m:*").Result()
//	if err != nil {
//		log.Printf("read redis error: %s", err)
//		return
//	}
//
//	var length int64
//	for _, key := range names {
//		length, err = rcs.client.LLen(key).Result()
//		if err != nil {
//			return
//		}
//		m_map[key[2:len(key)]], err = rcs.client.LRange(key, 0, length).Result()
//		if err != nil {
//			return
//		}
//	}
//	log.Printf("%d measurements loaded from redis.", len(m_map))
//	return
//}
