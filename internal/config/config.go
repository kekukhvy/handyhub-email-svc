package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	Database Database       `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Server   ServerSettings `mapstructure:"server"`
	App      Application    `mapstructure:"app"`
	Logs     LogsSettings   `mapstructure:"logs"`
	Security Security
	Redis    Redis
	Email    Email       `mapstructure:"email"`
	Queue    QueueConfig `mapstructure:"queue"`
	SMTP     SMTPConfig  `mapstructure:"smtp"`
}

type Database struct {
	Url             string `mapstructure:"url"`
	DbName          string `mapstructure:"dbname"`
	EmailCollection string `mapstructure:"email-collection"`
	Timeout         int    `mapstructure:"timeout"`
}

type ServerSettings struct {
	Port         string `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  int    `mapstructure:"read-timeout"`
	WriteTimeout int    `mapstructure:"write-timeout"`
	IdleTimeout  int    `mapstructure:"idle-timeout"`
}

type Application struct {
	Name    string `mapstructure:"name"`
	Timeout int    `mapstructure:"timeout"`
}

type LogsSettings struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"log-path"`
}

type Security struct {
	JwtKey           string `mapstructure:"jwt-key"`
	ExpirationPerion int    `mapstructure:"expiration_period"`
}

type Redis struct {
	Url      string `mapstructure:"url"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

type Email struct {
	SMTPHost     string `mapstructure:"smtp-host"`
	SMTPPort     int    `mapstructure:"smtp-port"`
	SMTPUser     string `mapstructure:"smtp-user"`
	SMTPPassword string `mapstructure:"smtp-password"`
	FromEmail    string `mapstructure:"from-email"`
	FromName     string `mapstructure:"from-name"`
}

type StorageConfig struct {
	Type string            `mapstructure:"type"`
	File FileStorageConfig `mapstructure:"file"`
}

type FileStorageConfig struct {
	Path      string `mapstructure:"path"`
	MaxSizeMB int    `mapstructure:"max-size-mb"`
	MaxFiles  int    `mapstructure:"max-files"`
}

type QueueConfig struct {
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
}

type RabbitMQConfig struct {
	Url            string `mapstructure:"url"`
	Exchange       string `mapstructure:"exchange"`
	ExchangeType   string `mapstructure:"exchange-type"`
	EmailQueue     string `mapstructure:"email-queue"`
	PrefetchCount  int    `mapstructure:"prefetch-count"`
	ReconnectDelay int    `mapstructure:"reconnect-delay"`
	Timeout        int    `mapstructure:"timeout"`
	RoutingKey     string `mapstructure:"routing-key"`
	PrefetchSize   int    `mapstructure:"prefetch-size"`
	Global         bool   `mapstructure:"global"`
	Durable        bool   `mapstructure:"durable"`
	AutoDelete     bool   `mapstructure:"auto-delete"`
	Internal       bool   `mapstructure:"internal"`
	NoWait         bool   `mapstructure:"no-wait"`
	Exclusive      bool   `mapstructure:"exclusive"`
	AutoAck        bool   `mapstructure:"auto-ack"`
	NoLocal        bool   `mapstructure:"no-local"`
	Consumer       string `mapstructure:"consumer"`
}

type SMTPConfig struct {
	Provider    string         `mapstructure:"provider"`
	DefaultFrom string         `mapstructure:"default-from"`
	Gmail       GmailConfig    `mapstructure:"gmail"`
	SendGrid    SendGridConfig `mapstructure:"sendgrid"`
	MailHog     MailHogConfig  `mapstructure:"mailhog"`
}

type GmailConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
}

type SendGridConfig struct {
	ApiKey string `mapstructure:"api-key"`
	Url    string `mapstructure:"url"`
}

type MailHogConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func Load() *Configuration {

	cfg := read()
	logrus.Info("Configuration loaded")

	mongoUri := os.Getenv("MONGODB_URL")
	if mongoUri != "" {
		cfg.Database.Url = mongoUri
	}

	dbName := os.Getenv("DB_NAME")
	if dbName != "" {
		cfg.Database.DbName = dbName
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost != "" {
		cfg.Redis.Url = redisHost
	}

	redisDB := os.Getenv("REDIS_DB")
	if redisDB != "" {
		cfg.Redis.Db, _ = strconv.Atoi(redisDB)
	}

	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey != "" {
		cfg.Security.JwtKey = jwtKey
	}

	return cfg
}

func read() *Configuration {
	viper.SetConfigFile("internal/config/cfg.yml")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var config Configuration

	err := viper.ReadInConfig()

	if err != nil {
		logrus.Panic("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		logrus.Panic("Error unmarshalling config file, %s", err)
	}

	return &config
}
