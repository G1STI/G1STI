package singbox

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

var ErrAlreadyRunning = errors.New("sing-box process already running")
var ErrNotRunning = errors.New("sing-box process not running")

type Manager struct {
	mu         sync.Mutex
	cmd        *exec.Cmd
	binaryPath string
	configPath string
	logPath    string
}

func NewManager(binaryPath, configPath, logPath string) *Manager {
	return &Manager{
		binaryPath: binaryPath,
		configPath: configPath,
		logPath:    logPath,
	}
}

func (m *Manager) IsRunning() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cmd != nil && m.cmd.Process != nil
}

func (m *Manager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd != nil && m.cmd.Process != nil {
		return ErrAlreadyRunning
	}
	if m.binaryPath == "" {
		return fmt.Errorf("sing-box binary path is empty")
	}

	if err := os.MkdirAll(filepath.Dir(m.logPath), 0o755); err != nil {
		return fmt.Errorf("create sing-box log dir: %w", err)
	}

	logFile, err := os.OpenFile(m.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("open sing-box log file: %w", err)
	}

	cmd := exec.CommandContext(ctx, m.binaryPath, "run", "-c", m.configPath)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	cmd.SysProcAttr = nil

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start sing-box: %w", err)
	}

	m.cmd = cmd
	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cmd == nil || m.cmd.Process == nil {
		return ErrNotRunning
	}

	if err := m.cmd.Process.Kill(); err != nil {
		return fmt.Errorf("stop sing-box: %w", err)
	}
	_, _ = m.cmd.Process.Wait()
	m.cmd = nil
	return nil
}

func (m *Manager) Restart(ctx context.Context) error {
	if m.IsRunning() {
		if err := m.Stop(); err != nil {
			return err
		}
		// Allow process to release resources.
		time.Sleep(250 * time.Millisecond)
	}
	return m.Start(ctx)
}

func (m *Manager) BinaryPath() string {
	return m.binaryPath
}

func (m *Manager) ConfigPath() string {
	return m.configPath
}

func (m *Manager) LogPath() string {
	return m.logPath
}
