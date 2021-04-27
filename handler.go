package actioncable

type Handler interface {
	Handle(*Client, *Event)
}

type HandlerFunc func(*Client, *Event)

func (f HandlerFunc) Handle(c *Client, e *Event) {
	f(c, e)
}
