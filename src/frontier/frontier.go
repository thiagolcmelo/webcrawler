package frontier

type Frontier interface {
	Push(string) error
	Pop() <-chan string
}
