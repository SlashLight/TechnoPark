package main

func ExecutePipeline(jobs ...job) {
	out := make(chan interface{})
	for _, worker := range jobs {
		in := out
		out := make(chan interface{})
		go func() {
			worker(in, out)
			close(out)
		}()
	}
}
func SingleHash(in, out chan interface{})     {}
func MultiHash(in, out chan interface{})      {}
func CombineResults(in, out chan interface{}) {}
