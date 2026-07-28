package main

import (
	"fmt"
	"otter/pkg"
	"time"
)

func main() {
	res := pkg.Run(pkg.Config{
		RootFS:      "/tmp/rootfs",
		Command:     []string{"/bin/sh"},
		Env:         []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		MemoryLimit: 512 * 1024 * 1024,
		CPUQuota:    50000,
		TimeLimit:   10 * time.Second,
	})

	fmt.Println(res.ExitCode, res.TimedOut, string(res.Stdout))
}
