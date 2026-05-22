package osqueryplugin

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// daemon supervises a private osqueryd subprocess that backs the plugin.
// It is intentionally minimal: spawn, wait for the extensions socket to
// appear, log stdout/stderr, and tear down cleanly on Stop().
type daemon struct {
	cmd        *exec.Cmd
	socketPath string
	stdoutLog  func(line string)
	stderrLog  func(line string)
}

// resolveBinary returns the absolute path to osqueryd. Lookup order:
//   1. config.binary_path (operator-supplied override),
//   2. embedded binary extracted to the cache dir (when built with
//      `-tags embed_osquery` for a supported platform),
//   3. <dir-of-this-binary>/osquery-bin/osqueryd (install-bundle layout for
//      deployments that prefer to ship the daemon alongside the plugin).
func resolveBinary(override string) (string, error) {
	if override != "" {
		if _, err := os.Stat(override); err != nil {
			return "", fmt.Errorf("osqueryd binary not found at %q: %w", override, err)
		}
		return override, nil
	}
	if path, err := extractEmbedded(); err == nil {
		return path, nil
	} else if !errors.Is(err, ErrNoEmbeddedBinary) {
		return "", fmt.Errorf("extracting embedded osqueryd: %w", err)
	}
	self, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("locating plugin executable: %w", err)
	}
	candidate := filepath.Join(filepath.Dir(self), "osquery-bin", "osqueryd")
	if _, err := os.Stat(candidate); err != nil {
		return "", fmt.Errorf("no osqueryd available: set config.binary_path, build with -tags embed_osquery, or install at %q: %w", candidate, err)
	}
	return candidate, nil
}

// resolveSocketPath returns the unix socket the daemon and client will share.
func resolveSocketPath(override string) string {
	if override != "" {
		return override
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("dtac-osquery-%d.sock", os.Getpid()))
}

// startDaemon spawns osqueryd and waits for its extensions socket to become
// readable. The returned daemon owns the subprocess; the caller must call
// Stop().
func startDaemon(cfg Config, stdoutLog, stderrLog func(line string)) (*daemon, error) {
	bin, err := resolveBinary(cfg.BinaryPath)
	if err != nil {
		return nil, err
	}
	socketPath := resolveSocketPath(cfg.SocketPath)
	// osqueryd refuses to start if the socket file already exists.
	_ = os.Remove(socketPath)

	args := []string{
		"--ephemeral",
		"--disable_database",
		"--disable_logging",
		"--disable_events",
		"--extensions_socket=" + socketPath,
	}
	args = append(args, cfg.ExtraFlags...)

	cmd := exec.Command(bin, args...)
	// Put osqueryd in its own process group so signals we send don't leak
	// to the parent (and so we can SIGKILL the whole group if needed).
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("osqueryd stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("osqueryd stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting osqueryd: %w", err)
	}

	d := &daemon{
		cmd:        cmd,
		socketPath: socketPath,
		stdoutLog:  stdoutLog,
		stderrLog:  stderrLog,
	}
	go d.relay(stdoutPipe, stdoutLog)
	go d.relay(stderrPipe, stderrLog)

	if err := d.waitForSocket(cfg.StartupTimeout); err != nil {
		_ = d.Stop(context.Background())
		return nil, err
	}
	return d, nil
}

func (d *daemon) relay(r io.Reader, sink func(string)) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		if sink != nil {
			sink(scanner.Text())
		}
	}
}

// waitForSocket polls for the extensions socket file to be readable, which
// is the cheapest signal that osqueryd is far enough along to accept clients.
func (d *daemon) waitForSocket(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if _, err := os.Stat(d.socketPath); err == nil {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("osqueryd extensions socket %q did not appear within %s", d.socketPath, timeout)
		}
		// Cheap signal that the process died early.
		if d.cmd.ProcessState != nil && d.cmd.ProcessState.Exited() {
			return errors.New("osqueryd exited during startup")
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// SocketPath returns the unix socket osqueryd is listening on.
func (d *daemon) SocketPath() string { return d.socketPath }

// Stop sends SIGTERM, waits up to ~3s (or until ctx cancels), then SIGKILLs
// the whole process group. Safe to call more than once.
func (d *daemon) Stop(ctx context.Context) error {
	if d == nil || d.cmd == nil || d.cmd.Process == nil {
		return nil
	}
	_ = d.cmd.Process.Signal(syscall.SIGTERM)

	done := make(chan error, 1)
	go func() { done <- d.cmd.Wait() }()

	select {
	case <-done:
	case <-ctx.Done():
	case <-time.After(3 * time.Second):
		// Negative pid signals the process group.
		_ = syscall.Kill(-d.cmd.Process.Pid, syscall.SIGKILL)
		<-done
	}
	_ = os.Remove(d.socketPath)
	return nil
}
