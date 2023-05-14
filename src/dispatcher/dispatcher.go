package dispatcher

type Dispatcher interface {
	DispatchNewUrls([]string) (int, error)
}
