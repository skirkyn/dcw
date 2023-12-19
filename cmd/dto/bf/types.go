package bf

type Request struct {
	BatchSize int `json:"batchSize"`
}

type Response[T any] struct {
	Data T      `json:"data"`
	Done bool   `json:"done"`
	Err  string `json:"err"`
}

func (r Request) Body() int {
	return r.BatchSize
}

func (r Response[T]) Error() string {
	return r.Err
}
