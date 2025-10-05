package bloom

type Bloom struct {
	id     string
	n_bits uint32
	n_hash uint32
	seeds  [2]uint32
	filter []uint64
}

// `NewBloomDefault` return a default `Bloom` object
func NewBloomDefault(id string, n_bits, n_hash uint32) Bloom {
	return NewBloomCustom(id, n_bits, n_hash, [2]uint32{DefaultSeed1, DefaultSeed2})
}

// `NewBloomDefault` return a custom `Bloom` object
func NewBloomCustom(id string, n_bits, n_hash uint32, seeds [2]uint32) Bloom {

	n_words := (n_bits + 63) / 64
	bloom := Bloom{
		id:     id,
		n_bits: n_bits,
		n_hash: n_hash,
		seeds:  seeds,
		filter: make([]uint64, n_words),
	}
	return bloom
}

// `Add`: add a value to the set
func (b *Bloom) Add(value any) {
	// find the indices
	indices := b.getIndices(value)

	// find word index and offset, and set it to true
	for _, index := range indices {
		wi := index / 64
		off := index % 64
		b.filter[wi] |= (1 << off)
	}
}

// `Add`: check a value to the set (false negative: never, false positives: maybe)
func (b *Bloom) Check(value any) bool {
	// find the indices
	indices := b.getIndices(value)

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
func (b *Bloom) getIndices(value any) []uint32 {
	// get bytes
	data := toBytes(value)

	// get primary hashes
	h1 := hash(b.seeds[0], data) % b.n_hash
	h2 := hash(b.seeds[1], data) % b.n_hash

	// use double hashing to generate n_hash indices
	inidices := make([]uint32, b.n_hash)
	m := b.n_bits

	for i := uint32(0); i < uint32(b.n_hash); i++ {
		index := (h1 + ((i*h2)%m)%m + m) % m
		inidices[i] = index
	}

	return inidices
}

// complie-time check
var _ IBloom = (*Bloom)(nil)
