package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Camera string `yaml:"camera"`
	Plc    string `yaml:"plc"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")    // Имя конфигурационного файла без расширения
	viper.SetConfigType("yaml")      // Формат конфигурационного файла
	viper.AddConfigPath("./configs") // Путь к директории, где находится конфигурационный файл

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}

	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
		return nil, err
	}

	return &config, nil
}

func (c *Config) String() string {
	return fmt.Sprintf("CameraAddress: %v\n PlcAddress: %v", c.Camera, c.Plc)
}
