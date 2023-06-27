package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

// a subsystem is bound to a specified cgroup
type Subsystem interface {
	Name() string
	// set the resource config for a cgroup
	SetResourceConfig(cgPath string, res *ResourceConfig) error
	AddProc(cgPath string, pid int) error
	Remove(cgPath string) error
}

var (
	Subsystems = []Subsystem{
		&MemorySubsystem{},
	}
)
