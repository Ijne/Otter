package sandbox

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Ijne/Otter/internal/ns"
	"github.com/Ijne/Otter/internal/params"
)

var (
	registry = map[string]func(){}
)

type Config struct {
	EntryPoint  string
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

func copySelfBinary(rootfs string) error {
	selfPath, err := os.Executable()
	if err != nil {
		return err
	}

	src, err := os.Open(selfPath)
	if err != nil {
		return err
	}
	defer src.Close()

	destPath := filepath.Join(rootfs, "bin", "self")
	dst, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}

func copyFileIntoRootFS(rootfs, sourcePath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	destinationPath := filepath.Join(rootfs, sourcePath)
	if err := os.MkdirAll(filepath.Dir(destinationPath), 0755); err != nil {
		return err
	}

	destination, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return err
	}

	return nil
}

func copySelfDependencies(rootfs string) error {
	selfPath, err := os.Executable()
	if err != nil {
		return err
	}

	command := exec.Command("ldd", selfPath)
	output, err := command.CombinedOutput()
	if err != nil && len(output) == 0 {
		return err
	}

	dependencies := map[string]struct{}{}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.Contains(line, "statically linked") {
			continue
		}

		if strings.Contains(line, "=>") {
			parts := strings.SplitN(line, "=>", 2)
			if len(parts) != 2 {
				continue
			}
			candidate := strings.Fields(strings.TrimSpace(parts[1]))
			if len(candidate) == 0 {
				continue
			}
			if strings.HasPrefix(candidate[0], "/") {
				dependencies[candidate[0]] = struct{}{}
			}
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if strings.HasPrefix(fields[0], "/") {
			dependencies[fields[0]] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for dependency := range dependencies {
		if err := copyFileIntoRootFS(rootfs, dependency); err != nil {
			return err
		}
	}

	return nil
}

func prepareRootFS(rootfs string) error {
	fmt.Println("Preparing root filesystem at:", rootfs)

	dirs := []string{"bin", "proc", "oldroot", "dev", "tmp"}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(rootfs, d), 0755); err != nil {
			return err
		}
	}
	fmt.Println("Directories created successfully.")

	busyboxPath := filepath.Join(rootfs, "bin", "busybox")
	if err := downloadFile("https://busybox.net/downloads/binaries/1.35.0-x86_64-linux-musl/busybox", busyboxPath); err != nil {
		return err
	}
	if err := os.Chmod(busyboxPath, 0755); err != nil {
		return err
	}

	links := []string{"sh", "ls", "ps", "mount", "cat", "id", "echo", "clear"}
	for _, name := range links {
		target := filepath.Join(rootfs, "bin", name)
		os.Remove(target)
		if err := os.Symlink("busybox", target); err != nil {
			return err
		}
	}

	fmt.Println("Copying self binary to root filesystem...")
	if err := copySelfBinary(rootfs); err != nil {
		fmt.Println("Error copying self binary:", err)
		return err
	}
	if err := copySelfDependencies(rootfs); err != nil {
		fmt.Println("Error copying self dependencies:", err)
		return err
	}
	fmt.Println("Self binary copied successfully.")

	return nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func Run(cfg Config) *Result {
	prepareRootFS(cfg.RootFS)

	p := params.Params{
		EntryPoint:  cfg.EntryPoint,
		Flags:       uintptr(syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID | syscall.CLONE_NEWIPC | syscall.CLONE_NEWCGROUP | syscall.SIGCHLD),
		RootFS:      cfg.RootFS,
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

func Register(name string, fn func()) {
	fmt.Println("Registering function:", name)
	registry[name] = fn
	fmt.Println("Function registered successfully.")
}

func Boot() {
	name := os.Getenv("OTTER_REEXEC")
	if name == "" {
		return
	}
	fmt.Println("Re-executing started")
	fn, ok := registry[name]
	if !ok {
		os.Exit(127)
	}
	fmt.Println("Re-executing finished")
	fn()
	os.Exit(0)
}
