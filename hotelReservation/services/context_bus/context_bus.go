package context_bus

import (
	"fmt"
	cb "github.com/AleckDarcy/ContextBus/proto"
	"os"
	"strconv"

	"github.com/AleckDarcy/ContextBus/background"
	"github.com/AleckDarcy/ContextBus/configure"
	cb_http "github.com/AleckDarcy/ContextBus/third-party/go/net/http"
)

var HOSTNAME = os.Getenv("HOSTNAME")
var GOLANG_VERSION = os.Getenv("GOLANG_VERSION")

var CONTEXTBUS_ON bool
var CONTEXTBUS_TRACE_SAMPLE_RATIO float64

func init() {
	tmpInt, err := strconv.Atoi(os.Getenv("CONTEXTBUS_ON"))
	if err != nil {
		fmt.Println("lookup CONTEXTBUS_ON from env fail:", err)
	} else {
		CONTEXTBUS_ON = tmpInt == 1
	}

	CONTEXTBUS_TRACE_SAMPLE_RATIO, err = strconv.ParseFloat(os.Getenv("CONTEXTBUS_TRACE_SAMPLE_RATIO"), 64)
	if err != nil {
		fmt.Println("lookup CONTEXTBUS_TRACE_SAMPLE_RATIO from env fail:", err)
		CONTEXTBUS_TRACE_SAMPLE_RATIO = 0.01
	}
}

func Set(cfg *configure.ServerConfigure, init ...func()) {
	// read env
	if !CONTEXTBUS_ON {
		fmt.Println("ContextBus is turned off")
		return
	} else if len(HOSTNAME) == 0 {
		fmt.Println("lookup HOSTNAME from env fail")
		return
	} else if len(GOLANG_VERSION) == 0 {
		fmt.Println("lookup GOLANG_VERSION from env fail")
		return
	}

	fmt.Printf("Initialize ContextBus(HOSTNAME=%s, GOLANG_VERSION=%s, CONTEXTBUS_ON=%v, PERF_METRIC=%v)\n", HOSTNAME, GOLANG_VERSION, CONTEXTBUS_ON, cb.PERF_METRIC)

	// run background tasks
	background.Run(cfg)

	// turn on http switches
	cb_http.TurnOn()

	for _, f := range init {
		f()
	}
}
