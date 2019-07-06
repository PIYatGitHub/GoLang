package main

import "fmt"

//main.go -- START CONFIG

//PostgresConfig holds the DB config fieldset
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

//Dialect will set the type of db we use, i.e. postgres, mongoDb, ect...
func (c PostgresConfig) Dialect() string {
	return "postgres"
}

//ConnectionInfo is the string used to connect to the db
func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port =%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port =%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

//DefaultPostgresConfiguration is the core value for all configs...
func DefaultPostgresConfiguration() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "DHFywnnrgui237357",
		Name:     "lenslocked_dev",
	}
}

//Config is the default configuration for the environment
type Config struct {
	Port    int    `json:"port"`
	Env     string `json:"env"`
	Pepper  string `json:"pepper"`
	HMACKey string `json:"hmac_key"`
}

//DefaultConfig will get you the port and the environment vars
func DefaultConfig() Config {
	return Config{
		Port:    8080,
		Env:     "dev",
		Pepper:  "secret-random-string",
		HMACKey: "secret-hmac-key",
	}
}

//IsProd will set the production flag
func (c Config) IsProd() bool {
	return c.Env == "prod"
}

//main.go -- END CONFIG

//modes/users.go
// const userPwP = "wrjg82j8#$%^&#Rweg4128y8y8suTO(24#%9ghsdbu"
// const hmacSecretKey = "4wjht8wywr!^Y@$Yggwj8qeyrh139hSFYHEYFehjeo235"

//services.go -- START CONFIG

//services.go -- END CONFIG
