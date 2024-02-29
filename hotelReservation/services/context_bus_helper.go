package services

import (
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"go.opentelemetry.io/otel/trace"
)

var _ = push.Pusher{}
var _ = cpu.InfoStat{}
var _ = host.InfoStat{}
var _ = mem.VirtualMemoryStat{}
var _ = net.InterfaceStat{}
var _ = trace.TraceID{}
