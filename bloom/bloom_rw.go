package bloom

import "sync"

type BloomRW struct {
	id     string
	n_bits uint64
	n_hash uint64
	seeds  [2]uint64
	filter []uint64

	mu sync.RWMutex
}

// `NewBloomRWDefault` return a default `BloomRW` object
func NewBloomRWDefault(id string, n_bits, n_hash uint64) *BloomRW {
	return NewBloomRWCustom(id, n_bits, n_hash, [2]uint64{DefaultSeed1, DefaultSeed2})
}

// `NewBloomRWCustom` return a custom `BloomRW` object
func NewBloomRWCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) *BloomRW {

	n_words := (n_bits + 63) / 64
	bloom := BloomRW{
		id:     id,
		n_bits: n_bits,
		n_hash: n_hash,
		seeds:  seeds,
		filter: make([]uint64, n_words),
	}
	return &bloom
}

// `Add`: add a value to the set
func (b *BloomRW) Add(value any) {
	// find the indices
	indices := b.getIndices(value)

	// get write lock
	b.mu.Lock()
	defer b.mu.Unlock()

	// find word index and offset, and set it to true
	for _, index := range indices {
		wi := index / 64
		off := index % 64
		b.filter[wi] |= (1 << off)
	}
}

// `Check`: check a value to the set (false negative: never, false positives: maybe)
func (b *BloomRW) Check(value any) bool {
	// find the indices
	indices := b.getIndices(value)

	// get read lock
	b.mu.RLock()
	defer b.mu.RUnlock()

	// // find word index and offset, and check if it is false
	for _, index := range indices {
		wi := index / 64
		off := index % 64

		if (b.filter[wi] & (1 << off)) == 0 {
			return false
		}
	}

	return true
}

// `getIndices`: find filter indices
func (b *BloomRW) getIndices(value any) []uint64 {
	// get bytes
	data := toBytes(value)

	// get primary hashes
	h1 := hash(b.seeds[0], data) % b.n_bits
	h2 := hash(b.seeds[1], data) % b.n_bits

	// use double hashing to generate n_hash indices
	indices := make([]uint64, b.n_hash)
	m := b.n_bits

	for i := uint64(0); i < uint64(b.n_hash); i++ {
		index := (h1 + ((i*h2)%m)%m + m) % m
		indices[i] = index
	}

	return indices
}

// complie-time check
var _ IBloom = (*BloomRW)(nil)
