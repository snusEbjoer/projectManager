package screen

type Handler = func(...any)

type Controls struct {
	handlers  map[Event]Handler
	eventChan <-chan Event
}

func NewControls(eventChan <-chan Event) *Controls {
	c := &Controls{
		handlers:  make(map[Event]Handler),
		eventChan: eventChan,
	}
	c.run()
	return c
}

func (c *Controls) AddHandler(event Event, handler Handler) {
	c.handlers[event] = handler
}

func (c *Controls) run() {
	go func() {
		for e := range c.eventChan {
			handler, ok := c.handlers[e]
			if ok {
				handler()
			}
		}
	}()
}
