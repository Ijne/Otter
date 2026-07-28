package params

import "time"

type Params struct {
	Flags       uintptr
	RootFS      string
	Command     []string
	Env         []string
	MemoryLimit int64
	CPUQuota    int64
	TimeLimit   time.Duration
}
