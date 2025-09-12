package nats

type Connector interface {
	Publish(subj string, data []byte) error
	Close()
}
