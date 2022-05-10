package appch

type IValue struct {
	Index int
	Value interface{}
}

func Merge(size int, done <-chan struct{}, cs ...chan interface{}) chan IValue {
	out := make(chan IValue, size)
	for i, c := range cs {
		go func(i int, c <-chan interface{}) {
			for {
				select {
				case v, ok := <-c:
					if !ok {
						return
					}
					select {
					case out <- IValue{i, v}:
					case <-done:
						return
					}
				case <-done:
					return
				}
			}
		}(i, c)
	}
	return out
}
