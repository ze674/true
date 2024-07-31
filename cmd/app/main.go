package main

import (
	"log"
	"time"

	"true/internal/adapters"
	"true/internal/config"

	"go.uber.org/zap"
)

func main() {
	// Настройка логгера
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	logger.Debug("Loaded config", zap.String("config", cfg.String()))

	// Инициализация адаптера камеры
	camera := adapters.NewСamera(cfg.Camera, logger)
	logger.Debug("Initializing camera adapter", zap.String("address", cfg.Camera))

	// Подключение к камере
	logger.Debug("Connecting to camera", zap.String("address", cfg.Camera))
	if err := camera.Connect(); err != nil {
		logger.Fatal("Failed to connect to camera", zap.Error(err))
	}
	defer func() {
		if err := camera.Close(); err != nil {
			logger.Error("Failed to close camera connection", zap.Error(err))
		}
	}()

	logger.Debug("Camera connected", zap.String("address", cfg.Camera))

	// Отправка команды на камеру
	logger.Debug("Sending command to camera")
	if err := camera.SendCommand(" "); err != nil {
		logger.Error("Failed to send command", zap.Error(err))
	}

	// Чтение ответа от камеры
	logger.Debug("Reading response from camera")
	response, err := camera.Read()
	if err != nil {
		logger.Error("Failed to read response from camera", zap.Error(err))
	} else {
		logger.Info("Received response from camera", zap.String("response", response))
	}

	// Ожидание для демонстрации работы программы
	time.Sleep(5 * time.Second)
}
