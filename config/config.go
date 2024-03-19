package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type config struct {
	HttpPort string `mapstructure:"HTTP_PORT"`
	Env      string `mapstructure:"ENV"`

	Database   database   `mapstructure:",squash"`
	Service    service    `mapstructure:",squash"`
	JWT        jwt        `mapstructure:",squash"`
	S3         s3         `mapstructure:",squash"`
	Otel       otel       `mapstructure:",squash"`
	Prometheus prometheus `mapstructure:",squash"`
	Bcrypt     bcrypt     `mapstructure:",squash"`
}

type service struct {
	Timeout int    `mapstructure:"SERVICE_TIMEOUT"`
	Name    string `mapstructure:"SERVICE_NAME"`
	Version string `mapstructure:"SERVICE_VERSION"`
}

type database struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
	SSLMode  string `mapstructure:"DB_SSL_MODE"`
}

func (d database) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", d.Host, d.Port, d.Username, d.Password, d.Name)
}

type jwt struct {
	Secret   string `mapstructure:"JWT_SECRET"`
	ExpireIn int    `mapstructure:"JWT_EXPIRE_IN"`
}

type otel struct {
	ExporterOTLPEndpoint   string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	Insecure               bool   `mapstructure:"OTEL_INSECURE"`
	ExporterPrometheusPath string `mapstructure:"OTEL_EXPORTER_PROMETHEUS_PATH"`
	ExporterPrometheusPort string `mapstructure:"OTEL_EXPORTER_PROMETHEUS_PORT"`
	EnableTracing          bool   `mapstructure:"OTEL_ENABLE_TRACING"`
	EnableMetrics          bool   `mapstructure:"OTEL_ENABLE_METRICS"`
	EnableLogging          bool   `mapstructure:"OTEL_ENABLE_LOGGING"`
	OnlyPrometheusExporter bool   `mapstructure:"OTEL_ONLY_PROMETHEUS_EXPORTER"`
}

type prometheus struct {
	Address string `mapstructure:"PROMETHEUS_ADDRESS"`
}

type bcrypt struct {
	Salt int `mapstructure:"BCRYPT_SALT"`
}

type s3 struct {
	EndPoint  string `mapstructure:"S3_BASE_URL"`
	AccessKey string `mapstructure:"S3_ID"`
	SecretKey string `mapstructure:"S3_SECRET_KEY"`
	Bucket    string `mapstructure:"S3_BUCKET_NAME"`
	UseSSL    bool   `mapstructure:"S3_USESSL"`
	Region    string `mapstructure:"S3_REGION"`
}

var configInstance *config
var viperInstance *viper.Viper

func LoadConfig(filenames ...string) (*viper.Viper, error) {
	if viperInstance != nil {
		return viperInstance, nil
	}
	v := viper.New()
	if len(filenames) > 0 {
		// v.SetConfigName("app")
		v.SetConfigFile(filenames[0])
	} else {
		// check .env file exist
		if _, err := os.Stat(".env"); err == nil {
			v.SetConfigFile(".env")
		}
	}

	initDefaultValue(v)
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil && !strings.Contains(err.Error(), "Not Found in") {
		err = fmt.Errorf("error read config file: %s", err)
		return nil, err
	}

	viperInstance = v
	return viperInstance, nil
}

func ParseConfig(v *viper.Viper) (*config, error) {
	if configInstance != nil {
		return configInstance, nil
	}
	var c config
	var out map[string]interface{}
	err := mapstructure.Decode(&c, &out)
	if err != nil {
		err = fmt.Errorf("error decode config: %s", err)
		return nil, err
	}

	for key := range out {
		vKey := strings.ToLower(strings.ReplaceAll(key, ".", "_"))
		err = v.BindEnv(vKey, key)
		if err != nil {
			err = fmt.Errorf("error bind env: %s", err)
			return nil, err
		}
	}

	err = v.Unmarshal(&c)
	if err != nil {
		err = fmt.Errorf("error unmarshal config: %s", err)
		return nil, err
	}

	configInstance = &c
	return configInstance, nil
}

func Get(filenames ...string) *config {
	if configInstance == nil {
		LoadConfig(filenames...)
		ParseConfig(viperInstance)
	}
	return configInstance
}

func GetViper(filenames ...string) *viper.Viper {
	if viperInstance == nil {
		LoadConfig(filenames...)
		ParseConfig(viperInstance)
	}
	return viperInstance
}

func initDefaultValue(v *viper.Viper) {
	v.SetDefault("HTTP_PORT", "8000")
	v.SetDefault("ENV", "dev")
	v.SetDefault("SERVICE_NAME", "project-sprint-social-media-api")
	v.SetDefault("SERVICE_TIMEOUT", 30)
	v.SetDefault("OTEL_INSECURE", true)
	v.SetDefault("OTEL_EXPORTER_PROMETHEUS_PATH", "/metrics")
	v.SetDefault("OTEL_EXPORTER_PROMETHEUS_PORT", "2223")
	v.SetDefault("JWT_SECRET", "secret")
	v.SetDefault("JWT_EXPIRE_IN", 120)
	v.SetDefault("S3_REGION", "ap-southeast-1")
	v.SetDefault("OTEL_ENABLE_METRICS", true)
	v.SetDefault("OTEL_ONLY_PROMETHEUS_EXPORTER", true)
}
