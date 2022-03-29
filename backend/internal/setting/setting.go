package setting

import (
	"os"

	"github.com/ZmaximillianZ/local-chain/internal/db"
)

// App is a structure for storage app configuration
type App struct {
	JwtSecret string

	RuntimeRootPath string

	ExportSavePath string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

// DBSetting is a structure for storage db configuration
type DBSetting struct {
	Database         string
	Username         string
	Password         string
	PostgresPassword string
	URL              string
	MaxIdleCons      string
	MaxOpenCons      string
}

// ServerSetting is a structure for storage user_protobuf configuration
type ServerSetting struct {
	RunMode string
	Host    string
	Port    string
	// ReadTimeout  time.Duration
	// WriteTimeout time.Duration
	// Path string
}

// NodeSetting is a structure for storage blockchain configuration
type NodeSetting struct {
	Host string
	Port string
}

type Setting struct {
	ServerConfig ServerSetting
	NodeConfig   NodeSetting
	DBConfig     db.ConnectionSettions
	App          App
}

// LoadSetting loads configuration from env variables
func LoadSetting() *Setting {
	// TODO: Try use go-env for easy unmarshalling https://github.com/Netflix/go-env
	return &Setting{
		ServerConfig: ServerSetting{
			Host: getEnv("HOST"),
			Port: getEnv("PORT"),
		},
		NodeConfig: NodeSetting{
			Host: getEnv("NODE_SERVICE"),
			Port: getEnv("ETH_PORT"),
		},
		DBConfig: db.ConnectionSettions{
			Database:    "postgres",
			URL:         getEnv("DATABASE_URL"), // for production DATABASE_CONTAINER_URL
			MaxIdleCons: 100,
			MaxOpenCons: 10,
		},
		App: App{
			getEnv("JWT_SECRET"),
			getEnv("ROOT_PATH"),
			getEnv("SAVE_PATH"),
			getEnv("LOG_PATH"),
			getEnv("LOG_NAME"),
			getEnv("LOG_EXT"),
			getEnv(""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return ""
}
