package paths

import (
	"fmt"
	"os"
	"path/filepath"
)

const AppName = "G1STI"

func RootDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}
	return filepath.Join(base, AppName), nil
}

func LogsDir() (string, error) {
	root, err := RootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "logs"), nil
}

func SingBoxConfigPath() (string, error) {
	root, err := RootDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "sing-box.json"), nil
}

func AppLogPath() (string, error) {
	logs, err := LogsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(logs, "app.log"), nil
}

func SingBoxLogPath() (string, error) {
	logs, err := LogsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(logs, "singbox.log"), nil
}
