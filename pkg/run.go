package pkg

import (
	"os/exec"
	"otter/internal/ns"
	"otter/internal/params"
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

func Run(cfg Config) *Result {
	cmd := exec.Command("bash", "fsinit.sh", cfg.RootFS)
	cmd.Run()

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
