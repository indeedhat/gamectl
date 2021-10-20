package performance

const (
	PollingInterval = 2
)

var (
	networkHistory  map[string]NetworkInterface
	cpuHistory      map[string]CpuCore
	monitorInstance *Monitor
)

func init() {
	networkHistory = make(map[string]NetworkInterface)
	cpuHistory = make(map[string]CpuCore)
	monitorInstance = newMonitor()
}
