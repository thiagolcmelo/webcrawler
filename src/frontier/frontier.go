package frontier

type Frontier interface {
	Publish(string) error
	Consume() <-chan string
}
