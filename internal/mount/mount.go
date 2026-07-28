package mount

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/Ijne/Otter/internal/params"
)

func PivotRoot(p params.Params) {
	targetPath, _ := syscall.BytePtrFromString("/")
	_, _, errno := syscall.Syscall6(
		syscall.SYS_MOUNT,
		0,
		uintptr(unsafe.Pointer(targetPath)),
		0,
		uintptr(syscall.MS_PRIVATE|syscall.MS_REC),
		0,
		0,
	)
	if errno != 0 {
		fmt.Println("Error mounting root filesystem:", errno)
		return
	}

	rootFSPath, _ := syscall.BytePtrFromString(p.RootFS)
	_, _, errno = syscall.Syscall6(
		syscall.SYS_MOUNT,
		uintptr(unsafe.Pointer(rootFSPath)),
		uintptr(unsafe.Pointer(rootFSPath)),
		0,
		uintptr(syscall.MS_BIND|syscall.MS_REC),
		0,
		0,
	)
	if errno != 0 {
		fmt.Println("Error mounting root filesystem:", errno)
		return
	}

	oldRootPath, _ := syscall.BytePtrFromString(p.RootFS + "/oldroot")
	_, _, errno = syscall.Syscall(
		syscall.SYS_PIVOT_ROOT,
		uintptr(unsafe.Pointer(rootFSPath)),
		uintptr(unsafe.Pointer(oldRootPath)),
		0,
	)
	if errno != 0 {
		fmt.Println("Error pivoting root filesystem:", errno)
		return
	}

	syscall.Syscall(
		syscall.SYS_CHDIR,
		uintptr(unsafe.Pointer(targetPath)),
		0,
		0,
	)
}

func MountProc() {
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
	fmt.Println("Mounting proc filesystem...")
	if errno != 0 {
		fmt.Println("Error mounting proc filesystem:", errno)
		return
	}
	fmt.Println("Proc filesystem mounted successfully.")
}

func Cleanup() {
	fmt.Println("Cleaning up...")
	oldRootPath, _ := syscall.BytePtrFromString("/oldroot")
	_, _, errno := syscall.Syscall(
		syscall.SYS_UMOUNT2,
		uintptr(unsafe.Pointer(oldRootPath)),
		uintptr(syscall.MNT_DETACH),
		0,
	)
	if errno != 0 {
		fmt.Println("Error cleaning up:", errno)
	}
	fmt.Println("Cleanup completed successfully.")
}
