package performance

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/indeedhat/command-center/app/config"
)

const (
	PollingInterval = 2
)

var (
	networkIntfCache map[string]NetworkInterface
	cpuTickCache     map[string]CpuCore
	monitorCache     *Monitor
)

func init() {
	networkIntfCache = make(map[string]NetworkInterface)
	cpuTickCache = make(map[string]CpuCore)
	monitorCache, _ = newMonitor()
}

type CpuCore struct {
	Total uint64 `json:"total"`
	Idle  uint64 `json:"idle"`
}

type UsageEntry struct {
	Used  uint64 `json:"used"`
	Total uint64 `json:"total"`
}

type NetworkInterface struct {
	Rx uint64 `json:"rx"`
	Tx uint64 `json:"tx"`
}

// Message will be serialised and passed via the http stream to the client
type Message struct {
	Uptime  uint64                      `json:"uptime"`
	Memory  UsageEntry                  `json:"memory"`
	Cpu     map[string]CpuCore          `json:"cpu"`
	Mount   map[string]UsageEntry       `json:"mount"`
	Network map[string]NetworkInterface `json:"network"`
}

// Json marshal the message struct to json ready for transit
func (m Message) Json() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

type Monitor struct {
	message Message
	conf    *config.Performance
}

// GetMonitor will recieve the currently running monitor instance
func GetMonitor() *Monitor {
	return monitorCache
}

// NewMonitor create a new Monitor instance and start it running
func newMonitor() (*Monitor, error) {
	conf, _ := config.LoadPerformanceConfig()
	monitor := &Monitor{
		conf: conf,
	}

	if err := monitor.run(); err != nil {
		return nil, err
	}

	return monitor, nil
}

// Read will return the current message
func (m *Monitor) Read() Message {
	return m.message
}

// run start the monitor watching
func (m *Monitor) run() error {
	go m.loop()

	return nil
}

// loop runs the logic for actually watching the server and building the message
func (m *Monitor) loop() {
	for range time.NewTicker(time.Second * PollingInterval).C {
		message := Message{
			Uptime:  m.uptime(),
			Cpu:     m.cpu(),
			Memory:  m.memory(),
			Mount:   m.mount(),
			Network: m.network(),
		}

		m.message = message
	}
}

// uptime calculate the system uptime
func (m *Monitor) uptime() (uptime uint64) {
	line, err := ioutil.ReadFile("/proc/uptime")
	if nil != err {
		return
	}

	parts := strings.Split(string(line), ".")
	uptime, _ = strconv.ParseUint(parts[0], 10, 64)

	return
}

// memory check the current memory usage
func (m *Monitor) memory() (memory UsageEntry) {
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
func (m *Monitor) mount() (mounts map[string]UsageEntry) {
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
		free := fs.Bavail * uint64(fs.Bsize)

		mounts[mount] = UsageEntry{
			Total: total,
			Used:  total - free,
		}
	}

	return
}

// cpu calculate the current usage of the cpu cores
func (m *Monitor) cpu() (cores map[string]CpuCore) {
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

		core := CpuCore{}
		for i, field := range fields[1:] {
			val, _ := strconv.ParseUint(field, 10, 64)

			core.Total += val

			if 3 == i {
				core.Idle = val
			}
		}

		if _, ok := cpuTickCache[fields[0]]; !ok {
			cpuTickCache[fields[0]] = CpuCore{core.Total, core.Idle}
		}

		prev := cpuTickCache[fields[0]]

		cores[fields[0]] = CpuCore{
			Total: (core.Total - prev.Total) / PollingInterval,
			Idle:  (core.Idle - prev.Idle) / PollingInterval,
		}

		cpuTickCache[fields[0]] = core
	}

	return
}

// network calculate the usage on each network interface since the last poll
func (m *Monitor) network() (interfaces map[string]NetworkInterface) {
	interfaces = make(map[string]NetworkInterface)

	device, err := ioutil.ReadFile("/proc/net/dev")
	if nil != err {
		return
	}

	for _, line := range strings.Split(string(device), "\n") {
		if !strings.Contains(line, ":") {
			continue
		}

		fields := strings.Fields(line)
		name := strings.Trim(strings.Split(line, ":")[0], " ")

		if m.skipNetworkInterface(name) {
			continue
		}

		rx, _ := strconv.ParseUint(fields[1], 10, 64)
		tx, _ := strconv.ParseUint(fields[9], 10, 64)

		if _, ok := networkIntfCache[name]; !ok {
			networkIntfCache[name] = NetworkInterface{rx, tx}
		}

		interfaces[name] = NetworkInterface{
			Rx: (rx - networkIntfCache[name].Rx) / PollingInterval,
			Tx: (tx - networkIntfCache[name].Tx) / PollingInterval,
		}

		networkIntfCache[name] = NetworkInterface{rx, tx}
	}

	return
}

// skipNetworkInterface check the network interface name against the whitelist
// and see if it should be skipped or not
func (m *Monitor) skipNetworkInterface(name string) bool {
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
