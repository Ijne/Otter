package main

import (
	"fmt"
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

	targetPath, _ := syscall.BytePtrFromString("/")
	syscall.Syscall6(
		syscall.SYS_MOUNT,
		0,
		uintptr(unsafe.Pointer(targetPath)),
		0,
		uintptr(syscall.MS_PRIVATE|syscall.MS_REC),
		0,
		0,
	)

	rootFSPath, _ := syscall.BytePtrFromString("/tmp/rootfs")
	syscall.Syscall6(
		syscall.SYS_MOUNT,
		uintptr(unsafe.Pointer(rootFSPath)),
		uintptr(unsafe.Pointer(rootFSPath)),
		0,
		uintptr(syscall.MS_BIND|syscall.MS_REC),
		0,
		0,
	)

	oldRootPath, _ := syscall.BytePtrFromString("/tmp/rootfs/oldroot")
	syscall.Syscall(
		syscall.SYS_PIVOT_ROOT,
		uintptr(unsafe.Pointer(rootFSPath)),
		uintptr(unsafe.Pointer(oldRootPath)),
		0,
	)

	syscall.Syscall(
		syscall.SYS_CHDIR,
		uintptr(unsafe.Pointer(targetPath)),
		0,
		0,
	)

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

	oldRootPath, _ = syscall.BytePtrFromString("/oldroot")
	syscall.Syscall(
		syscall.SYS_UMOUNT2,
		uintptr(unsafe.Pointer(oldRootPath)),
		uintptr(syscall.MNT_DETACH),
		0,
	)

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
