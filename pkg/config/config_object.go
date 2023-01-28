package config

type ConfigObject struct {
	AppHost string `mapstructure:"APP_HOST"`
	AppPort string `mapstructure:"APP_PORT"`
	AppName string `mapstructure:"APP_NAME"`

	//	Redis
	RedisHost                  string `mapstructure:"REDIS_HOST"`
	RedisPort                  string `mapstructure:"REDIS_PORT"`
	RedisDialConnectTimeout    string `mapstructure:"REDIS_DIAL_CONNECT_TIMEOUT"`
	RedisReadTimeout           string `mapstructure:"REDIS_READ_TIMEOUT"`
	RedisWriteTimeout          string `mapstructure:"REDIS_WRITE_TIMEOUT"`
	RedisIdleTimeout           string `mapstructure:"REDIS_IDLE_TIMEOUT"`
	RedisConnLifetimeMax       string `mapstructure:"REDIS_CONN_LIFETIME_MAX"`
	RedisConnIdleMax           string `mapstructure:"REDIS_CONN_IDLE_MAX"`
	RedisConnActiveMax         string `mapstructure:"REDIS_CONN_ACTIVE_MAX"`
	RedisIsWait                string `mapstructure:"REDIS_IS_WAIT"`
	RedisNamespace             string `mapstructure:"REDIS_NAMESPACE"`
	RedisPassword              string `mapstructure:"REDIS_PASSWORD"`
	RedisLockerTries           string `mapstructure:"REDIS_LOCKER_TRIES"`
	RedisLockerTriesRetryDelay string `mapstructure:"REDIS_LOCKER_TRIES_RETRY_DELAY"`
	RedisLockerExpiry          string `mapstructure:"REDIS_LOCKER_EXPIRY"`
}
