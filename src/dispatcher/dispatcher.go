package dispatcher

// Dispatcher defines an interface for dispatching new URLs found in downloaded content
type Dispatcher interface {
	DispatchNewUrls([]string) (int, error)
}
