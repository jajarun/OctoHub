package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	// 服务器配置
	Server ServerConfig `mapstructure:"server"`

	// 签名配置
	Signature SignatureConfig `mapstructure:"signature"`

	// WebSocket配置
	WebSocket WebSocketConfig `mapstructure:"websocket"`

	// 连接配置
	Connection ConnectionConfig `mapstructure:"connection"`

	// 日志配置
	Logging LoggingConfig `mapstructure:"logging"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// SignatureConfig 签名配置
type SignatureConfig struct {
	Key     string `mapstructure:"key"`
	Timeout int    `mapstructure:"timeout"` // 秒
}

// WebSocketConfig WebSocket配置
type WebSocketConfig struct {
	ReadTimeout    int   `mapstructure:"read_timeout"`     // 秒
	WriteTimeout   int   `mapstructure:"write_timeout"`    // 秒
	MaxMessageSize int64 `mapstructure:"max_message_size"` // 字节
}

// ConnectionConfig 连接配置
type ConnectionConfig struct {
	MaxConnections  int `mapstructure:"max_connections"`
	BufferSize      int `mapstructure:"buffer_size"`
	CleanupInterval int `mapstructure:"cleanup_interval"` // 清理检查间隔（秒）
	StaleTimeout    int `mapstructure:"stale_timeout"`    // 连接过期时间（秒）
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// LoadConfig 加载配置
func LoadConfig() *Config {
	v := viper.New()

	// 直接设置配置文件
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".") // 只在当前目录查找 config.yaml

	// 设置环境变量前缀
	v.SetEnvPrefix("OCTOHUB")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 设置默认值
	setDefaults(v)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("警告: 配置文件未找到，使用默认配置: %v", err)
		} else {
			log.Printf("警告: 读取配置文件失败，使用默认配置: %v", err)
		}
	} else {
		log.Printf("使用配置文件: %s", v.ConfigFileUsed())
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("无法解析配置: %v", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	return &config
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 服务器默认配置
	v.SetDefault("server.port", "8080")

	// 签名默认配置
	v.SetDefault("signature.key", "your-secret-key-here")
	v.SetDefault("signature.timeout", 300) // 5分钟

	// WebSocket默认配置
	v.SetDefault("websocket.read_timeout", 60)          // 60秒
	v.SetDefault("websocket.write_timeout", 10)         // 10秒
	v.SetDefault("websocket.max_message_size", 1048576) // 1MB

	// 连接默认配置
	v.SetDefault("connection.max_connections", 10000)
	v.SetDefault("connection.buffer_size", 256)
	v.SetDefault("connection.cleanup_interval", 30) // 30秒
	v.SetDefault("connection.stale_timeout", 120)   // 2分钟

	// 日志默认配置
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "text")
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Server.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}

	if config.Signature.Key == "" {
		return fmt.Errorf("签名密钥不能为空")
	}

	if config.Signature.Timeout <= 0 {
		return fmt.Errorf("签名超时时间必须大于0")
	}

	if config.WebSocket.ReadTimeout <= 0 {
		return fmt.Errorf("WebSocket读取超时时间必须大于0")
	}

	if config.WebSocket.WriteTimeout <= 0 {
		return fmt.Errorf("WebSocket写入超时时间必须大于0")
	}

	if config.Connection.MaxConnections <= 0 {
		return fmt.Errorf("最大连接数必须大于0")
	}

	if config.Connection.BufferSize <= 0 {
		return fmt.Errorf("缓冲区大小必须大于0")
	}

	return nil
}

// GetSignatureTimeout 获取签名超时时间
func (c *Config) GetSignatureTimeout() time.Duration {
	return time.Duration(c.Signature.Timeout) * time.Second
}

// GetReadTimeout 获取读取超时时间
func (c *Config) GetReadTimeout() time.Duration {
	return time.Duration(c.WebSocket.ReadTimeout) * time.Second
}

// GetWriteTimeout 获取写入超时时间
func (c *Config) GetWriteTimeout() time.Duration {
	return time.Duration(c.WebSocket.WriteTimeout) * time.Second
}

// IsDefaultSignatureKey 检查是否使用默认签名密钥
func (c *Config) IsDefaultSignatureKey() bool {
	return c.Signature.Key == "your-secret-key-here"
}
