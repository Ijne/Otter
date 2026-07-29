package main

import (
	"fmt"
	"time"

	"github.com/Ijne/Otter/sandbox"
)

func init() {
	sandbox.Register("default", func() {
		fmt.Println("Hello from the default entry point!")
		fmt.Println(1 + 2)
	})
}

func main() {
	sandbox.Boot()

	res := sandbox.Run(sandbox.Config{
		EntryPoint:  "default",
		RootFS:      "/tmp/rootfs",
		MemoryLimit: 512 * 1024 * 1024,
		CPUQuota:    50000,
		TimeLimit:   10 * time.Second,
	})

	fmt.Println(res.ExitCode, res.TimedOut, string(res.Stdout))
}
