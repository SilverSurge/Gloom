package bloom

type IBloom interface {
	Add(any)
	Check(any) bool
	Reset()
	Union(*BloomDS) bool
	GetState() BloomDS
}
