// Package config is responsible for reading all environment
// variables and set up the base configuration for a
// functional application
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Activated - "yes" keyword to activate a service
const Activated string = "yes"

// Configuration - server and db configuration variables
type Configuration struct {
	Database DatabaseConfig
	Logger   LoggerConfig
	Aws      AwsConfig
	Kafka    KafkaConfig

	Scheduler SchedulerConfig
}

var configAll *Configuration

// Env - load the configurations from .env
func Env() error {
	// Load environment variables
	return godotenv.Load()
}

// Config - load all the configurations
func Config() (err error) {
	var configuration Configuration

	configuration.Database, err = database()
	if err != nil {
		return
	}

	configuration.Logger, err = logger()
	if err != nil {
		return
	}

	configuration.Aws, err = aws()
	if err != nil {
		return
	}

	configuration.Kafka, err = kafka()
	if err != nil {
		return
	}

	configuration.Scheduler, err = scheduler()
	if err != nil {
		return
	}

	configAll = &configuration

	return
}

func scheduler() (schedulerConfig SchedulerConfig, err error) {
	err = Env()
	if err != nil {
		return
	}

	schedulerConfig.MaxRetries, _ = strconv.Atoi(strings.TrimSpace(os.Getenv("MAX_RETRIES")))
	schedulerConfig.MaxDelay, _ = strconv.Atoi(strings.TrimSpace(os.Getenv("MAX_DELAY")))
	schedulerConfig.SchedulerCrons = make(map[string]string)
	for _, cron := range strings.Split(strings.TrimSpace(os.Getenv("SCHEDULER_INTERVAL_1")), ",") {
		parts := strings.Split(cron, ":")
		if len(parts) == 2 {
			schedulerConfig.SchedulerCrons[parts[0]] = parts[1]
		}
	}
	return schedulerConfig, nil
}

func kafka() (kafkaConfig KafkaConfig, err error) {
	err = Env()
	if err != nil {
		return
	}

	kafkaConfig.Activate = strings.TrimSpace(os.Getenv("ACTIVATE_KAFKA"))
	if kafkaConfig.Activate == Activated {
		kafkaConfig.Env.Host = strings.TrimSpace(os.Getenv("KAFKA_HOST"))
		kafkaConfig.Env.Port = strings.TrimSpace(os.Getenv("KAFKA_PORT"))
	}

	return kafkaConfig, nil
}

// GetConfig - return all the config variables
func GetConfig() *Configuration {
	return configAll
}

// database - all DB variables
func database() (databaseConfig DatabaseConfig, err error) {
	// Load environment variables
	err = Env()
	if err != nil {
		return
	}

	// RDBMS
	activateRDBMS := strings.TrimSpace(os.Getenv("ACTIVATE_RDBMS"))
	if activateRDBMS == Activated {
		dbRDBMS, errThis := databaseRDBMS()
		if errThis != nil {
			err = errThis
			return
		}
		databaseConfig.RDBMS = dbRDBMS.RDBMS
	}
	databaseConfig.RDBMS.Activate = activateRDBMS

	// REDIS
	activateRedis := strings.TrimSpace(os.Getenv("ACTIVATE_REDIS"))
	if activateRedis == Activated {
		dbRedis, errThis := databaseRedis()
		if errThis != nil {
			err = errThis
			return
		}
		databaseConfig.REDIS = dbRedis.REDIS
	}
	databaseConfig.REDIS.Activate = activateRedis

	return
}

