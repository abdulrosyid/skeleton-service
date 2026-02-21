package configs

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/caarlos0/env"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT"`

	MySQLHost     string `env:"MYSQL_HOST"`
	MySQLPort     int    `env:"MYSQL_PORT"`
	MySQLUser     string `env:"MYSQL_USER"`
	MySQLPassword string `env:"MYSQL_PASSWORD"`
	MySQLDatabase string `env:"MYSQL_DATABASE"`

	JWTSecret        string `env:"JWT_SECRET"`
	JWTExpireMinutes int    `env:"JWT_EXPIRE_MINUTES"`
}

var (
	configuration Config
	once          sync.Once
)

func GetConfig() Config {
	once.Do(func() {
		loadConfig()
	})
	return configuration
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	_ = viper.ReadInConfig()

	t := reflect.TypeOf(configuration)
	v := reflect.ValueOf(&configuration)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("env")
		if key == "" {
			continue
		}
		val := viper.GetString(key)
		if val == "" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			v.Elem().FieldByName(field.Name).SetString(val)
		case reflect.Int:
			if iVal, err := strconv.Atoi(val); err == nil {
				v.Elem().FieldByName(field.Name).SetInt(int64(iVal))
			}
		}
	}

	// env vars have highest priority for deployment/runtime configuration.
	if err := env.Parse(&configuration); err != nil {
		log.Printf("warning parse env: %v", err)
	}

	if configuration.HTTPPort == "" {
		configuration.HTTPPort = "8080"
	}
	if configuration.JWTSecret == "" {
		configuration.JWTSecret = "super-secret-change-me"
	}
	if configuration.JWTExpireMinutes <= 0 {
		configuration.JWTExpireMinutes = 60
	}
}

func (c Config) MySQLDSN() string {
	port := c.MySQLPort
	if port == 0 {
		port = 3306
	}
	return c.MySQLUser + ":" + c.MySQLPassword +
		"@tcp(" + c.MySQLHost + ":" + strconv.Itoa(port) + ")/" +
		c.MySQLDatabase + "?parseTime=true&loc=Asia%2FJakarta"
}

func ServerShutdownTimeout() time.Duration {
	return 10 * time.Second
}

func init() {
	_ = os.Setenv("TZ", "Asia/Jakarta")
}
