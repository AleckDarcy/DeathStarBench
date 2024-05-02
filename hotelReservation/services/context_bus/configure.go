package context_bus

import (
	cb_configure "github.com/AleckDarcy/ContextBus/configure"
	cb "github.com/AleckDarcy/ContextBus/proto"
)

var defaultConfigure = &cb.Configure{
	Observations: map[string]*cb.ObservationConfigure{},
}

func SetDefaultConfigure() {
	cb_configure.Store.SetDefault(defaultConfigure)
}
