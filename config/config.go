package config

import (
	"fmt"
	"runtime"

	"github.com/spf13/viper"
)

type dbConfig struct {
	Host      string `mapstructure:"host"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DbName    string `mapstructure:"dbname"`
	Port      int    `mapstructure:"port"`
	EnableSsl bool   `mapstructure:"enable_ssl"`
}

type config struct {
	Database  dbConfig `mapstructure:"database"`
	JwtSecret string   `mapstructure:"jwt_secret"`
}

func Load(file string) (*config, error) {
	_, pathname, _, _ := runtime.Caller(0)
	fmt.Println(pathname)

	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	p := pathname + "/../"
	fmt.Println(p)
	viper.AddConfigPath(p)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
