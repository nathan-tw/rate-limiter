package limiter

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Limiter struct {
	RedisClient *redis.Client
	LimitConfig *LimitConfig
}

func NewLimiter(redisClient *redis.Client) (*Limiter, error) {
	var limitConfig LimitConfig
	yfile, err := os.ReadFile("config/limit_config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yfile, &limitConfig)
	if err != nil {
		return nil, err
	}
	return &Limiter{
		RedisClient: redisClient,
		LimitConfig: &limitConfig,
	}, nil
}

func (l *Limiter) ReloadConfig() error {
	var limitConfig LimitConfig
	yfile, err := os.ReadFile("config/limit_config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yfile, &limitConfig)
	l.LimitConfig = &limitConfig
	return nil
}

func (l *Limiter) LimitAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		account := c.Request.Header.Get("account")
		param := &LimitParam{
			Key:       "limiter:account:" + account,
			Timestamp: time.Now().UnixMicro(),
			Duration:  l.LimitConfig.AccountLimit.Duration,
			Value:     l.LimitConfig.AccountLimit.Value,
			LimitType: l.LimitConfig.AccountLimit.LimitType,
		}
		if isOver, err := l.limit(c, param); err != nil {
			c.AbortWithError(500, err)
		} else if isOver {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, over account limit"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func (l *Limiter) LimitEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.RequestURI
		param := &LimitParam{
			Key:       "limiter:endpoint:" + uri,
			Timestamp: time.Now().UnixMicro(),
			Duration:  l.LimitConfig.EndpointLimit.Duration,
			Value:     l.LimitConfig.EndpointLimit.Value,
			LimitType: l.LimitConfig.EndpointLimit.LimitType,
		}
		if isOver, err := l.limit(c, param); err != nil {
			c.AbortWithError(500, err)
		} else if isOver {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, over endpoint limit"})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func (l *Limiter) limit(ctx context.Context, param *LimitParam) (isOver bool, err error) {
	var (
		key        = param.Key
		duration   = param.Duration
		limitValue = param.Value
		ts         = param.Timestamp
		r          = l.RedisClient
	)
	if duration == 0 { // no rate limit exist
		return false, nil
	}
	switch param.LimitType {
	case SlidingWindow:
		lowerBound := strconv.FormatInt(ts-duration.Microseconds(), 10)
		pipe := r.TxPipeline()

		zRemRangeByScore := pipe.ZRemRangeByScore(ctx, key, "0", lowerBound)
		if zRemRangeByScore.Err() != nil {
			log.Printf("Error removing old entries: %v", zRemRangeByScore.Err())
			return false, err
		}

		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(ts),
			Member: ts,
		})
		pipe.Expire(ctx, key, duration)
		zcard := pipe.ZCard(ctx, key)

		_, err = pipe.Exec(ctx)

		if err != nil {
			log.Printf("Error executing pipeline: %v", err)
			return
		}

		if zRemRangeByScore.Val() != 0 {
			log.Printf("Removed %v entries for %v requests", key, zRemRangeByScore.Val())
		}

		log.Printf("Added %v for 1 request", key)
		log.Printf("%v's current requests count is %v", key, zcard.Val())

		// Check the rate limit
		isOver = limitValue < int(zcard.Val())

	case FixedWindow:
		script := redis.NewScript(luaScript)
		result, err := script.Run(ctx, r, []string{key}, limitValue, int(duration.Seconds())).Result()
		if err != nil {
			return false, err
		}

		switch result.(type) {
		case int64:
			return result.(int64) == 1, nil
		default:
			return false, errors.New("unexpected Lua script result")
		}

	}
	return
}
