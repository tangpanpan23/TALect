package config

import (
	"flag"

	"github.com/zeromicro/go-zero/rest"
)

var ConfigFile = flag.String("f", "etc/talink.yaml", "the config file")

type Config struct {
	rest.RestConf
	Database struct {
		Host         string `json:",default=localhost"`
		Port         int    `json:",default=5432"`
		User         string `json:",default=future_mcp"`
		Password     string `json:",default=password"`
		Dbname       string `json:",default=future_mcp"`
		Sslmode      string `json:",default=disable"`
		MaxIdleConns int    `json:",default=10"`
		MaxOpenConns int    `json:",default=100"`
	} `json:",optional"`
	Redis struct {
		Host      string `json:",default=localhost:6379"`
		Password  string `json:",default="`
		Db        int    `json:",default=0"`
		PoolSize  int    `json:",default=10"`
	} `json:",optional"`
	Auth struct {
		JwtSecret  string `json:",default=your-jwt-secret-key"`
		JwtExpire  int    `json:",default=86400"`
	} `json:",optional"`
	VectorSearch struct {
		Provider   string `json:",default=pinecone"`
		ApiKey     string `json:",optional"`
		Environment string `json:",default=us-east-1"`
		IndexName  string `json:",default=future-materials"`
		Dimension  int    `json:",default=768"`
	} `json:",optional"`
	Storage struct {
		Provider   string `json:",default=local"`
		Bucket     string `json:",default=future-materials"`
		Region     string `json:",default=us-east-1"`
		Endpoint   string `json:",optional"`
		AccessKey  string `json:",optional"`
		SecretKey  string `json:",optional"`
	} `json:",optional"`
}
