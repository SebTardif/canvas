package browser

import (
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"
)

func TestWaitProcessPreservesExitError(t *testing.T) {
	cmd := exec.Command("sh", "-c", "exit 7")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	var exitErr *exec.ExitError
	if err := waitProcess(cmd, 5*time.Second); !errors.As(err, &exitErr) || exitErr.ExitCode() != 7 {
		t.Fatalf("expected exit code 7, got %v", err)
	}
}

func TestWaitProcessTimeoutReapsProcess(t *testing.T) {
	cmd := exec.Command("sleep", "60")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	if err := waitProcess(cmd, 10*time.Millisecond); err != context.DeadlineExceeded {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
	if cmd.ProcessState == nil {
		t.Fatal("timed-out process was not reaped")
	}
}
