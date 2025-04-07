package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/oliashish/vofo/logger"
	"github.com/tklauser/go-sysconf"
)

var log = logger.Logger()

// TODO: Match the type below to actual types, not all string
type Proc struct {
	pid       string
	process   string
	uptime    int
	stime     int
	starttime int
}

type CpuUsage struct {
	pid     string
	process string
	usage   float64
}

var jiffies = sysconf.SC_CLK_TCK

func getSystemUptime() float64 {
	data, _ := os.ReadFile("/proc/uptime")
	uptimeStr := strings.Fields(string(data))[0]
	uptime, _ := strconv.ParseFloat(uptimeStr, 64)
	return uptime
}

func parseProcDir(procDir []os.DirEntry) []Proc {
	processes := []Proc{}

	for _, p := range procDir {

		info, _ := p.Info()

		path := filepath.Join("/proc", info.Name(), "stat")

		procStat, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Unable to read stat file: ", err)
			continue
		}

		// TODO: Handle process name as it is enclosed in paranthesis ()
		// and it can have space in this which can break this " " delimiter
		stat := strings.Split(string(procStat), " ")
		proc := Proc{
			pid:     stat[0],
			process: stat[1],
			// TODO: Better handle this conversion in it's own function which handles error and returns only int
			uptime:    func() int { i, _ := strconv.Atoi(stat[13]); return i }(),
			stime:     func() int { i, _ := strconv.Atoi(stat[14]); return i }(),
			starttime: func() int { i, _ := strconv.Atoi(stat[21]); return i }(),
		}

		processes = append(processes, proc)

	}
	return processes
}

func getCpuUsage(processes []Proc) []CpuUsage {
	sysUptime := getSystemUptime()

	usage := []CpuUsage{}

	for _, process := range processes {
		totalTime := float64(process.uptime + process.stime)
		seconds := sysUptime - (float64(process.starttime) / float64(jiffies))

		if seconds <= 0 {
			continue
		}

		cpuUsage := 100 * ((totalTime / float64(jiffies)) / seconds)
		usage = append(usage, CpuUsage{
			pid:     process.pid,
			process: process.process,
			usage:   cpuUsage,
		})
	}
	return usage
}

func CPU() {
	procDir, err := os.ReadDir("/proc")
	if err != nil {
		log.Error(fmt.Sprintf("Error while reading /proc Directory: %s\n", err))
	}

	processes := parseProcDir(procDir)
	cpuPercentage := getCpuUsage(processes)

	fmt.Println(cpuPercentage)

}
