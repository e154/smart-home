package events

// EventLoadedPlugin ...
type EventLoadedPlugin struct {
	PluginName string `json:"plugin_name"`
}

// EventUnloadedPlugin ...
type EventUnloadedPlugin struct {
	PluginName string `json:"plugin_name"`
}
