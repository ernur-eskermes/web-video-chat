package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultHTTPPort               = "8000"
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1
	defaultAccessTokenTTL         = 15 * time.Minute
	defaultRefreshTokenTTL        = 24 * time.Hour * 30
	defaultLimiterRPS             = 10
	defaultLimiterBurst           = 2
	defaultLimiterTTL             = 10 * time.Minute

	EnvLocal = "local"
	Prod     = "prod"
)

type (
	Config struct {
		Environment string
		Mongo       MongoConfig
		HTTP        HTTPConfig
		Auth        AuthConfig
		LiveKit     LiveKitConfig
		Limiter     LimiterConfig
		GoogleOauth GoogleOauthConfig
	}

	GoogleOauthConfig struct {
		ClientId     string
		ClientSecret string
		CallbackURL  string
	}

	MongoConfig struct {
		URI      string
		User     string
		Password string
		Name     string `mapstructure:"databaseName"`
	}

	AuthConfig struct {
		JWT           JWTConfig
		SessionSecret string
		PasswordSalt  string
	}

	LiveKitConfig struct {
		Host      string
		ApiKey    string
		ApiSecret string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	HTTPConfig struct {
		Schema             string
		Host               string
		Port               string        `mapstructure:"port"`
		ReadTimeout        time.Duration `mapstructure:"readTimeout"`
		WriteTimeout       time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderMegabytes int           `mapstructure:"maxHeaderBytes"`
	}

	LimiterConfig struct {
		RPS   int
		Burst int
		TTL   time.Duration
	}
)

// Init populates Config struct with values from config file
// located at filepath and environment variables.
func Init(configsDir string) (*Config, error) {
	populateDefaults()

	if err := parseConfigFile(configsDir, os.Getenv("APP_ENV")); err != nil {
		return nil, err
	}

	var cfg Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}

	setFromEnv(&cfg)

	return &cfg, nil
}

func unmarshal(cfg *Config) error {
	if err := viper.UnmarshalKey("mongo", &cfg.Mongo); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("rest", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}

	return viper.UnmarshalKey("limiter", &cfg.Limiter)
}

func setFromEnv(cfg *Config) {
	cfg.GoogleOauth.ClientId = os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	cfg.GoogleOauth.ClientSecret = os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	cfg.GoogleOauth.CallbackURL = os.Getenv("GOOGLE_OAUTH_CALLBACK_URL")

	cfg.LiveKit.Host = os.Getenv("LIVEKIT_HOST")
	cfg.LiveKit.ApiKey = os.Getenv("LIVEKIT_APIKEY")
	cfg.LiveKit.ApiSecret = os.Getenv("LIVEKIT_APISECRET")

	cfg.Mongo.URI = os.Getenv("MONGO_URI")
	cfg.Mongo.User = os.Getenv("MONGO_USER")
	cfg.Mongo.Password = os.Getenv("MONGO_PASS")

	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")
	cfg.Auth.SessionSecret = os.Getenv("SESSION_SECRET")

	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.HTTP.Schema = os.Getenv("HTTP_SCHEMA")

	cfg.Environment = os.Getenv("APP_ENV")
}

func parseConfigFile(folder, env string) error {
	viper.AddConfigPath(folder)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if env == EnvLocal {
		return nil
	}

	viper.SetConfigName(env)

	return viper.MergeInConfig()
}

func populateDefaults() {
	viper.SetDefault("rest.port", defaultHTTPPort)
	viper.SetDefault("rest.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("rest.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("rest.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
	viper.SetDefault("limiter.rps", defaultLimiterRPS)
	viper.SetDefault("limiter.burst", defaultLimiterBurst)
	viper.SetDefault("limiter.ttl", defaultLimiterTTL)
}
