package appch

func ToInterfaceChan[T interface{}](c <-chan T) chan interface{} {
	out := make(chan interface{})
	go func() {
		for v := range c {
			out <- v
		}
		close(out)
	}()
	return out
}
