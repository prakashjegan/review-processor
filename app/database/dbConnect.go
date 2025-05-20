// Package database handles connections to different
// types of databases
package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	//"github.com/olivere/elastic"
	"github.com/prakashjegan/review-processor/app/config"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// Import PostgreSQL database driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"

	// Import Redis Driver
	"github.com/mediocregopher/radix/v4"

	log "github.com/sirupsen/logrus"
	//"google.golang.org/api/option"
)

// dbClient variable to access gorm
var dbClient *gorm.DB

var sqlDB *sql.DB
var err error

// redisClient variable to access the redis client
var redisClient *radix.Client

// RedisConnTTL - context deadline in second
var RedisConnTTL int

// InitDB - function to initialize db
func InitDB() *gorm.DB {
	var db = dbClient

	configureDB := config.GetConfig().Database.RDBMS

	driver := configureDB.Env.Driver
	username := configureDB.Access.User
	password := configureDB.Access.Pass
	database := configureDB.Access.DbName
	host := configureDB.Env.Host
	port := configureDB.Env.Port
	sslmode := configureDB.Ssl.Sslmode
	timeZone := configureDB.Env.TimeZone
	maxIdleConns := configureDB.Conn.MaxIdleConns
	maxOpenConns := configureDB.Conn.MaxOpenConns
	connMaxLifetime := configureDB.Conn.ConnMaxLifetime
	logLevel := configureDB.Log.LogLevel
	log.Info("DriverName := ", driver)
	switch driver {
	case "mysql":
		dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
		if sslmode != "disable" {
			dsn += "&tls=custom"
			err = InitTLSMySQL()
			if err != nil {
				log.WithError(err).Panic("panic code: 150")
			}
		}
		sqlDB, err = sql.Open(driver, dsn)
		if err != nil {
			log.WithError(err).Panic("panic code: 151")
		}
		sqlDB.SetMaxIdleConns(maxIdleConns)       // max number of connections in the idle connection pool
		sqlDB.SetMaxOpenConns(maxOpenConns)       // max number of open connections in the database
		sqlDB.SetConnMaxLifetime(connMaxLifetime) // max amount of time a connection may be reused

		db, err = gorm.Open(mysql.New(mysql.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
		})
		if err != nil {
			log.WithError(err).Panic("panic code: 152")
		}
		// Only for debugging
		if err == nil {
			fmt.Println("DB connection successful!")
		}

	case "postgres":
		log.Info("PostgreSQL driver Started!")
		sslmode = "disable"
		dsn := "host=" + host + " port=" + port + " user=" + username + " dbname=" + database + " password=" + password + " sslmode=" + sslmode + " TimeZone=" + timeZone
		sqlDB, err = sql.Open(driver, dsn)
		if err != nil {
			log.WithError(err).Panic("panic code: 153")
		}
		sqlDB.SetMaxIdleConns(maxIdleConns)       // max number of connections in the idle connection pool
		sqlDB.SetMaxOpenConns(maxOpenConns)       // max number of open connections in the database
		sqlDB.SetConnMaxLifetime(connMaxLifetime) // max amount of time a connection may be reused

		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(logLevel)),
		})
		if err != nil {
			log.WithError(err).Panic("panic code: 154")
		}
		// Only for debugging
		if err == nil {
			fmt.Println("DB connection successful!")
		}
		log.Info("PostgreSQL driver Completed")

	case "sqlite3":
		db, err = gorm.Open(sqlite.Open(database), &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
		})
		if err != nil {
			log.WithError(err).Panic("panic code: 155")
		}
		// Only for debugging
		if err == nil {
			fmt.Println("DB connection successful!")
		}

	default:
		log.Fatal("The driver " + driver + " is not implemented yet")
	}

	dbClient = db

	return dbClient
}

// GetDB - get a connection
func GetDB() *gorm.DB {
	return dbClient
}

// InitRedis - function to initialize redis client
func InitRedis() (*radix.Client, error) {
	configureRedis := config.GetConfig().Database.REDIS

	RedisConnTTL = configureRedis.Conn.ConnTTL
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(RedisConnTTL)*time.Second)
	defer cancel()

	rClient, err := (radix.PoolConfig{
		Size: configureRedis.Conn.PoolSize,
	}).New(ctx, "tcp", fmt.Sprintf("%v:%v",
		configureRedis.Env.Host,
		configureRedis.Env.Port))
	if err != nil {
		log.WithError(err).Panic("panic code: 161")
		return &rClient, err
	}
	// Only for debugging
	if err == nil {
		fmt.Println("REDIS pool connection successful!")
	}

	redisClient = &rClient

	return redisClient, nil
}

// GetRedis - get a connection
func GetRedis() *radix.Client {
	return redisClient
}
