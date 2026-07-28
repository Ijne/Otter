package ns

import (
	"fmt"
	"otter/internal/mount"
	"otter/internal/params"
	"strconv"
	"syscall"
	"unsafe"
)

func Start(p params.Params) {
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
		p.Flags,
		0,
		0,
	)
	if errno != 0 {
		println("Error:", errno)
		return
	}

	if r1 == 0 {
		childEntry(p, fdCHW, fdPW)
	} else {
		parentEntry(fdCHW, fdPW, r1)
	}
}

func childEntry(p params.Params, fdCHW, fdPW []int32) {
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

	mount.PivotRoot(p)

	mount.MountProc()

	mount.Cleanup()

	execInto(p)
}

func execInto(p params.Params) {
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

func parentEntry(fdCHW, fdPW []int32, chID uintptr) {
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

	syscall.Syscall(syscall.SYS_WAIT4, chID, 0, 0)

	syscall.Syscall6(syscall.SYS_CLOSE, uintptr(fdCHW[0]), 0, 0, 0, 0, 0)
}
