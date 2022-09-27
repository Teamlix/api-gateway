package config

type Config struct {
	App  struct{} `yaml:"app"`
	Http struct {
		Host string `yaml:"host" env:"HTTP_HOST" env-description:"Http host"`
		Port string `yaml:"port" env:"HTTP_PORT" env-description:"Http port"`
	} `yaml:"http"`
	Grpc struct {
		Client struct {
			User string `yaml:"user" env:"GRPC_CLIENT_USER_SERVICE_URL" env-description:"gRPC client for user-service"`
		} `yaml:"client"`
	} `yaml:"grpc"`
}
