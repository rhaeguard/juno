package config

import (
	"os"
	"strings"
)

// Config is the general application config
var Config = struct {
	ApiVersion    string
	FileUploadDir string
	JwtRealm      string
	JwtSecret     string
	Port          string
}{
	ApiVersion:    "v1",
	FileUploadDir: getEnv("FILE_UPLOAD_DIRECTORY"),
	JwtRealm:      getEnv("JWT_REALM"),
	JwtSecret:     getEnv("JWT_SECRET"),
	Port:          getEnv("APPLICATION_PORT"),
}

// DatabaseConfig is the database specific config
var DatabaseConfig = struct {
	Username,
	Password,
	Endpoint,
	Database string
}{
	Username: getEnv("DATABASE_USERNAME"),
	Password: getEnv("DATABASE_PASSWORD"),
	Endpoint: getEnv("DATABASE_ENDPOINT"),
	Database: getEnv("DATABASE_NAME"),
}

func getEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic("Missing value for the key : " + key)
	}
	return value
}
