package pool

type DefaultPool struct{}

func NewDefaultPool() *DefaultPool { return &DefaultPool{} }
func (DefaultPool) New(fn func()) {
	go fn()
}
