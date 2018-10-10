package process

// A Processor service processes messages from streaming message brokers.
type Processor interface {
	Process(b []byte) (oBytes []byte, err error)
}

type basicProcessorService struct{}

// Process messages from streaming services.
func (bps *basicProcessorService) Process(b []byte) (oBytes []byte, err error) {
	// TODO: implement based on use case
	oBytes = b
	return oBytes, err
}

// NewProcessorService processes and returns messages.
func NewProcessorService() Processor {
	return &basicProcessorService{}
}
