package performance

import (
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/indeedhat/gamectl/internal/config"
)

type Monitor struct {
	stats Stats
	conf  *config.Performance
}

// GetMonitor will recieve the currently running monitor instance
func GetMonitor() *Monitor {
	return monitorInstance
}

// NewMonitor create a new Monitor instance and start it running
func newMonitor() *Monitor {
	conf, _ := config.LoadPerformanceConfig()
	monitor := &Monitor{
		conf: conf,
	}

	monitor.run()

	return monitor
}

// Read will return the current stats
func (m *Monitor) Read() Stats {
	return m.stats
}

// run start the monitor watching
func (m *Monitor) run() {
	go func() {
		for range time.NewTicker(time.Second * PollingInterval).C {
			m.stats = Stats{
				Uptime:  uptime(),
				Cpu:     cpu(),
				Memory:  memory(),
				Mount:   mount(m),
				Network: network(m),
			}
		}
	}()

	return
}

// uptime calculate the system uptime
func uptime() (uptime uint64) {
	line, err := ioutil.ReadFile("/proc/uptime")
	if nil != err {
		return
	}

	parts := strings.Split(string(line), ".")
	uptime, _ = strconv.ParseUint(parts[0], 10, 64)

	return
}

// memory check the current memory usage
func memory() (memory UsageEntry) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if nil != err {
		return
	}

	var free uint64
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal") {
			fields := strings.Fields(line)
			memory.Total, _ = strconv.ParseUint(fields[1], 10, 64)
			memory.Total *= 1024
		} else if strings.HasPrefix(line, "MemAvailable") {
			fields := strings.Fields(line)
			free, _ = strconv.ParseUint(fields[1], 10, 64)
			free *= 1024
		}
	}

	memory.Used = memory.Total - free

	return
}

// mount scans the drive mounts based on the configured whitelist
// and returns their space stats
func mount(m *Monitor) (mounts map[string]UsageEntry) {
	mounts = make(map[string]UsageEntry)

	if m.conf == nil || len(m.conf.Mount) == 0 {
		return
	}

	for _, mount := range m.conf.Mount {
		fs := syscall.Statfs_t{}
		err := syscall.Statfs(mount, &fs)
		if err != nil {
			mounts[mount] = UsageEntry{0, 0}
			continue
		}

		total := fs.Blocks * uint64(fs.Bsize)
		free := fs.Bfree * uint64(fs.Bsize)

		mounts[mount] = UsageEntry{
			Total: total,
			Used:  total - free,
		}
	}

	return
}

// cpu calculate the current usage of the cpu cores
func cpu() (cores map[string]CpuCore) {
	cores = make(map[string]CpuCore)

	contents, err := ioutil.ReadFile("/proc/stat")
	if nil != err {
		return
	}

	for _, line := range strings.Split(string(contents), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || "cpu" == fields[0] || !strings.HasPrefix(fields[0], "cpu") {
			continue
		}

		key := fields[0]
		core := CpuCore{}

		for i, field := range fields[1:] {
			val, _ := strconv.ParseUint(field, 10, 64)

			core.Total += val

			if 3 == i {
				core.Idle = val
			}
		}

		if _, ok := cpuHistory[key]; !ok {
			cpuHistory[key] = CpuCore{core.Total, core.Idle}
		}

		cores[key] = CpuCore{
			Total: (core.Total - cpuHistory[key].Total) / PollingInterval,
			Idle:  (core.Idle - cpuHistory[key].Idle) / PollingInterval,
		}

		cpuHistory[key] = core
	}

	return
}

// network calculate the usage on each network interface since the last poll
func network(m *Monitor) (interfaces map[string]NetworkInterface) {
	interfaces = make(map[string]NetworkInterface)

	device, err := ioutil.ReadFile("/proc/net/dev")
	if nil != err {
		return
	}

	for _, line := range strings.Split(string(device), "\n") {
		if !strings.Contains(line, ":") {
			continue
		}

		name := strings.Trim(strings.Split(line, ":")[0], " ")
		if skipNetworkInterface(m, name) {
			continue
		}

		fields := strings.Fields(line)
		rx, _ := strconv.ParseUint(fields[1], 10, 64)
		tx, _ := strconv.ParseUint(fields[9], 10, 64)

		if _, ok := networkHistory[name]; !ok {
			networkHistory[name] = NetworkInterface{rx, tx}
		}

		interfaces[name] = NetworkInterface{
			Rx: (rx - networkHistory[name].Rx) / PollingInterval,
			Tx: (tx - networkHistory[name].Tx) / PollingInterval,
		}

		networkHistory[name] = NetworkInterface{rx, tx}
	}

	return
}

// skipNetworkInterface check the network interface name against the whitelist
// and see if it should be skipped or not
func skipNetworkInterface(m *Monitor, name string) bool {
	if m.conf == nil {
		return false
	}

	for _, tracked := range m.conf.Network {
		if name == tracked {
			return false
		}
	}

	return true
}
