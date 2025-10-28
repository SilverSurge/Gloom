package bloom

import "sync"

type BloomRW struct {
	State    BloomDS
	Mu       sync.RWMutex
	rareMu sync.RWMutex
}

// `NewBloomRWDefault` return a default `BloomRW` object
func NewBloomRWDefault(id string, n_bits, n_hash uint64) *BloomRW {
	return NewBloomRWCustom(id, n_bits, n_hash, [2]uint64{DefaultSeed1, DefaultSeed2})
}

// `NewBloomRWCustom` return a custom `BloomRW` object
func NewBloomRWCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) *BloomRW {
	bloom := BloomRW{
		State: NewBloomDSCustom(id, n_bits, n_hash, seeds),
	}
	return &bloom
}

// `NewBloomRWFromBloomDS`: return a `BloomRW` using the data from the bloom_ds
func NewBloomRWFromBloomDS(b *BloomDS) *BloomRW {
	bloom := NewBloomRWCustom(b.ID, b.NBits, b.NHash, b.Seeds)
	bloom.Union(b)
	return bloom
}

// `Add`: add a value to the set
func (b *BloomRW) Add(value any) {
	// unionRW mutex
	b.rareMu.RLock()
	defer b.rareMu.RUnlock()

	// find the indices
	indices := b.State.GetIndices(value)

	// get write lock
	b.Mu.Lock()
	defer b.Mu.Unlock()

	// find word index and offset, and set it to true
	for _, index := range indices {
		wi := index / 64
		off := index % 64
		b.State.Filter[wi] |= (1 << off)
	}
}

// `Check`: check a value to the set (false negative: never, false positives: maybe)
func (b *BloomRW) Check(value any) bool {
	// unionRW mutex
	b.rareMu.RLock()
	defer b.rareMu.RUnlock()

	// find the indices
	indices := b.State.GetIndices(value)

	// get read lock
	b.Mu.RLock()
	defer b.Mu.RUnlock()

	// find word index and offset, and check if it is false
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
func (b *BloomRW) Reset() {
	// unionRW mutex
	b.rareMu.Lock()
	defer b.rareMu.Unlock()

	b.State.Reset()
}

// `Union`: tries state union
func (b1 *BloomRW) Union(b2 *BloomDS) bool {
	// unionRW mutex
	b1.rareMu.Lock()
	defer b1.rareMu.Unlock()

	return b1.State.Union(b2)
}

// `GetState`: return current State bool
func (b *BloomRW) GetState() BloomDS {
	// unionRW mutex
	b.rareMu.Lock()
	defer b.rareMu.Unlock()

	return b.State
}

// complie-time check
var _ IBloom = (*BloomRW)(nil)
