package context_bus

import (
	cb_configure "github.com/AleckDarcy/ContextBus/configure"
	cb "github.com/AleckDarcy/ContextBus/proto"
)

// todo parent name: event name of parent span
var defaultConfigure = &cb.Configure{
	Observations: map[string]*cb.ObservationConfigure{
		"hotels.1": {
			Type: cb.ObservationType_ObservationStart,
			Tracing: &cb.TracingConfigure{
				Start: true,
			},
		},
		"hotels.2": {
			Type:    cb.ObservationType_ObservationEnd,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "hotels",
				PrevEventName: "hotels.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
		"frontend.searchHandler.1": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"frontend.searchHandler.3": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start:      true,
				ParentName: "hotels.1",
			},
		},
		"frontend.searchHandler.4": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start:         true,
				End:           true,
				SpanName:      "searchClient.Nearby",
				PrevEventName: "frontend.searchHandler.3",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "hotels.1",
			},
		},
		"frontend.searchHandler.5": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "reservationClient.CheckAvailability",
				PrevEventName: "frontend.searchHandler.4",
				Attrs:         nil,
				Stacktrace:    nil,
			},
		},
		"Search.Nearby.Handler.1": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start: true,
			},
		},
		"Search.Nearby.Handler.2": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "Search.Nearby.Handler",
				PrevEventName: "Search.Nearby.Handler.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
	},
}

func SetConfigureForTesting() {
	cb_configure.Store.SetDefault(defaultConfigure)

	// bypass tracing
	cfg := &cb.Configure{
		Observations: map[string]*cb.ObservationConfigure{},
	}

	for key, val := range defaultConfigure.Observations {
		newVal := &cb.ObservationConfigure{
			Type:    val.Type,
			Logging: val.Logging,
			Metrics: val.Metrics,
		}

		cfg.Observations[key] = newVal
	}

	cb_configure.Store.SetConfigure(cb_configure.CBCID_TRACEBYPASS, cfg)
}
