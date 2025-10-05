package bloom

import (
	"sync"
)

type BloomShard struct {
	id       string
	n_bits   uint64
	n_hash   uint64
	seeds    [2]uint64
	filter   []uint64
	n_shards uint64
	shards   []sync.RWMutex

	len_long       uint64
	len_short      uint64
	n_long         uint64
	n_short        uint64
	boundary_index uint64
}

// `NewBloomShardDefault` return a default `BloomShard` object
func NewBloomShardDefault(id string, n_bits, n_hash, n_shards uint64) *BloomShard {
	return NewBloomShardCustom(id, n_bits, n_hash, n_shards, [2]uint64{DefaultSeed1, DefaultSeed2})
}

// `NewBloomShardCustom` return a custom `BloomShard` object
func NewBloomShardCustom(id string, n_bits, n_hash, n_shards uint64, seeds [2]uint64) *BloomShard {

	n_words := (n_bits + 63) / 64
	bloom := BloomShard{
		id:       id,
		n_bits:   n_bits,
		n_hash:   n_hash,
		seeds:    seeds,
		filter:   make([]uint64, n_words),
		n_shards: min(n_shards, n_bits),
		shards:   make([]sync.RWMutex, n_shards),

		len_long:  (n_bits + n_shards - 1) / n_shards,
		len_short: n_bits / n_shards,
		n_long:    n_bits % n_shards,
	}
	bloom.n_short = bloom.n_bits - bloom.n_long
	bloom.boundary_index = bloom.n_long * bloom.len_long
	return &bloom
}

// `Add`: add a value to the set
func (b *BloomShard) Add(value any) {
	// find the indices
	indices := b.getIndices(value)

	// find word index and offset, and set it to true
	for _, index := range indices {
		wi := index / 64
		off := index % 64

		si := b.getShardId(index)
		b.shards[si].Lock()
		b.filter[wi] |= (1 << off)
		b.shards[si].Unlock()
	}
}

// `Check`: check a value to the set (false negative: never, false positives: maybe)
func (b *BloomShard) Check(value any) bool {
	// find the indices
	indices := b.getIndices(value)

	// find word index and offset, and check if it is false
	for _, index := range indices {
		wi := index / 64
		off := index % 64

		si := b.getShardId(index)
		b.shards[si].RLock()
		is_reset := ((b.filter[wi] & (1 << off)) == 0)
		b.shards[si].RUnlock()

		if is_reset {
			return false
		}
	}

	return true
}

// `getIndices`: find filter indices
func (b *BloomShard) getIndices(value any) []uint64 {
	// get bytes
	data := toBytes(value)

	// get primary hashes
	h1 := hash(b.seeds[0], data) % b.n_bits
	h2 := hash(b.seeds[1], data) % b.n_bits

	// use double hashing to generate n_hash indices
	indices := make([]uint64, b.n_hash)
	m := b.n_bits

	for i := uint64(0); i < uint64(b.n_hash); i++ {
		index := (h1 + (i*h2)%m) % m
		indices[i] = index
	}

	return indices
}

// `getShardId`: find the shard id for a given index
func (b *BloomShard) getShardId(idx uint64) uint64 {

	// shard layout
	// [0, n_long-1] - longer shards
	// [n_long, n_shards-1] - shorter shards

	if idx < b.boundary_index {
		return idx / b.len_long
	}
	return b.n_long + (idx-b.boundary_index)/b.len_short
}

func (b *BloomShard) Reset() {
	for i := range b.filter {
		b.filter[i] = 0
	}
}

// complie-time check
var _ IBloom = (*BloomShard)(nil)
