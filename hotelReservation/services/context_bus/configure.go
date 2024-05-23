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
		"frontend.searchHandler.1": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"frontend.searchHandler.2": {
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
				Start:         true,
				End:           true,
				SpanName:      "reservationClient.CheckAvailability",
				PrevEventName: "frontend.searchHandler.4",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "hotels.1",
			},
		},
		"frontend.searchHandler.6": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "profileClient.GetProfiles",
				PrevEventName: "frontend.searchHandler.5",
				Attrs:         nil,
				Stacktrace:    nil,
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
		"_Search_Nearby_Handler.1": {
			Type:    cb.ObservationType_ObservationStart,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start: true,
			},
		},
		"search.Nearby.1": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"search.Nearby.2": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"search.Nearby.3": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start:      true,
				ParentName: "_Search_Nearby_Handler.1",
			},
		},
		"search.Nearby.4": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start:         true,
				End:           true,
				SpanName:      "geoClient.Nearby",
				PrevEventName: "search.Nearby.3",
				ParentName:    "_Search_Nearby_Handler.1",
			},
		},
		"search.Nearby.5": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "rateClient.GetRates",
				PrevEventName: "search.Nearby.4",
				ParentName:    "_Search_Nearby_Handler.1",
			},
		},
		"search.Nearby.6": {
			Type: cb.ObservationType_ObservationInter,
			// Logging: cb_configure.DefaultJSONLogging, // for loop, too many log prints
		},
		"_Search_Nearby_Handler.2": {
			Type:    cb.ObservationType_ObservationEnd,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "Search.Nearby.Handler",
				PrevEventName: "_Search_Nearby_Handler.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
		"_Reservation_CheckAvailability_Handler.1": {
			Type:    cb.ObservationType_ObservationStart,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start: true,
			},
		},
		"reservation.CheckAvailability.1": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				Start:      true,
				ParentName: "_Reservation_CheckAvailability_Handler.1",
			},
		},
		"reservation.CheckAvailability.2": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "memcached_capacity_get_multi_number",
				PrevEventName: "reservation.CheckAvailability.1",
			},
		},
		"reservation.CheckAvailability.3": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				Start:      true,
				ParentName: "_Reservation_CheckAvailability_Handler.1",
			},
		},
		"reservation.CheckAvailability.4": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "mongodb_capacity_get_multi_number",
				PrevEventName: "reservation.CheckAvailability.3",
			},
		},
		"reservation.CheckAvailability.5": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				Start:      true,
				ParentName: "_Reservation_CheckAvailability_Handler.1",
			},
		},
		"reservation.CheckAvailability.6": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "memcached_reserve_get_multi_number",
				PrevEventName: "reservation.CheckAvailability.5",
			},
		},
		"_Reservation_CheckAvailability_Handler.2": {
			Type:    cb.ObservationType_ObservationEnd,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "Reservation.CheckAvailability.Handler",
				PrevEventName: "_Reservation_CheckAvailability_Handler.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
		"_Profile_GetProfiles_Handler.1": {
			Type:    cb.ObservationType_ObservationStart,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start: true,
			},
		},
		"profile.GetProfiles.1": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"profile.GetProfiles.2": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				Start:      true,
				ParentName: "_Profile_GetProfiles_Handler.1",
			},
		},
		"profile.GetProfiles.3": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "memcached_get_profile",
				PrevEventName: "profile.GetProfiles.2",
			},
		},
		"profile.GetProfiles.4": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"_Profile_GetProfiles_Handler.2": {
			Type:    cb.ObservationType_ObservationEnd,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "Profile.GetProfiles.Handler",
				PrevEventName: "_Profile_GetProfiles_Handler.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
		"_Geo_Nearby_Handler.1": {
			Type:    cb.ObservationType_ObservationStart,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start: true,
			},
		},
		"geo.Nearby.1": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"geo.Nearby.2": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"geo.Nearby.3": {
			Type: cb.ObservationType_ObservationInter,
			// Logging: cb_configure.DefaultJSONLogging, // for loop, too many log prints
		},
		"_Geo_Nearby_Handler.2": {
			Type:    cb.ObservationType_ObservationEnd,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "Geo.Nearby.Handler",
				PrevEventName: "_Geo_Nearby_Handler.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
		"_Rate_GetRates_Handler.1": {
			Type:    cb.ObservationType_ObservationStart,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				Start: true,
			},
		},
		"rate.GetRates.1": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				Start:      true,
				ParentName: "_Rate_GetRates_Handler.1",
			},
		},
		"rate.GetRates.2": {
			Type: cb.ObservationType_ObservationInter,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "memcached_get_multi_rate",
				PrevEventName: "rate.GetRates.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
		"rate.GetRates.3": {
			Type: cb.ObservationType_ObservationInter,
			// Logging: cb_configure.DefaultJSONLogging, // too long
		},
		"rate.GetRates.4": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"rate.GetRates.5": {
			Type:    cb.ObservationType_ObservationInter,
			Logging: cb_configure.DefaultJSONLogging,
		},
		"_Rate_GetRates_Handler.2": {
			Type:    cb.ObservationType_ObservationEnd,
			Logging: cb_configure.DefaultJSONLogging,
			Tracing: &cb.TracingConfigure{
				End:           true,
				SpanName:      "Rate.GetRates.Handler",
				PrevEventName: "_Rate_GetRates_Handler.1",
				Attrs:         nil,
				Stacktrace:    nil,
				ParentName:    "",
			},
		},
	},
	Reactions: map[string]*cb.ReactionConfigure{
		"_Search_Nearby_Handler.2": {
			Type:   cb.ReactionType_ReactionPrintLog,
			Params: nil,
			PreTree: &cb.PrerequisiteTree{
				Nodes: []*cb.PrerequisiteNode{
					{
						Id:   0,
						Type: cb.PrerequisiteNodeType_PrerequisiteAfterObservation_,
						PrevEvent: &cb.PrerequisiteEvent{
							Name:    "_Search_Nearby_Handler.1",
							Latency: 50,
						},
					},
				},
				LeafIDs: []int64{0},
			},
		},
	},
}

func SetConfigureForTesting() {
	cb_configure.Store.SetDefault(defaultConfigure)

	// bypass observation
	{
		cfg := &cb.Configure{
			Observations: map[string]*cb.ObservationConfigure{},
			Reactions:    defaultConfigure.Reactions,
		}

		for key, val := range defaultConfigure.Observations {
			newVal := &cb.ObservationConfigure{
				Type: val.Type,
			}

			cfg.Observations[key] = newVal
		}

		cb_configure.Store.SetConfigure(cb_configure.CBCID_OBSERVATIONBYPASS, cfg)
	}

	// bypass logging
	{
		cfg := &cb.Configure{
			Observations: map[string]*cb.ObservationConfigure{},
			Reactions:    defaultConfigure.Reactions,
		}

		for key, val := range defaultConfigure.Observations {
			newVal := &cb.ObservationConfigure{
				Type:    val.Type,
				Tracing: val.Tracing,
				Metrics: val.Metrics,
			}

			cfg.Observations[key] = newVal
		}

		cb_configure.Store.SetConfigure(cb_configure.CBCID_LOGGINGYPASS, cfg)
	}

	// bypass tracing
	{
		cfg := &cb.Configure{
			Observations: map[string]*cb.ObservationConfigure{},
			Reactions:    defaultConfigure.Reactions,
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
}
