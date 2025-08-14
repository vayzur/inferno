package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	InfernoDir = "/etc/inferno"
	ConfigFile = "inferno.yml"
)

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	CertFile string `mapstructure:"certFile" yaml:"certFile"`
	KeyFile  string `mapstructure:"keyFile" yaml:"keyFile"`
}

type Config struct {
	Address                string        `mapstructure:"address" yaml:"address"`
	Port                   uint16        `mapstructure:"port" yaml:"port"`
	Prefork                bool          `mapstructure:"prefork" yaml:"prefork"`
	Token                  string        `mapstructure:"token" yaml:"token"`
	TLS                    TLSConfig     `mapstructure:"tls" yaml:"tls"`
	EtcdEndpoints          []string      `mapstructure:"etcd" yaml:"etcd"`
	NodeMonitorPeriod      time.Duration `mapstructure:"nodeMonitorPeriod" yaml:"nodeMonitorPeriod"`
	NodeMonitorGracePeriod time.Duration `mapstructure:"nodeMonitorGracePeriod" yaml:"nodeMonitorGracePeriod"`
	InboundMonitorPeriod   time.Duration `mapstructure:"inboundMonitorPeriod" yaml:"inboundMonitorPeriod"`
}

var AppConfig *Config

func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&AppConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}
