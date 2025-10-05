package bloom

type IBloom interface {
	Add(any)
	Check(any) bool
	Reset()
}
