package utility

import (
	"fmt"
	"github.com/go-gorm/caches/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"time"
)

// NewRDB create rational database instance
func NewRDB(driver RdbProperties, plugins ...gorm.Plugin) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	var cachePlugin gorm.Plugin
	//
	if driver.Name == "mysql" {
		db, err = NewMySQL(driver)
	} else if driver.Name == "postgres" {
		db, err = NewPostgres(driver)
	} else {
		db, err = nil, errors.New("not support driver: "+driver.Name)
	}
	if err != nil {
		return nil, err
	} else if len(driver.Cache) > 0 { // cache
		if strings.HasPrefix(driver.Cache, "mem://") {
			if memProtocol, err := url.Parse(driver.Cache); err == nil {
				cachePlugin = &caches.Caches{Conf: &caches.Config{
					Cacher: (&memoryCacher{}).size(QueryInt(memProtocol.Query(), "size")),
				}}
			} else {
				DefaultLogger().Error("rdb cache inactivate, because of create failed: " + driver.Cache)
			}
		} else if strings.HasPrefix(driver.Cache, "redis://") {
			cachePlugin = &caches.Caches{Conf: &caches.Config{
				Cacher: &redisCacher{rdb: NewRedis(driver.Cache)},
			}}
		} else {
			DefaultLogger().Error("rdb not support this cache: " + driver.Cache)
		}
	}
	//
	if db != nil {
		if cachePlugin != nil { // // cache plugin
			_ = db.Use(cachePlugin)
		}
		if plugins != nil && len(plugins) > 0 {
			for _, plugin := range plugins {
				_ = db.Use(plugin)
			}
		}
	}
	//
	return initRbdConnPool(db, driver)
}

// NewMySQL create gorm.DB instance based on mysql database
func NewMySQL(driver RdbProperties) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", driver.Username, driver.Password, driver.Host, driver.Port, driver.Database)
	if len(driver.Parameters) > 0 {
		dsn += "&" + driver.Parameters
	}
	//
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // autoconfigure based on currently MySQL version
	}), defaultConfig())
}

// NewPostgres create gorm.DB instance based on postgres database
func NewPostgres(driver RdbProperties) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s TimeZone=Asia/Shanghai",
		driver.Host, driver.Port, driver.Username, driver.Password, driver.Database)
	if len(driver.Parameters) > 0 {
		dsn += " " + driver.Parameters
	}
	//
	return gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), defaultConfig())
}

func initRbdConnPool(gdb *gorm.DB, driver RdbProperties) (*gorm.DB, error) {
	db, err := gdb.DB()
	if err != nil {
		DefaultLogger().Error("get db instance from gorm error", zap.Any("driver", driver), zap.Error(err))
		return nil, err
	}
	db.SetMaxIdleConns(MaxInt(driver.Pool.MaxIdleCon, 5))
	db.SetMaxOpenConns(MaxInt(driver.Pool.MaxOpenCon, 5))
	db.SetConnMaxLifetime(time.Duration(MaxInt64(1, driver.Pool.MaxLifeTimeInMinute) * int64(time.Minute)))
	if err = db.Ping(); err != nil {
		DefaultLogger().Error("connect to db instance failed", zap.Any("driver", driver), zap.Error(err))
		return nil, err
	}
	return gdb, nil
}

func defaultConfig() *gorm.Config {
	return &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	}
}
