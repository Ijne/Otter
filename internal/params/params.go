package params

import "time"

type Params struct {
	EntryPoint  string
	Flags       uintptr
	RootFS      string
	MemoryLimit int64
	CPUQuota    int64
	TimeLimit   time.Duration
}
