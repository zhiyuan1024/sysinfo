package cpuinfo

import (
	"log"
	"time"
)

type CPUInfo struct {
	interval int
	data     *cpuData
}

func NewCPUInfo() *CPUInfo {
	var cpu = CPUInfo{
		interval: 1,
	}
	cpu.data = NewCPUData()
	return &cpu
}

func (c *CPUInfo) Collecting(lineChan chan string) {
	for {
		if err := c.data.Collecting(); err != nil {
			log.Printf("collecting cpuinfo error, err = %v", err)
			time.Sleep(time.Duration(c.interval) * time.Second)
			continue
		}
		if line, err := c.data.Dump(); err != nil {
			log.Printf("collecting cpuinfo error, err = %v", err)
			time.Sleep(time.Duration(c.interval) * time.Second)
			continue
		} else {
			lineChan <- line
		}
		time.Sleep(time.Duration(c.interval) * time.Second)
	}
}
