package terraform

// Ops interface for the operation
type Ops interface {
	Apply() error
	Init() error
}
