package config

import "time"

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config/config.yaml"`
	HTTPServer `yaml:"httpServer"`
	AuthJWT    `yaml:"authJwt"`
	CSRFJWT    `yaml:"csrfJwt"`
	GRPC

	Database
	Enviroment     string `env:"ENVIROMENT" env-default:"prod" env-description:"avalible: local, dev, prod"`
	LogFilePath    string `env:"LOG_FILE_PATH" env-default:"zuzu.log"`
	PhotosFilePath string `env:"PHOTOS_FILE_PATH" env-default:"photos/"`
}

func (c Config) GetPhotosFilePath() string {
	return c.PhotosFilePath
}

type HTTPServer struct {
	Address           string        `yaml:"address" yaml-default:"localhost:8080"`
	Timeout           time.Duration `yaml:"timeout" yaml-default:"4s"`
	IdleTimeout       time.Duration `yaml:"idleTimeout" yaml-default:"60s"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout" yaml-default:"10s"`
}

type Database struct {
	DBName string `env:"POSTGRES_DB" env-required:"true"`
	DBPass string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBHost string `env:"DB_HOST" env-default:"0.0.0.0"`
	DBPort int    `env:"DB_PORT" env-required:"true"`
	DBUser string `env:"POSTGRES_USER" env-required:"true"`
}

type AuthJWT struct {
	JwtAccess            string        `env:"AUTH_JWT_SECRET_KEY" env-required:"true"`
	AccessExpirationTime time.Duration `yaml:"accessExpirationTime" yaml-default:"6h"`
	Issuer               string
}

func (a AuthJWT) GetTTL() time.Duration {
	return a.AccessExpirationTime
}
func (a AuthJWT) GetSecret() string {
	return a.JwtAccess
}
func (a AuthJWT) GetIssuer() string {
	return "auth"
}

type GRPC struct {
	AuthPort            int    `env:"GRPC_AUTH_PORT" env-default:"8011"`
	OrderPort           int    `env:"GRPC_ORDER_PORT" env-default:"8012"`
	ProductsPort        int    `env:"GRPC_PRODUCTS_PORT" env-default:"8013"`
	AuthContainerIP     string `env:"GRPC_AUTH_CONTAINER_IP" env-default:"zuzu-auth"`
	OrderContainerIP    string `env:"GRPC_ORDER_CONTAINER_IP" env-default:"zuzu-order"`
	ProductsContainerIP string `env:"GRPC_PRODUCTS_CONTAINER_IP" env-default:"zuzu-products"`
}

type CSRFJWT struct {
	JwtAccess            string        `env:"CSRF_JWT_SECRET_KEY" env-required:"true"`
	AccessExpirationTime time.Duration `yaml:"accessExpirationTime" yaml-default:"6h"`
	Issuer               string
}

func (C CSRFJWT) GetTTL() time.Duration {
	return C.AccessExpirationTime
}

func (C CSRFJWT) GetSecret() string {
	//TODO implement me
	return C.JwtAccess
}

func (C CSRFJWT) GetIssuer() string {
	return C.Issuer
}
