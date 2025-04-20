package redis

import (
	"os"
	"testing"
	"time"

	test "github.com/mahopon/gobackend/testing"
	"github.com/stretchr/testify/assert"
)

func init() {
	loadEnv()
}

func TestRedisIntegration(t *testing.T) {
	test.CheckEnvironment(t, "INTEGRATION")
	os.Setenv("REDIS_URL", "redis://:dellmewhatdodo@localhost:6379/0?protocol=3")

	client := GetClient()
	defer client.CloseClient()

	err := client.Set("testkey", "testvalue", 5*time.Second)
	assert.NoError(t, err)

	val, err := client.Get("testkey")
	assert.NoError(t, err)
	assert.Equal(t, "testvalue", val)

	_ = client.Del("testkey")
}
