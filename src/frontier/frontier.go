package frontier

// Frontier defines an interface for a queue for download jobs
type Frontier interface {
	Publish(string) error
	Consume() <-chan string
}
