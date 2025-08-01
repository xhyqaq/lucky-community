package config

// 读取配置信息并且给各个配置类赋值

import (
	"os"
	"strconv"
)

type AppConfig struct {
	ServerBind      string      `yaml:"serverBind" default:":8080"`
	DbConfig        DbConfig    `yaml:"db"`
	PGDbConfig      PGDbConfig  `yaml:"pgdb"`
	OssConfig       OssConfig   `yaml:"oss"`
	EmailConfig     EmailConfig `yaml:"email"`
	LLMConfig       LLMConfig
	EmbeddingConfig EmbeddingConfig
	GitHubConfig    GitHubConfig `yaml:"github"`
}

type EmbeddingConfig struct {
	ApiKey string `yaml:"apikey"`
	Model  string `yaml:"model"`
	Url    string `yaml:"url"`
}

type GitHubConfig struct {
	ClientID     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	RedirectURL  string `yaml:"redirectUrl"`
}

type LLMConfig struct {
	ApiKey string `yaml:"apikey"`
	Model  string `yaml:"model"`
	Url    string `yaml:"url"`
}
type PGDbConfig struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type DbConfig struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type OssConfig struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
	Cdn       string `yaml:"cdn"`
	Callback  string `json:"callback"`
	Endpoint  string `yaml:"endpoint"`
}

type EmailConfig struct {
	Address   string `yaml:"address"`
	PollCount int    `yaml:"pollCount"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
}

var instance *AppConfig

func GetInstance() *AppConfig {
	return instance
}

func Init() {
	pollCount, _ := strconv.Atoi(os.Getenv("POLLCOUNT"))
	if pollCount == 0 {
		pollCount = 10
	}

	// 设置数据库默认值，如果环境变量为空
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}
	dbDatabase := os.Getenv("DB_DATABASE")
	if dbDatabase == "" {
		dbDatabase = "community"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPass := os.Getenv("DB_PASS")

	serverBind := os.Getenv("SERVER_BIND")
	if serverBind == "" {
		serverBind = ":8080"
	}

	appConfig := &AppConfig{
		ServerBind: serverBind,
		DbConfig: DbConfig{
			Address:  dbHost,
			Database: dbDatabase,
			Username: dbUser,
			Password: dbPass,
		},
		PGDbConfig: PGDbConfig{
			Address:  os.Getenv("PG_DB_HOST"),
			Database: os.Getenv("PG_DB_DATABASE"),
			Username: os.Getenv("PG_DB_USER"),
			Password: os.Getenv("PG_DB_PASS"),
		},
		OssConfig: OssConfig{
			AccessKey: os.Getenv("OSS_ACCESS_KEY"),
			Bucket:    os.Getenv("OSS_BUCKET"),
			SecretKey: os.Getenv("OSS_SECRET_KEY"),
			Cdn:       os.Getenv("OSS_CDN"),
			Callback:  os.Getenv("OSS_CALLBACK"),
			Endpoint:  os.Getenv("OSS_ENDPOINT"),
		},
		EmailConfig: EmailConfig{
			Address:   os.Getenv("ADDRESS"),
			Username:  os.Getenv("USERNAME"),
			Password:  os.Getenv("PASSWORD"),
			Host:      os.Getenv("HOST"),
			PollCount: pollCount,
		},
		LLMConfig: LLMConfig{
			ApiKey: os.Getenv("LLM_API_KEY"),
			Model:  os.Getenv("LLM_MODEL"),
			Url:    os.Getenv("LLM_URL"),
		},
		EmbeddingConfig: EmbeddingConfig{
			ApiKey: os.Getenv("EMBEDDING_API_KEY"),
			Model:  os.Getenv("EMBEDDING_LLM_MODEL"),
			Url:    os.Getenv("EMBEDDING_LLM_URL"),
		},
		GitHubConfig: GitHubConfig{
			ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
			ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GITHUB_REDIRECT_URL"),
		},
	}
	instance = appConfig

}
