package performance

import "encoding/json"

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

// Stats will be serialised and passed via the http stream to the client
type Stats struct {
	Uptime  uint64                      `json:"uptime"`
	Memory  UsageEntry                  `json:"memory"`
	Cpu     map[string]CpuCore          `json:"cpu"`
	Mount   map[string]UsageEntry       `json:"mount"`
	Network map[string]NetworkInterface `json:"network"`
}

// Json marshal the message struct to json ready for transit
func (m Stats) Json() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
