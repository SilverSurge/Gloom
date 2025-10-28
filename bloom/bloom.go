package bloom

type Bloom struct {
	State BloomDS
}

// `NewBloomDefault` return a default `Bloom` object
func NewBloomDefault(id string, n_bits, n_hash uint64) *Bloom {
	return NewBloomCustom(id, n_bits, n_hash, [2]uint64{DefaultSeed1, DefaultSeed2})
}

// `NewBloomCustom` return a custom `Bloom` object
func NewBloomCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) *Bloom {
	bloom := Bloom{
		State: NewBloomDSCustom(id, n_bits, n_hash, seeds),
	}
	return &bloom
}

// `NewBloomFromBloomDS`: return a `Bloom` using the data from bloom_ds
func NewBloomFromBloomDS(b *BloomDS) *Bloom {
	bloom := NewBloomCustom(b.ID, b.NBits, b.NHash, b.Seeds)
	bloom.Union(b)
	return bloom
}

// `Add`: add a value to the set
func (b *Bloom) Add(value any) {
	// find the indices
	indices := b.State.GetIndices(value)

	// find word index and offset, and set it to true
	for _, index := range indices {
		wi := index / 64
		off := index % 64
		b.State.Filter[wi] |= (1 << off)
	}
}

// `Check`: check a value to the set (false negative: never, false positives: maybe)
func (b *Bloom) Check(value any) bool {
	// find the indices
	indices := b.State.GetIndices(value)

	// // find word index and offset, and check if it is false
	for _, index := range indices {
		wi := index / 64
		off := index % 64

		if (b.State.Filter[wi] & (1 << off)) == 0 {
			return false
		}
	}

	return true
}

// `Reset`: resets bloom_ds
func (b *Bloom) Reset() {
	b.State.Reset()
}

// `Union`: tries state union
func (b1 *Bloom) Union(b2 *BloomDS) bool {
	return b1.State.Union(b2)
}

// `GetState`: return current State bool
func (b *Bloom) GetState() BloomDS {
	return b.State
}

// complie-time check
var _ IBloom = (*Bloom)(nil)
