package collect

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// The path of the proc filesystem.
var procPath = "/proc/loadavg"
var macLoadCmd = "sysctl -n vm.loadavg"

const shellToUse = "/bin/bash"

// GetLoad Read loadavg from /proc.
func GetLoad() (loads []float64, err error) {
	osName := runtime.GOOS
	switch osName {
	case "linux":
		data, err := ioutil.ReadFile(procPath)
		if err != nil {
			return nil, err
		}
		loads, err = parseLoad(string(data))
		if err != nil {
			return nil, err
		}
		return loads, nil
	case "darwin":
		cmd, err := exec.Command(shellToUse, "-c", macLoadCmd).Output()
		if err != nil {
			return nil, err
		}
		out := strings.Trim(string(cmd), "{}")
		loads, err = parseLoad(out)
		if err != nil {
			return nil, err
		}
		return loads, nil
	default:
		return nil, fmt.Errorf("not support os: %s", osName)
	}
}

// Parse /proc loadavg and return 1m, 5m and 15m.
func parseLoad(data string) (loads []float64, err error) {
	loads = make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected content in %s", procPath)
	}
	for i, load := range parts[0:3] {
		loads[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse load '%s': %w", load, err)
		}
	}
	return loads, nil
}
