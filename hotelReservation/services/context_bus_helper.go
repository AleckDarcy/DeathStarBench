package services

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

var _ = mem.VirtualMemoryStat{}
var _ = host.InfoStat{}
var _ = cpu.InfoStat{}
var _ = net.InterfaceStat{}