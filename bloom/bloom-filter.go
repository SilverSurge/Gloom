package bloom

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/twmb/murmur3"
)

// -----------------------------------------------------
//
//	Structs
//
// -----------------------------------------------------
type Bloom struct {
	id     string
	filter []bool
	n_bits uint32
	n_hash uint32
	seeds  [2]uint32
}

// -----------------------------------------------------
//
//	Constants
//
// -----------------------------------------------------
const (
	DefaultSeed1 uint32 = 6269
	DefaultSeed2 uint32 = 4241
)

// -----------------------------------------------------
//
//	Functions and Methods
//
// -----------------------------------------------------
func NewBloomDefault(id string, n_bits, n_hash uint32) Bloom {
	return NewBloomCustom(id, n_bits, n_hash, [2]uint32{DefaultSeed1, DefaultSeed2})
}

func NewBloomCustom(id string, n_bits, n_hash uint32, seeds [2]uint32) Bloom {
	bloom := Bloom{
		id:     id,
		n_bits: n_bits,
		n_hash: n_hash,
		seeds:  seeds,
		filter: make([]bool, n_bits),
	}
	return bloom
}

func (b *Bloom) Add(value any) {
	// find the indices
	indices := b.getIndices(value)

	// set them to true
	for _, index := range indices {
		b.filter[index] = true
	}
}

func (b *Bloom) Check(value any) bool {
	// find the indices
	indices := b.getIndices(value)

	// check if false
	for _, index := range indices {
		if !b.filter[index] {
			return false
		}
	}

	return true
}

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

// -----------------------------------------------------
// 	Helper Functions
// -----------------------------------------------------

func hash(seed uint32, data []byte) uint32 {
	return murmur3.SeedSum32(seed, data)
}

func toBytes(value any) []byte {
	switch v := value.(type) {
	case nil:
		return []byte("null")

	case []byte:
		return v

	case string:
		return []byte(v)

	case int, int8, int16, int32, int64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(reflect.ValueOf(v).Int()))
		return buf

	case uint, uint8, uint16, uint32, uint64:
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, reflect.ValueOf(v).Uint())
		return buf

	case float32:
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.BigEndian, v)
		return buf.Bytes()

	case float64:
		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.BigEndian, v)
		return buf.Bytes()

	case fmt.Stringer:
		return []byte(v.String())

	default:
		// fallback: deterministic JSON encoding
		data, err := json.Marshal(v)
		if err != nil {
			// as a last resort, use fmt.Sprintf (non-deterministic but safe)
			return []byte(fmt.Sprintf("%v", v))
		}
		return data
	}
}
