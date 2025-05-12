package configs

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"kite/pkg/config"
	"os"
	"strings"
)

type Config struct {
	Server   ServerConfig `mapstructure:"server"`
	Database MySQLConfig  `mapstructure:"database"`
	Log      LogConfig    `mapstructure:"log"`
}

type ServerConfig struct {
	Port            int `mapstructure:"port"`
	ShutdownTimeout int `mapstructure:"shutdown_timeout"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	Encoding string `mapstructure:"encoding"`
	Dev      bool   `mapstructure:"dev"`
}

type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Name        string `mapstructure:"name"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Charset     string `mapstructure:"charset"`
	MaxIdleConn int    `mapstructure:"max_idle_conn"`
	MaxOpenConn int    `mapstructure:"max_open_conn"`
	MaxLifeConn int    `mapstructure:"max_life_conn"`
	SSLMode     string `mapstructure:"ssl_mode"`
}

func Load() (*Config, error) {
	v := viper.New()

	configPath := getConfigPath()
	v.AddConfigPath(configPath)
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// 读取基础配置
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	// 读取环境相关配置
	environment := config.GetEnvironment()
	envConfigName := fmt.Sprintf("config.%s", environment)
	v.SetConfigName(envConfigName)

	// 合并基础配置和环境相关配置
	if err := v.MergeInConfig(); err != nil {
		// 文件如果不存在，则忽略错误，否则提示合并失败
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to merge enviromment config: %w", err)
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.AllowEmptyEnv(true)

	// 反序列化配置
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	// 验证配置
	if err := validateConfig(&c); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}
	return &c, nil
}

func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return "./configs"
}
