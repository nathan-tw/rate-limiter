package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nathan-tw/swif_devops_assignment/limiter"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

type RedisConfig struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
}

func main() {

	yamlFile, err := os.ReadFile("config/redis_config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
		return
	}
	var config RedisConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
		return
	}
	redisAddr := fmt.Sprintf("%v:%v", config.Host, config.Port)
	fmt.Println(redisAddr)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0,
	})
	if redisClient == nil {
		log.Fatal("Failed to new redis")
		return
	}
	limiter, err := limiter.NewLimiter(redisClient)
	if err != nil {
		log.Fatalf("Failed to new rate limiter: %v", err)
		return
	}
	r := gin.Default()
	r.GET("/flush_redis", func(c *gin.Context) {
		if status := redisClient.FlushAll(c); status.Err() != nil {
			c.String(http.StatusInternalServerError, status.Err().Error())
		} else {
			c.String(http.StatusOK, "successfully flush redis")
		}
	})
	r.GET("/reload_limit_config", func(c *gin.Context) {
		err := limiter.ReloadConfig()
		if err != nil {
			c.String(http.StatusInternalServerError, "failed to reload rate limit config")
		}
		c.String(http.StatusOK, "config reloaded")
	})
	apiGroup := r.Group("/api", limiter.LimitEndpoint(), limiter.LimitAccount())
	apiGroup.GET("/path1", func(c *gin.Context) {
		c.String(http.StatusOK, "path1 visited")
	})
	apiGroup.GET("/path2", func(c *gin.Context) {
		c.String(http.StatusOK, "path2 visited")
	})

	r.Run("0.0.0.0:8080")
}
