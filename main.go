package main

import (
	"os/exec"
	"otter/internal/ns"
	"otter/internal/params"
	"syscall"
)

var (
	DEBUG bool = true
)

func main() {
	cmd := exec.Command("bash", "fsinit.sh")
	cmd.Run()
	ns.Start(params.Params{
		Flags:   uintptr(syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWCGROUP | syscall.SIGCHLD),
		RootFS:  "/tmp/otter",
		Command: []string{"/bin/bash"},
		Env:     []string{"PATH=/bin:/usr/bin:/sbin:/usr/sbin"},
	})
}
