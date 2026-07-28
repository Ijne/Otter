package pkg

import (
	"io"
	"net/http"
	"os"
	"otter/internal/ns"
	"otter/internal/params"
	"path/filepath"
	"syscall"
	"time"
)

type Config struct {
	RootFS      string
	Command     []string
	Env         []string
	MemoryLimit int64
	CPUQuota    int64
	TimeLimit   time.Duration
}

type Result struct {
	ExitCode int
	TimedOut bool
	Stdout   []byte
}

func prepareRootFS(rootfs string) error {
	marker := filepath.Join(rootfs, "bin", "busybox")
	if _, err := os.Stat(marker); err == nil {
		return nil
	}

	dirs := []string{"bin", "proc", "oldroot", "dev", "tmp"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(rootfs, d), 0755); err != nil {
			return err
		}
	}

	busyboxPath := filepath.Join(rootfs, "bin", "busybox")
	if err := downloadFile("https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox", busyboxPath); err != nil {
		return err
	}
	if err := os.Chmod(busyboxPath, 0755); err != nil {
		return err
	}

	links := []string{"sh", "ls", "ps", "mount", "cat", "id", "echo", "clear"}
	for _, name := range links {
		target := filepath.Join(rootfs, "bin", name)
		os.Remove(target)
		if err := os.Symlink("busybox", target); err != nil {
			return err
		}
	}
	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Run(cfg Config) *Result {
	prepareRootFS(cfg.RootFS)

	p := params.Params{
		Flags:       uintptr(syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWCGROUP | syscall.SIGCHLD),
		RootFS:      cfg.RootFS,
		Command:     cfg.Command,
		Env:         cfg.Env,
		MemoryLimit: cfg.MemoryLimit,
		CPUQuota:    cfg.CPUQuota,
		TimeLimit:   cfg.TimeLimit,
	}

	r := ns.Start(p)
	if r == nil {
		return nil
	}

	result := &Result{
		ExitCode: r.ExitCode,
		TimedOut: r.TimedOut,
		Stdout:   r.Stdout,
	}

	return result
}