// databaseRDBMS - all RDBMS variables
func databaseRDBMS() (databaseConfig DatabaseConfig, err error) {
	// Load environment variables
	err = Env()
	if err != nil {
		return
	}

	// Env
	databaseConfig.RDBMS.Env.Driver = strings.TrimSpace(os.Getenv("DB_DRIVER"))
	databaseConfig.RDBMS.Env.Host = strings.TrimSpace(os.Getenv("DB_HOST"))
	databaseConfig.RDBMS.Env.Port = strings.TrimSpace(os.Getenv("DB_PORT"))
	databaseConfig.RDBMS.Env.TimeZone = strings.TrimSpace(os.Getenv("DB_TIMEZONE"))
	// Access
	databaseConfig.RDBMS.Access.DbName = strings.TrimSpace(os.Getenv("DB_NAME"))
	databaseConfig.RDBMS.Access.User = strings.TrimSpace(os.Getenv("DB_USER"))
	databaseConfig.RDBMS.Access.Pass = strings.TrimSpace(os.Getenv("DB_PASSWORD"))
	// SSL
	databaseConfig.RDBMS.Ssl.Sslmode = strings.TrimSpace(os.Getenv("DB_SSLMODE"))
	databaseConfig.RDBMS.Ssl.MinTLS = strings.TrimSpace(os.Getenv("DB_SSL_TLS_MIN"))
	databaseConfig.RDBMS.Ssl.RootCA = strings.TrimSpace(os.Getenv("DB_SSL_ROOT_CA"))
	databaseConfig.RDBMS.Ssl.ServerCert = strings.TrimSpace(os.Getenv("DB_SSL_SERVER_CERT"))
	databaseConfig.RDBMS.Ssl.ClientCert = strings.TrimSpace(os.Getenv("DB_SSL_CLIENT_CERT"))
	databaseConfig.RDBMS.Ssl.ClientKey = strings.TrimSpace(os.Getenv("DB_SSL_CLIENT_KEY"))
	// Conn
	dbMaxIdleConns := strings.TrimSpace(os.Getenv("DB_MAXIDLECONNS"))
	dbMaxOpenConns := strings.TrimSpace(os.Getenv("DB_MAXOPENCONNS"))
	dbConnMaxLifetime := strings.TrimSpace(os.Getenv("DB_CONNMAXLIFETIME"))
	databaseConfig.RDBMS.Conn.MaxIdleConns, err = strconv.Atoi(dbMaxIdleConns)
	if err != nil {
		return
	}
	databaseConfig.RDBMS.Conn.MaxOpenConns, err = strconv.Atoi(dbMaxOpenConns)
	if err != nil {
		return
	}
	databaseConfig.RDBMS.Conn.ConnMaxLifetime, err = time.ParseDuration(dbConnMaxLifetime)
	if err != nil {
		return
	}

	// Logger
	dbLogLevel := strings.TrimSpace(os.Getenv("DB_LOG_LEVEL"))
	databaseConfig.RDBMS.Log.LogLevel, err = strconv.Atoi(dbLogLevel)
	if err != nil {
		return
	}

	return
}

// databaseRedis - all REDIS DB variables
func databaseRedis() (databaseConfig DatabaseConfig, err error) {
	// Load environment variables
	err = Env()
	if err != nil {
		return
	}

	// REDIS
	poolSize, errThis := strconv.Atoi(strings.TrimSpace(os.Getenv("POOL_SIZE")))
	if errThis != nil {
		err = errThis
		return
	}
	connTTL, errThis := strconv.Atoi(strings.TrimSpace(os.Getenv("CONN_TTL")))
	if errThis != nil {
		err = errThis
		return
	}

	databaseConfig.REDIS.Env.Host = strings.TrimSpace(os.Getenv("REDIS_HOST"))
	databaseConfig.REDIS.Env.Port = strings.TrimSpace(os.Getenv("REDIS_PORT"))
	databaseConfig.REDIS.Conn.PoolSize = poolSize
	databaseConfig.REDIS.Conn.ConnTTL = connTTL

	return
}

// logger - config for sentry.io
func logger() (loggerConfig LoggerConfig, err error) {
	// Load environment variables
	err = Env()
	if err != nil {
		return
	}

	loggerConfig.Activate = strings.TrimSpace(os.Getenv("ACTIVATE_SENTRY"))
	if loggerConfig.Activate == Activated {
		loggerConfig.SentryDsn = strings.TrimSpace(os.Getenv("SentryDSN"))
	}

	return
}

// view - HTML renderer
func aws() (awsConfig AwsConfig, err error) {
	// Load environment variables
	err = Env()
	if err != nil {
		return
	}

	awsConfig.Activate = strings.TrimSpace(os.Getenv("ACTIVATE_AWS"))
	if awsConfig.Activate == Activated {
		awsConfig.AccessKey = strings.TrimSpace(os.Getenv("AWS_ACCESS_KEY_ID"))
		awsConfig.DocumentBucketName = strings.TrimSpace(os.Getenv("S3_DOCUMENT_BUCKET_NAME"))
		awsConfig.Region = strings.TrimSpace(os.Getenv("AWS_REGION"))
		awsConfig.SecreteAccessKey = strings.TrimSpace(os.Getenv("AWS_SECRET_ACCESS_KEY"))
	}

	return
}
