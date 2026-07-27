package main

import (
	"os/exec"
	"syscall"
	"unsafe"
)

var (
	DEBUG bool = true
)

func main() {
	cmd := exec.Command("bash", "fsinit.sh")
	cmd.Run()

	fdCHW := make([]int32, 2)
	_, _, errno := syscall.Syscall(
		syscall.SYS_PIPE,
		uintptr(unsafe.Pointer(&fdCHW[0])),
		0,
		0,
	)
	if errno != 0 {
		println("Error:", errno)
		return
	}

	_, _, errno = syscall.Syscall(
		syscall.SYS_FCNTL,
		uintptr(fdCHW[1]),
		uintptr(syscall.F_SETFD),
		uintptr(syscall.FD_CLOEXEC),
	)

	fdPW := make([]int32, 2)
	_, _, errno = syscall.Syscall(
		syscall.SYS_PIPE,
		uintptr(unsafe.Pointer(&fdPW[0])),
		0,
		0,
	)
	if errno != 0 {
		println("Error:", errno)
		return
	}

	r1, _, errno := syscall.Syscall(
		syscall.SYS_CLONE,
		uintptr(syscall.CLONE_NEWUSER|syscall.CLONE_NEWUTS|syscall.CLONE_NEWNS|syscall.CLONE_NEWPID),
		0,
		0,
	)
	if errno != 0 {
		println("Error:", errno)
		return
	}

	if r1 == 0 {
		Child(fdCHW, fdPW)
	} else {
		Parent(fdCHW, fdPW, r1)
	}
}
