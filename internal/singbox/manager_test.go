package singbox

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestManagerStartStop(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("mock script uses sh")
	}

	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "mock-singbox.sh")
	configPath := filepath.Join(tmpDir, "config.json")
	logPath := filepath.Join(tmpDir, "singbox.log")

	script := "#!/bin/sh\ntrap 'exit 0' TERM INT\nsleep 60\n"
	if err := os.WriteFile(binaryPath, []byte(script), 0o755); err != nil {
		t.Fatalf("write mock binary: %v", err)
	}
	if err := os.WriteFile(configPath, []byte("{}"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	manager := NewManager(binaryPath, configPath, logPath)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := manager.Start(ctx); err != nil {
		t.Fatalf("start: %v", err)
	}

	time.Sleep(50 * time.Millisecond)
	if !manager.IsRunning() {
		t.Fatalf("expected running manager")
	}

	if err := manager.Stop(); err != nil {
		t.Fatalf("stop: %v", err)
	}

	if manager.IsRunning() {
		t.Fatalf("expected manager to stop")
	}
}
