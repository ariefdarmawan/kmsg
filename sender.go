package kmsg

type Sender interface {
	Send(msg *Message) error
	Close()
}
