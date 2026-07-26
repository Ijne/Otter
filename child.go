package main

import (
	"syscall"
	"unsafe"
)

func Child(fdCHW []int32, fdPW []int32) {
	syscall.Syscall(
		syscall.SYS_CLOSE,
		uintptr(fdCHW[0]),
		0,
		0,
	)

	syscall.Syscall(
		syscall.SYS_CLOSE,
		uintptr(fdPW[1]),
		0,
		0,
	)
	tmp := []byte{0}
	syscall.Syscall(
		syscall.SYS_READ,
		uintptr(fdPW[0]),
		uintptr(unsafe.Pointer(&tmp[0])),
		1,
	)

	if DEBUG {
		println("Child process received signal to start")
	}

	path, _ := syscall.BytePtrFromString("/bin/sh")
	argv, _ := syscall.SlicePtrFromStrings([]string{"sh"})
	envp, _ := syscall.SlicePtrFromStrings([]string{"PATH=/bin:usr/bin", "TERM=xterm"})
	_, _, errno := syscall.Syscall(
		syscall.SYS_EXECVE,
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(&argv[0])),
		uintptr(unsafe.Pointer(&envp[0])),
	)
	if errno != 0 {
		tmp := []byte{1}
		syscall.Syscall(
			syscall.SYS_WRITE,
			uintptr(fdCHW[1]),
			uintptr(unsafe.Pointer(&tmp[0])),
			1,
		)
	}
}
