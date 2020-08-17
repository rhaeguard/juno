package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetEnv(t *testing.T) {
	key := "JUNO_RANDOM_KEY_ENV"

	t.Run("Successfully retrieves and trims an existing env. var. value", func(t *testing.T) {
		failIfError(t, os.Setenv(key, "      hello"))
		val := getEnv(key)
		assert.Equal(t, "hello", val)
		t.Cleanup(func() {
			_ = os.Unsetenv(key)
		})
	})

	t.Run("Panics when key does not exist", func(t *testing.T) {
		failIfError(t, os.Unsetenv(key))
		assert.Panics(t, func() {
			getEnv(key)
		})
	})

	t.Run("Panics when key is empty(only spaces)", func(t *testing.T) {
		failIfError(t, os.Setenv(key, "                    "))
		assert.Panics(t, func() {
			getEnv(key)
		})
		t.Cleanup(func() {
			_ = os.Unsetenv(key)
		})
	})
}

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Failed with error : %v", err)
	}
}
