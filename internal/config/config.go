package config

import (
	"google.golang.org/grpc"
	"os"
	"websocket-service/internal/utils"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	ServerAddress   string `mapstructure:"SERVER_ADDRESS"`
	GrpcServer      string `mapstructure:"GRPC_SERVER"`
	DBDriver        string `mapstructure:"DB_DRIVER"`
	DBSource        string `mapstructure:"DB_SOURCE"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	GrpcClient      *grpc.ClientConn
	RabbitMQUtils   *utils.RabbitMQ
	RabbitMQAddress string `mapstructure:"RABBITMQ_ADDRESS"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return cfg, err
		}

		godotenv.Load()

		viper.SetDefault("SERVER_ADDRESS", os.Getenv("SERVER_ADDRESS"))
		viper.SetDefault("DB_DRIVER", os.Getenv("DB_DRIVER"))
		viper.SetDefault("DB_SOURCE", os.Getenv("DB_SOURCE"))
		viper.SetDefault("JWT_SECRET", os.Getenv("JWT_SECRET"))
		viper.SetDefault("GRPC_SERVER", os.Getenv("GRPC_SERVER"))
		viper.SetDefault("RABBITMQ_ADDRESS", os.Getenv("RABBITMQ_ADDRESS"))
	}

	err = viper.Unmarshal(&cfg)
	return cfg, err
}
