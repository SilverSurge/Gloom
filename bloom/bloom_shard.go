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
	}
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
		index := (h1 + ((i*h2)%m)%m + m) % m
		indices[i] = index
	}

	return indices
}

// `getShardId`: find the shard id for a given index
func (b *BloomShard) getShardId(idx uint64) uint64 {
	// setup
	n_bits := b.n_bits     // number of bits
	n_shards := b.n_shards // number of shards

	// solution

	len_short := n_bits / n_shards                 // shorter shard length
	len_long := (n_bits + n_shards - 1) / n_shards // longer shard length
	n_long := n_bits % n_shards                    // number of longer shards

	// shard layout
	// [0, n_long-1] - longer shards
	// [n_long, n_shards-1] - shorter shards

	boundary_index := n_long * len_long

	if idx < boundary_index {
		return idx / len_long
	} else {
		relative_index := idx - boundary_index

		if len_short == 0 {
			// this should never happen
			panic("BloomShard::getShardId: shorter shard length is 0\n")
		}

		shorter_shard_offset := relative_index / len_short
		return n_long + shorter_shard_offset
	}
}

// complie-time check
var _ IBloom = (*BloomShard)(nil)
