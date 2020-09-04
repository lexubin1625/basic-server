package db

import (
	"basic-server/core/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/url"
	"time"
)

var (
	clients map[string]map[string]*gorm.DB
)

const (
	Master              = "master"
	Slave               = "slave"
	DefaultCharset      = "utf8"
	DefaultTimetout     = "50ms"
	DefaultReadTimeout  = "50ms"
	DefaultWriteTimeout = "50ms"
	DefaultMaxIdleConns = 10
	DefaultMaxOpenConns = 100
	DefaultMaxLifetime  = "1h"
)

// Conf 数据库配置
type Conf struct {
	Dsn          string `mapstructure:"dsn"`
	SlaveDsn     string `mapstructure:"slave_dsn"`
	Charset      string `mapstructure:"charset"`
	Timeout      string `mapstructure:"timeout"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	MaxIdleConn  int    `mapstructure:"max_idle_conn"`
	MaxOpenConn  int    `mapstructure:"max_open_conn"`
	MaxLifetime  string `mapstructure:"max_lifetime"`
}

// DefaultConf 默认配置
func DefaultConf() Conf {
	return Conf{
		Dsn:          "",
		SlaveDsn:     "",
		Charset:      DefaultCharset,
		Timeout:      DefaultTimetout,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		MaxIdleConn:  DefaultMaxIdleConns,
		MaxOpenConn:  DefaultMaxOpenConns,
		MaxLifetime:  DefaultMaxLifetime,
	}
}

// 初始化配置
func initConf() map[string]Conf {
	databases := config.Viper.GetStringMap("database.mysql.node")
	conf := make(map[string]Conf)
	for k, v := range databases {
		addrMap := DefaultConf()
		err := mapstructure.Decode(v, &addrMap)
		if err != nil {
			//log.Panicln(err)
		}

		conf[k] = addrMap

	}
	return conf

}

// New 初始化数据库连接
func New() error {
	mysqlConf := initConf()
	clients = make(map[string]map[string]*gorm.DB, len(mysqlConf))
	for k, v := range mysqlConf {
		args := url.Values{}
		args.Set("charset", v.Charset)
		args.Set("parseTime", "true")
		args.Set("loc", "Local")
		args.Set("timeout", v.Timeout)
		args.Set("readTimeout", v.ReadTimeout)
		args.Set("writeTimeout", v.WriteTimeout)
		params := args.Encode()
		clients[k] = make(map[string]*gorm.DB, 2)
		conn, err := gorm.Open("mysql", v.Dsn+"?"+params)
		if err != nil {
			//log.Errorf(nil, "mysql connect error: %s", err.Error())
			return err
		}
		conn.DB().SetMaxIdleConns(v.MaxIdleConn)
		conn.DB().SetMaxOpenConns(v.MaxOpenConn)
		if d, err := time.ParseDuration(v.MaxLifetime); err == nil {
			conn.DB().SetConnMaxLifetime(d)
		}
		clients[k][Master] = conn
		if v.SlaveDsn != "" {
			conn, err := gorm.Open("mysql", v.SlaveDsn+"?"+params)
			if err != nil {
				//log.Errorf(nil, "mysql connect error: %s", err.Error())
				return err
			}
			conn.DB().SetMaxIdleConns(v.MaxIdleConn)
			conn.DB().SetMaxOpenConns(v.MaxOpenConn)
			if d, err := time.ParseDuration(v.MaxLifetime); err == nil {
				conn.DB().SetConnMaxLifetime(d)
			}
			clients[k][Slave] = conn
		}
	}
	return nil
}

// Get 获取数据库连接，第二个参数设置为true表示获取主库连接
func Get(name string, master bool) (*gorm.DB, bool) {

	conn, ok := clients[name]
	if !ok {
		return nil, false
	}
	key := Slave
	if master {
		key = Master
	}
	client, ok := conn[key]
	return client, ok
}
