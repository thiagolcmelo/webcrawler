package memory

// Frontier is an in memory implementation of the Frontier interface
type Frontier struct {
	jobs chan string
}

// NewFrontier is a factory for an in memory Frontier
func NewFrontier() *Frontier {
	return &Frontier{
		jobs: make(chan string),
	}
}

// Publish adds a job/message/url to the queue
func (mf *Frontier) Publish(address string) error {
	mf.jobs <- address
	return nil
}

// Consume returns a channel to read from the queue
func (mf *Frontier) Consume() <-chan string {
	return mf.jobs
}
