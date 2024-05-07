package context_bus

import (
	cb_configure "github.com/AleckDarcy/ContextBus/configure"
	cb "github.com/AleckDarcy/ContextBus/proto"
)

// todo parent name: event name of parent span
var defaultConfigure = &cb.Configure{
	Observations: map[string]*cb.ObservationConfigure{
		"frontend.searchHandler.1": {
			Type:    cb.ObservationType_ObservationStart,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"frontend.searchHandler.3": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Name:       "searchClient.Nearby",
				PrevName:   "",
				Attrs:      nil,
				Stacktrace: nil,
				ParentName: "",
			},
		},
		"frontend.searchHandler.4": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Name:       "searchClient.Nearby",
				PrevName:   "frontend.searchHandler.3",
				Attrs:      nil,
				Stacktrace: nil,
				ParentName: "",
			},
		},
		"frontend.searchHandler.5": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Name:       "reservationClient.CheckAvailability",
				PrevName:   "frontend.searchHandler.4",
				Attrs:      nil,
				Stacktrace: nil,
				ParentName: "",
			},
		},
		"Search.Nearby.Handler.1": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Name:       "Search.Nearby.Handler",
				PrevName:   "",
				Attrs:      nil,
				Stacktrace: nil,
				ParentName: "",
			},
		},
		"Search.Nearby.Handler.2": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Name:       "Search.Nearby.Handler",
				PrevName:   "Search.Nearby.Handler.1",
				Attrs:      nil,
				Stacktrace: nil,
				ParentName: "",
			},
		},
	},
}

func SetDefaultConfigure() {
	cb_configure.Store.SetDefault(defaultConfigure)
}
