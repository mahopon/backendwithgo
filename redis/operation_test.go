package redis

import (
	"log"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/joho/godotenv"
	test "github.com/mahopon/gobackend/testing"
	"github.com/stretchr/testify/assert"
)

// init runs before every test
// can use TestMain(m *testing.Main) if needed to
func init() {
	loadEnv()
}

func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func TestOperation(t *testing.T) {
	test.CheckEnvironment(t, "UNIT")
	db, mock := redismock.NewClientMock()
	defer db.Close()
	instance := &redisClient{client: db, raw: db}
	key, field, value := "key", "field", "value"
	t.Run("Set", func(t *testing.T) {
		mock.ExpectSet(key, value, 5*time.Second).SetVal("OK")
		err := instance.Set(key, value, 5*time.Second)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get", func(t *testing.T) {
		mock.ExpectGet(key).SetVal(value)
		val, err := instance.Get(key)
		assert.NoError(t, err)
		assert.Equal(t, value, val)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete", func(t *testing.T) {
		mock.ExpectDel(key).SetVal(1)
		err := instance.Delete(key)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exists", func(t *testing.T) {
		mock.ExpectExists(key).SetVal(1)
		exists, err := instance.Exists(key)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, true, exists)
	})

	t.Run("HSet", func(t *testing.T) {
		mock.ExpectHSet(key, field, value, 5*time.Second).SetVal(1)
		result, err := instance.HSet(key, field, value, 5*time.Second)
		assert.Equal(t, int64(1), result)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HGet", func(t *testing.T) {
		mock.ExpectHGet(key, field).SetVal("OK")
		result, err := instance.HGet(key, field)
		assert.Equal(t, "OK", result)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HDelete", func(t *testing.T) {
		mock.ExpectHDel(key, field).SetVal(1)
		result, err := instance.HDelete(key, field)
		assert.Equal(t, int64(1), result)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HExists", func(t *testing.T) {
		mock.ExpectHExists(key, field).SetVal(true)
		result, err := instance.HExists(key, field)
		assert.Equal(t, true, result)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("HGetAll", func(t *testing.T) {
		want := map[string]string{"id": "YAAA", "ASD": "YAAA"}
		mock.ExpectHGetAll(key).SetVal(want)
		result, err := instance.HGetAll(key)
		assert.Equal(t, want, result)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func BenchmarkSet(b *testing.B) {
	instance := GetClient()
	key, value := "key", "value"
	for b.Loop() {
		instance.Set(key, value, 10*time.Second)
	}
}
