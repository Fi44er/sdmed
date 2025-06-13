package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTPHost            string        `mapstructure:"HTTP_HOST"`
	HTTPPort            string        `mapstructure:"HTTP_PORT"`
	PostgresUrl         string        `mapstructure:"POSTGRES_URL"`
	RedisUrl            string        `mapstructure:"REDIS_URL"`
	VerifyCodeExpiredIn time.Duration `mapstructure:"VERIFY_CODE_EXPIRED_IN"`
	ImageDir            string        `mapstructure:"IMAGE_DIR"`

	CorsOrigin string `mapstructure:"CORS_ORIGIN"`

	PayKeeperUser   string `mapstructure:"PAY_KEEPER_USER"`
	PayKeeperPass   string `mapstructure:"PAY_KEEPER_PASS"`
	PayKeeperServer string `mapstructure:"PAY_KEEPER_SERVER"`

	SMTPTemplatePath string `mapstructure:"SMTP_TEMPLATE_PATH"`
	SMTPHost         string `mapstructure:"SMTP_HOST"`
	SMTPUser         string `mapstructure:"SMTP_USER"`
	SMTPPassword     string `mapstructure:"SMTP_PASSWORD"`
	SMTPFrom         string `mapstructure:"SMTP_FROM"`
	SMTPPort         string `mapstructure:"SMTP_PORT"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAX_AGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAX_AGE"`

	ResetPassTokenExpiredIn time.Duration `mapstructure:"RESET_PASS_TOKEN_EXPIRED_IN"`
	ResetPassURL            string        `mapstructure:"RESET_PASS_URL"`
	ClienUrl                string        `mapstructure:"CLIENT_URL"`
}

func validateConfig(config *Config) error {
	configMap := map[string]interface{}{
		"HTTP_HOST":              config.HTTPHost,
		"HTTP_PORT":              config.HTTPPort,
		"POSTGRES_URL":           config.PostgresUrl,
		"REDIS_URL":              config.RedisUrl,
		"VERIFY_CODE_EXPIRED_IN": config.VerifyCodeExpiredIn,
		"IMAGE_DIR":              config.ImageDir,

		"CORS_ORIGIN": config.CorsOrigin,

		"PAY_KEEPER_USER":   config.PayKeeperUser,
		"PAY_KEEPER_PASS":   config.PayKeeperPass,
		"PAY_KEEPER_SERVER": config.PayKeeperServer,

		"SMTP_TEMPLATE_PATH": config.SMTPTemplatePath,
		"SMTP_HOST":          config.SMTPHost,
		"SMTP_USER":          config.SMTPUser,
		"SMTP_PASSWORD":      config.SMTPPassword,
		"MAIL_FROM":          config.SMTPFrom,
		"MAIL_PORT":          config.SMTPPort,

		"ACCESS_TOKEN_PRIVATE_KEY":  config.AccessTokenPrivateKey,
		"REFRESH_TOKEN_PRIVATE_KEY": config.RefreshTokenPrivateKey,
		"ACCESS_TOKEN_PUBLIC_KEY":   config.AccessTokenPublicKey,
		"REFRESH_TOKEN_PUBLIC_KEY":  config.RefreshTokenPublicKey,
		"ACCESS_TOKEN_EXPIRED_IN":   config.AccessTokenExpiresIn,
		"REFRESH_TOKEN_EXPIRED_IN":  config.RefreshTokenExpiresIn,
		"ACCESS_TOKEN_MAX_AGE":      config.AccessTokenMaxAge,
		"REFRESH_TOKEN_MAX_AGE":     config.RefreshTokenMaxAge,

		"RESET_PASS_TOKEN_EXPIRED_IN": config.ResetPassTokenExpiredIn,
		"RESET_PASS_URL":              config.ResetPassURL,
		"CLIENT_URL":                  config.ClienUrl,
	}

	for key, value := range configMap {
		if isEmptyValue(value) {
			return fmt.Errorf("missing required configuration field: %s", key)
		}
	}

	return nil
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Automatically map environment variables
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	err = validateConfig(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

func isEmptyValue(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case int64:
		return v == 0
	case nil:
		return true
	default:
		return false
	}
}
