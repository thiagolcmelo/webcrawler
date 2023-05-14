package memory

type MemoryFrontier struct {
	jobs chan string
}

func NewMemoryFrontier() *MemoryFrontier {
	return &MemoryFrontier{
		jobs: make(chan string),
	}
}

func (mf *MemoryFrontier) Push(address string) error {
	mf.jobs <- address
	return nil
}

func (mf *MemoryFrontier) Pop() <-chan string {
	return mf.jobs
}
