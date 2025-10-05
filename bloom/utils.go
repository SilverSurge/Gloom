package bloom

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/twmb/murmur3"
)

const (
	DefaultSeed1 uint32 = 6269
	DefaultSeed2 uint32 = 4241
)

// `hash`: returns murmur3 using a seed and data
func hash(seed uint32, data []byte) uint32 {
	// hash function
	return murmur3.SeedSum32(seed, data)
}

// `toBytes`: converts any to []byte
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

// `PositiveProbablityEstimate`: return probability of returning true, for a filter with n_bits, n_hash, and n_add Add operations.
func FalsePositiveProbabilityEstimate(n_bits, n_hash, n_add uint32) float64 {
	b := float64(n_bits)
	h := float64(n_hash)
	a := float64(n_add)

	p := math.Pow(float64(1)-math.Pow(float64(1)-h/b, a), h)

	return p
}
