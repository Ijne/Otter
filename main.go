package main

import (
	"os/exec"
	"otter/internal/ns"
	"otter/internal/params"
	"syscall"
	"time"
)

var (
	DEBUG bool = true
)

func main() {
	config := params.Params{
		Flags:       uintptr(syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWCGROUP | syscall.SIGCHLD),
		RootFS:      "/tmp/otter",
		Command:     []string{"/bin/sh"},
		Env:         []string{"PATH=/bin:/usr/bin:/sbin:/usr/sbin"},
		MemoryLimit: 0,
		CPUQuota:    0,
		TimeLimit:   5 * time.Second,
	}
	cmd := exec.Command("bash", "fsinit.sh", config.RootFS)
	cmd.Run()
	ns.Start(config)
}
