package params

type Params struct {
	Flags   uintptr
	RootFS  string
	Command []string
	Env     []string
}
