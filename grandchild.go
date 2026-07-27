package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func Grandchild() {
	source, _ := syscall.BytePtrFromString("proc")
	target, _ := syscall.BytePtrFromString("/proc")
	fstype, _ := syscall.BytePtrFromString("proc")
	_, _, errno := syscall.Syscall6(
		syscall.SYS_MOUNT,
		uintptr(unsafe.Pointer(source)),
		uintptr(unsafe.Pointer(target)),
		uintptr(unsafe.Pointer(fstype)),
		0,
		0,
		0,
	)
	if errno != 0 {
		fmt.Println("Error mounting proc filesystem:", errno)
		return
	}

	path, _ := syscall.BytePtrFromString("/bin/sh")
	argv, _ := syscall.SlicePtrFromStrings([]string{"sh"})
	envp, _ := syscall.SlicePtrFromStrings([]string{"PATH=/bin:usr/bin", "TERM=xterm"})

	syscall.Syscall(
		syscall.SYS_EXECVE,
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(&argv[0])),
		uintptr(unsafe.Pointer(&envp[0])),
	)
}
