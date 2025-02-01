package core

type Driver interface {
	Out(interface {}) error
	Err(int, string) error
}
