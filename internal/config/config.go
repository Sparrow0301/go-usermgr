package config

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database   DatabaseConf   `json:"Database"`
	JWT        JWTConf        `json:"JWT"`
	Password   PasswordConf   `json:"Password"`
	Pagination PaginationConf `json:"Pagination"`
	Security   SecurityConf   `json:"Security"`
}

type DatabaseConf struct {
	DSN             string        `json:"DSN"`
	MaxIdleConns    int           `json:"MaxIdleConns"`
	MaxOpenConns    int           `json:"MaxOpenConns"`
	ConnMaxLifetime time.Duration `json:"ConnMaxLifetime"`
}

type JWTConf struct {
	AccessSecret  string        `json:"AccessSecret"`
	AccessExpire  time.Duration `json:"AccessExpire"`
	RefreshExpire time.Duration `json:"RefreshExpire"`
}

type PasswordConf struct {
	BcryptCost int `json:"BcryptCost"`
}

type PaginationConf struct {
	DefaultPageSize int `json:"DefaultPageSize"`
	MaxPageSize     int `json:"MaxPageSize"`
}

type SecurityConf struct {
	AllowOrigins []string `json:"AllowOrigins"`
}
