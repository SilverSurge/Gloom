package bloom

import "sync/atomic"

type BloomAtomic struct {
	State BloomDS
}

// `NewBloomAtomicDefault` return a default `BloomAtomic` object
func NewBloomAtomicDefault(id string, n_bits, n_hash uint64) *BloomAtomic {
	return NewBloomAtomicCustom(id, n_bits, n_hash, [2]uint64{DefaultSeed1, DefaultSeed2})
}

// `NewBloomAtomicCustom` return a custom `BloomAtomic` object
func NewBloomAtomicCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) *BloomAtomic {
	bloom := BloomAtomic{
		State: NewBloomDSCustom(id, n_bits, n_hash, seeds),
	}
	return &bloom
}

// `NewBloomAtomicFromBloomDS`: return a `BloomAtomic` using the data from the bloom_ds
func NewBloomAtomicFromBloomDS(b *BloomDS) *BloomAtomic {
	bloom := NewBloomAtomicCustom(b.ID, b.NBits, b.NHash, b.Seeds)
	bloom.Union(b)
	return bloom
}

// `Add`: add a value to the set
func (b *BloomAtomic) Add(value any) {
	// find the indices
	indices := b.State.GetIndices(value)

	// find word index and offset, and set it to true
	for _, index := range indices {
		wi := index / 64
		off := index % 64
		mask := uint64(1) << off

		for {
			old := atomic.LoadUint64(&b.State.Filter[wi])
			if old&mask != 0 {
				break
			}
			if atomic.CompareAndSwapUint64(&b.State.Filter[wi], old, old|mask) {
				break
			}
		}
	}
}

// `Check: check a value to the set (false negative: never, false positives: maybe)
func (b *BloomAtomic) Check(value any) bool {
	// find the indices
	indices := b.State.GetIndices(value)

	// // find word index and offset, and check if it is false
	for _, index := range indices {
		wi := index / 64
		off := index % 64
		mask := uint64(1) << off

		v := atomic.LoadUint64(&b.State.Filter[wi])
		if v&mask == 0 {
			return false
		}
	}

	return true
}

// `Reset`: resets bloom_ds
func (b *BloomAtomic) Reset() {
	b.State.Reset()
}

// `Union`: tries state union
func (b1 *BloomAtomic) Union(b2 *BloomDS) bool {
	return b1.State.Union(b2)
}

// `GetState`: return current State bool
func (b *BloomAtomic) GetState() BloomDS {
	return b.State
}

// complie-time check
var _ IBloom = (*BloomAtomic)(nil)
