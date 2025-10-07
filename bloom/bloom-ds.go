package bloom

import (
	"encoding/gob"
	"os"
	"path/filepath"
)

type BloomDS struct {
	ID     string
	NBits  uint64
	NHash  uint64
	Seeds  [2]uint64
	Filter []uint64
}

// `NewBloomDSDefault`: return default bloom_ds
func NewBloomDSDefault(id string, n_bits, n_hash uint64) BloomDS {
	return NewBloomDSCustom(id, n_bits, n_hash, [2]uint64{DefaultSeed1, DefaultSeed2})
}

// `NewBloomDSCustom`: return custom bloom_ds
func NewBloomDSCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) BloomDS {
	return BloomDS{
		ID:     id,
		NBits:  n_bits,
		NHash:  n_hash,
		Seeds:  seeds,
		Filter: make([]uint64, (n_bits+63)/64),
	}
}

// `Reset`: resets all bits
func (b *BloomDS) Reset() {
	for i := range b.Filter {
		b.Filter[i] = 0
	}
}

// `Union`: union with another bloom_ds with same n_bits and seeds
func (b1 *BloomDS) Union(b2 *BloomDS) bool {
	if b1.NBits != b2.NBits || b1.NHash != b2.NHash || b1.Seeds != b2.Seeds {
		return false
	}
	for i := range b1.Filter {
		b1.Filter[i] |= b2.Filter[i]
	}
	return true
}

// `GetIndices`: get indices that would be considered for a value
func (b *BloomDS) GetIndices(value any) []uint64 {
	// get bytes
	data := toBytes(value)

	// get primary hashes
	h1 := hash(b.Seeds[0], data) % b.NBits
	h2 := hash(b.Seeds[1], data) % b.NBits

	// use double hashing to generate n_hash indices
	indices := make([]uint64, b.NHash)
	m := b.NBits

	for i := uint64(0); i < uint64(b.NHash); i++ {
		index := (h1 + (i*h2)%m) % m
		indices[i] = index
	}

	return indices
}

// `Save`: save bloom_ds to dir/id.bloom
func (b *BloomDS) Save(dir string) error {
	// make dir if it doesnt exist
	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return err
	}

	// save to the file
	fname := filepath.Join(dir, b.ID+".bloom")
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(b)
}

// `Load`: load bloom_ds from dir/id.bloom
func (b *BloomDS) Load(dir string) error {
	fname := filepath.Join(dir, b.ID+".bloom")
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(b)
}
