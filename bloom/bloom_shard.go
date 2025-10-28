package bloom

import (
	"sync"
)

type BloomShard struct {
	State   BloomDS
	NShards uint64
	Shards  []sync.RWMutex

	len_long       uint64
	len_short      uint64
	n_long         uint64
	n_short        uint64
	boundary_index uint64
	rareMu         sync.RWMutex
}

// `NewBloomShardDefault` return a default `BloomShard` object
func NewBloomShardDefault(id string, n_bits, n_hash, n_shards uint64) *BloomShard {
	return NewBloomShardCustom(id, n_bits, n_hash, n_shards, [2]uint64{DefaultSeed1, DefaultSeed2})
}

// `NewBloomShardCustom` return a custom `BloomShard` object
func NewBloomShardCustom(id string, n_bits, n_hash, n_shards uint64, seeds [2]uint64) *BloomShard {

	bloom := BloomShard{
		State:   NewBloomDSCustom(id, n_bits, n_hash, seeds),
		NShards: min(n_shards, n_bits),
		Shards:  make([]sync.RWMutex, n_shards),

		len_long:  (n_bits + n_shards - 1) / n_shards,
		len_short: n_bits / n_shards,
		n_long:    n_bits % n_shards,
	}
	bloom.n_short = bloom.NShards - bloom.n_long
	bloom.boundary_index = bloom.n_long * bloom.len_long
	return &bloom
}

// `NewBloomShardFromBloomDS`: return a `BloomShard` using the data from the bloom_ds
func NewBloomShardFromBloomDS(b *BloomDS, n_shard uint64) *BloomShard {
	bloom := NewBloomShardCustom(b.ID, b.NBits, b.NHash, n_shard, b.Seeds)
	bloom.Union(b)
	return bloom
}

// `Add`: add a value to the set
func (b *BloomShard) Add(value any) {
	// unionRW mutex
	b.rareMu.RLock()
	defer b.rareMu.RUnlock()

	// find the indices
	indices := b.State.GetIndices(value)

	// find word index and offset, and set it to true
	for _, index := range indices {
		wi := index / 64
		off := index % 64

		si := b.getShardId(index)
		b.Shards[si].Lock()
		b.State.Filter[wi] |= (1 << off)
		b.Shards[si].Unlock()
	}
}

// `Check`: check a value to the set (false negative: never, false positives: maybe)
func (b *BloomShard) Check(value any) bool {
	// unionRW mutex
	b.rareMu.RLock()
	defer b.rareMu.RUnlock()

	// find the indices
	indices := b.State.GetIndices(value)

	// find word index and offset, and check if it is false
	for _, index := range indices {
		wi := index / 64
		off := index % 64

		si := b.getShardId(index)
		b.Shards[si].RLock()
		is_reset := ((b.State.Filter[wi] & (1 << off)) == 0)
		b.Shards[si].RUnlock()

		if is_reset {
			return false
		}
	}

	return true
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

// `Reset`: resets bloom_ds
func (b *BloomShard) Reset() {
	// unionRW mutex
	b.rareMu.Lock()
	defer b.rareMu.Unlock()

	b.State.Reset()
}

// `Union`: tries state union
func (b1 *BloomShard) Union(b2 *BloomDS) bool {
	// unionRW mutex
	b1.rareMu.RLock()
	defer b1.rareMu.RUnlock()

	return b1.State.Union(b2)
}

// `GetState`: return current State bool
func (b *BloomShard) GetState() BloomDS {
	// unionRW mutex
	b.rareMu.Lock()
	defer b.rareMu.Unlock()

	return b.State
}

// complie-time check
var _ IBloom = (*BloomShard)(nil)
