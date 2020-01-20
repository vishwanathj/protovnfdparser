package config

import (
	"sync"

	"github.com/kelseyhightower/envconfig"
	"github.com/vishwanathj/protovnfdparser/pkg/constants"
)

type Config struct {
	MongoDBConfig *MongoDBConfig    `envconfig:"mongodbcfg"`
	PgntConfig    *PaginationConfig `envconfig:"pgncfg"`
	WebConfig     *WebServerConfig  `envconfig:"webcfg"`
}

var cfgInstance *Config
var once sync.Once

type MongoDBConfig struct {
	MongoIP      string `envconfig:"MONGO_IP" default:"mongo"`
	MongoPort    uint   `envconfig:"MONGO_PORT" default:"27017"`
	MongoDBName  string `envconfig:"MONGO_DB_NAME" default:"VNFDSVCDB"`
	MongoColName string `envconfig:"MONGO_COL_NAME" default:"vnfd"`
}

type PaginationConfig struct {
	BaseURI        string 	`envconfig:"BASE_URI" default:"http://localhost:8080"`
	DefaultLimit   int 	`envconfig:"DEFAULT_LIMIT" default:"5"`
	MaxLimit       int 	`envconfig:"MAX_LIMIT" default:"10"`
	MinLimit       int 	`envconfig:"MIN_LIMIT" default:"1"`
	DefaultOrderBy string 	`envconfig:"DEFAULT_ORDER_BY" default:"name"`
}

type WebServerConfig struct {
	WebServerPort       uint   `envconfig:"WEB_SERVER_PORT" default:"8080"`
	WebServerSecurePort uint   `envconfig:"WEB_SERVER_SECURE_PORT" default:"443"`
	WebServerBasePath   string `envconfig:"WEB_SERVER_BASE_PATH" default:"/"`
}

func GetConfigInstance() *Config {
	once.Do(func() {
		var cfg Config
		envconfig.MustProcess(constants.ENVCONFIG_PREFIX, &cfg)
		cfgInstance = &Config{cfg.MongoDBConfig, cfg.PgntConfig, cfg.WebConfig}
	})
	return cfgInstance
}
