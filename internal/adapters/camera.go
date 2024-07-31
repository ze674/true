package adapters

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

var rwTimeout = 1 * time.Second   // Таймаут для чтения/записи
var dialTimeout = 5 * time.Second // Таймаут для подключения

var retryDelay = 2 * time.Second // Задержка между повторными попытками

type Camera struct {
	Address string
	Conn    net.Conn
	Logger  *zap.Logger
}

// NewCamera создает новый экземпляр камеры
func NewСamera(address string, logger *zap.Logger) *Camera {
	return &Camera{
		Address: address,
		Logger:  logger,
	}
}

// Connect пытается подключиться к камере
func (c *Camera) Connect() error {
	op := fmt.Sprintf("Camera.Connect: %s", c.Address)

	c.Logger.Info("Attempting to connect to camera", zap.String("operation", op))
	conn, err := net.DialTimeout("tcp", c.Address, dialTimeout)
	if err != nil {
		c.Logger.Error("Failed to connect to camera", zap.String("operation", op), zap.Error(err))
		return fmt.Errorf("operation: %s, error: %w", op, err)
	}

	c.Conn = conn
	c.Logger.Info("Successfully connected to camera", zap.String("operation", op))
	return nil
}

func (c *Camera) Reconnect() error {
	c.Close()
	for {
		err := c.Connect()
		if err == nil {
			return nil
		}
		c.Logger.Error("Reconnect failed, retrying...", zap.Error(err))
		time.Sleep(retryDelay)
	}
}

// Close closes the connection to the camera
func (c *Camera) Close() error {
	if c.Conn != nil {
		c.Logger.Info("Closing connection to camera", zap.String("address", c.Address))
		err := c.Conn.Close()
		if err != nil {
			c.Logger.Error("Failed to close connection to camera", zap.String("address", c.Address), zap.Error(err))
			return err
		}
		c.Logger.Info("Successfully closed connection to camera", zap.String("address", c.Address))
	}
	return nil
}

// SendCommand sends a command to the camera
func (c *Camera) SendCommand(cmd string) error {
	if c.Conn == nil {
		err := fmt.Errorf("no active connection")
		c.Logger.Error("Failed to send command", zap.String("command", cmd), zap.Error(err))
		return err
	}

	cmd = cmd + "\n"
	_, err := c.Conn.Write([]byte(cmd))
	if err != nil {
		c.Logger.Error("Failed to send command", zap.String("command", cmd), zap.Error(err))
		return fmt.Errorf("failed to send command '%s': %w", cmd, err)
	}

	c.Logger.Info("Command sent successfully", zap.String("command", cmd))
	return nil
}

// Read reads data from the camera connection
func (c *Camera) Read() (string, error) {
	if c.Conn == nil {
		err := fmt.Errorf("no active connection")
		c.Logger.Error("Failed to read data", zap.Error(err))
		return "", err
	}

	scanner := bufio.NewScanner(c.Conn)
	c.Conn.SetReadDeadline(time.Now().Add(rwTimeout))

	if scanner.Scan() {
		text := scanner.Text()
		c.Logger.Info("Data read from camera", zap.String("data", text))
		return text, nil

	}

	err := scanner.Err()
	if err != nil {
		c.Logger.Error("Failed to read data", zap.Error(err))
		return "", err
	}

	c.Logger.Info("No data read from camera")
	return "", nil
}
