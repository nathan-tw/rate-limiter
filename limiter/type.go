package limiter

import "time"

type LimitType int

const (
	SlidingWindow LimitType = iota
	FixedWindow
)

type LimitParam struct {
	Key       string
	Timestamp int64
	Duration  time.Duration
	Value     int
	LimitType
}

type LimitRules struct {
	Duration  time.Duration `yaml:"Duration"`
	Value     int           `yaml:"Value"`
	LimitType LimitType     `yaml:"LimitType"`
}

type LimitConfig struct {
	AccountLimit  LimitRules `yaml:"AccountLimit"`
	EndpointLimit LimitRules `yaml:"EndpointLimit"`
}
