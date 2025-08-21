package commons

type Form interface {
	Validate() error
	GetData() any
}
