package work

type Supplier interface {
	Supply() ([]byte, error)
}
