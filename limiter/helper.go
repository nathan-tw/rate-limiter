package limiter

import (
	"crypto/md5"
	"fmt"
)

func getMD5(data string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
