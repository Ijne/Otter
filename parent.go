package main

import (
	"fmt"
	"strconv"
	"syscall"
	"unsafe"
)

func writeProcFile(path string, value string) error {
	pathPtr, _ := syscall.BytePtrFromString(path)
	fdMap, _, errno := syscall.Syscall(
		syscall.SYS_OPEN,
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(syscall.O_WRONLY),
		0,
	)
	defer syscall.Syscall(
		syscall.SYS_CLOSE,
		fdMap,
		0,
		0,
	)
	if errno != 0 {
		return errno
	}
	syscall.Syscall(
		syscall.SYS_WRITE,
		uintptr(fdMap),
		uintptr(unsafe.Pointer(&[]byte(value)[0])),
		uintptr(len(value)),
	)
	return nil
}

func Parent(fdCHW []int32, fdPW []int32, chID uintptr) {
	syscall.Syscall(
		syscall.SYS_CLOSE,
		uintptr(fdCHW[1]),
		0,
		0,
	)

	syscall.Syscall(
		syscall.SYS_CLOSE,
		uintptr(fdPW[0]),
		0,
		0,
	)

	pidStr := strconv.Itoa(int(chID))
	uidMapPath := "/proc/" + pidStr + "/uid_map"
	setgroupsPath := "/proc/" + pidStr + "/setgroups"
	gidMapPath := "/proc/" + pidStr + "/gid_map"
	err := writeProcFile(setgroupsPath, "deny")
	if err != nil {
		fmt.Printf("Error writing to setgroups: %v\n", err)
		return
	}

	uid, _, _ := syscall.Syscall(syscall.SYS_GETUID, 0, 0, 0)
	err = writeProcFile(uidMapPath, fmt.Sprintf("0 %d 1", uid))
	if err != nil {
		fmt.Printf("Error writing to uid_map: %v\n", err)
		return
	}

	gid, _, _ := syscall.Syscall(syscall.SYS_GETGID, 0, 0, 0)
	err = writeProcFile(gidMapPath, fmt.Sprintf("0 %d 1", gid))
	if err != nil {
		fmt.Printf("Error writing to gid_map: %v\n", err)
		return
	}

	syscall.Syscall(
		syscall.SYS_WRITE,
		uintptr(fdPW[1]),
		uintptr(unsafe.Pointer(&[]byte{1}[0])),
		1,
	)

	tmp := []byte{0}
	syscall.Syscall(
		syscall.SYS_READ,
		uintptr(fdCHW[0]),
		uintptr(unsafe.Pointer(&tmp[0])),
		1,
	)
	if tmp[0] == 0 {
		if DEBUG {
			println("Child process started")
		}
	} else {
		println("Child process failed to start")
	}

	syscall.Syscall(
		syscall.SYS_WAIT4,
		chID,
		0,
		0,
	)

	if DEBUG {
		println("Child process finished")
	}

	syscall.Syscall6(
		syscall.SYS_CLOSE,
		uintptr(fdCHW[0]),
		0,
		0,
		0,
		0,
		0,
	)
}
