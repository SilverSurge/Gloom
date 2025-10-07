# 🌫️ Gloom — A Simple Bloom Filter in Go

**Gloom** is a lightweight, generic Bloom filter implementation in Go, designed for simplicity and clarity. It provides basic operations to add and check membership of arbitrary data types using efficient hashing with **Murmur3**.

## ✨ Features

- ✅ Generic support for any Go type (`any`)
- ⚡ Fast double hashing using Murmur3
- 🧩 Customizable bit size, number of hash functions, and seeds
- 🧠 Deterministic byte conversion with JSON fallback
- 🧍 Minimal dependencies and easy to use

## 📦 Installation

```bash
go get github.com/SilverSurge/Gloom
```

---

## 🧰 Usage

Here’s a simple example, a detailed example can be found in `main.go`:

```go
package main

import (
    "fmt"
    "github.com/SilverSurge/Gloom/bloom"
)

func main() {
    // Create a new Bloom filter with default seeds
    bf := bloom.NewBloomDefault("example", 1024, 3)

    // Add some values
    bf.Add("apple")
    bf.Add(42)

    // Check for membership
    fmt.Println(bf.Check("apple")) // true
    fmt.Println(bf.Check("banana")) // false
    fmt.Println(bf.Check(42)) // true
}
```

---

## ⚙️ Customization

You can also define custom seeds for more control over hash generation:

```go
bf := bloom.NewBloomCustom("custom", 2048, 4, [2]uint32{1111, 2222})
```

| Parameter | Description                                |
| --------- | ------------------------------------------ |
| `id`      | Optional name or identifier for the filter |
| `n_bits`  | Number of bits in the filter (size)        |
| `n_hash`  | Number of hash functions used              |
| `seeds`   | Two seeds for Murmur3 double hashing       |

---

## 🧪 Implementation Details

* Uses **double hashing** to derive multiple indices from two primary Murmur3 hashes.
* The conversion function `toBytes()` handles all common Go types deterministically.
* False positives are possible (as in all Bloom filters), but false negatives are not.


## 🧩 Dependencies

* [`github.com/twmb/murmur3`](https://pkg.go.dev/github.com/twmb/murmur3) — fast Murmur3 hashing library.

## 💡 Future Ideas

* ~~Add serialization/deserialization support~~✅
* ~~Introduce concurrency-safe version~~ ✅
* ~~Provide false-positive probability estimator~~✅
* add version headers to the save files

## 🌱 Inspiration

This project was built as a simple exploration of Bloom filters and hashing concepts in Go.
It’s intentionally minimal for clarity and educational purposes.

## ⚖️ License

Distributed under the MIT License. See `LICENSE` for more information.

## 📜 Change Log

1. Replace `[]bool` with `[]uint64` for memory optimization.
2. Add concurrency safe versions: `BloomRW`, `BloomAtomic`, and `BloomShard`.
3. Add `GetOptimalParameters` and `GetFalsePositiveProbabilityEstimate`.
4. Add `Save` and `Load` for `BloomDS`.

## 🗎 Documentation

```go
package bloom // import "gloom/bloom"


CONSTANTS

const (
        DefaultSeed1 uint64 = 6269
        DefaultSeed2 uint64 = 4241
)

FUNCTIONS

func GetFalsePositiveProbabilityEstimate(n_bits, n_hash, n_add uint64) float64
    """`GetPositiveProbablityEstimate`: returns probability of getting true,
    for a filter with n_bits, n_hash, and n_add Add operations."""

func GetOptimalParameters(n_add uint64, prob_fp float64) (uint64, uint64)
    """`GetOptimalParameters`: return optimal (n_bits, n_hash) for a given n_add,
    and prob_fp"""


TYPES

type BloomDS struct {
        ID     string
        NBits  uint64
        NHash  uint64
        Seeds  [2]uint64
        Filter []uint64
}

func NewBloomDSCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) BloomDS
    """`NewBloomDSCustom`: return custom bloom_ds"""

func NewBloomDSDefault(id string, n_bits, n_hash uint64) BloomDS
    """`NewBloomDSDefault`: return default bloom_ds"""

func (b *BloomDS) GetIndices(value any) []uint64
    """`GetIndices`: get indices that would be considered for a value"""

func (b *BloomDS) Load(dir string) error
    """`Load`: load bloom_ds from dir/id.bloom"""

func (b *BloomDS) Reset()
    """`Reset`: resets all bits"""

func (b *BloomDS) Save(dir string) error
    """`Save`: save bloom_ds to dir/id.bloom"""

func (b1 *BloomDS) Union(b2 *BloomDS) bool
    """`Union`: union with another bloom_ds with same n_bits and seeds"""


type Bloom struct {
        State BloomDS
}

func NewBloomCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) *Bloom
    """`NewBloomCustom` return a custom `Bloom` object"""

func NewBloomDefault(id string, n_bits, n_hash uint64) *Bloom
    """`NewBloomDefault` return a default `Bloom` object"""

func NewBloomFromBloomDS(b *BloomDS) *Bloom
    """`NewBloomFromBloomDS`: return a `Bloom` using the data from bloom_ds"""

func (b *Bloom) Add(value any)
    """`Add`: add a value to the set"""

func (b *Bloom) Check(value any) bool
    """`Check`: check a value to the set (false negative: never, false positives:
    maybe)"""

func (b *Bloom) Reset()
    """`Reset`: resets bloom_ds"""

func (b1 *Bloom) Union(b2 *BloomDS) bool
    """`Union`: tries state union"""

type BloomAtomic struct {
        State BloomDS
}

func NewBloomAtomicCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) *BloomAtomic
    """`NewBloomAtomicCustom` return a custom `BloomAtomic` object"""

func NewBloomAtomicDefault(id string, n_bits, n_hash uint64) *BloomAtomic
    """`NewBloomAtomicDefault` return a default `BloomAtomic` object"""

func NewBloomAtomicFromBloomDS(b *BloomDS) *BloomAtomic
    """`NewBloomAtomicFromBloomDS`: return a `BloomAtomic` using the data from the
    bl"""oom_ds

func (b *BloomAtomic) Add(value any)
    """`Add`: add a value to the set"""

func (b *BloomAtomic) Check(value any) bool
    """`Check`: check a value to the set (false negative: never, false positives:
    ma"""ybe)

func (b *BloomAtomic) Reset()
    """`Reset`: resets bloom_ds"""

func (b1 *BloomAtomic) Union(b2 *BloomDS) bool
    """`Union`: tries state union"""


type BloomRW struct {
        State BloomDS
        Mu    sync.RWMutex
}

func NewBloomRWCustom(id string, n_bits, n_hash uint64, seeds [2]uint64) *BloomRW
    """`NewBloomRWCustom` return a custom `BloomRW` object"""

func NewBloomRWDefault(id string, n_bits, n_hash uint64) *BloomRW
    """`NewBloomRWDefault` return a default `BloomRW` object"""

func NewBloomRWFromBloomDS(b *BloomDS) *BloomRW
    """`NewBloomRWFromBloomDS`: return a `BloomRW` using the data from the bloom_ds"""

func (b *BloomRW) Add(value any)
    """`Add`: add a value to the set"""

func (b *BloomRW) Check(value any) bool
    """`Check`: check a value to the set (false negative: never, false positives:
    ma"""ybe)

func (b *BloomRW) Reset()
    """`Reset`: resets bloom_ds"""

func (b1 *BloomRW) Union(b2 *BloomDS) bool
    """`Union`: tries state union"""

type BloomShard struct {
        State   BloomDS
        NShards uint64
        Shards  []sync.RWMutex

        // Has unexported fields.
}

func NewBloomShardCustom(id string, n_bits, n_hash, n_shards uint64, seeds [2]uint64) *BloomShard
    """`NewBloomShardCustom` return a custom `BloomShard` object"""

func NewBloomShardDefault(id string, n_bits, n_hash, n_shards uint64) *BloomShard
    """`NewBloomShardDefault` return a default `BloomShard` object"""

func NewBloomShardFromBloomDS(b *BloomDS, n_shard uint64) *BloomShard
    """`NewBloomShardFromBloomDS`: return a `BloomShard` using the data from the
    bl"""oom_ds

func (b *BloomShard) Add(value any)
    """`Add`: add a value to the set"""

func (b *BloomShard) Check(value any) bool
    """`Check`: check a value to the set (false negative: never, false positives:
    ma"""ybe)

func (b *BloomShard) Reset()
    """`Reset`: resets bloom_ds"""

func (b1 *BloomShard) Union(b2 *BloomDS) bool
    """`Union`: tries state union"""

type IBloom interface {
        Add(any)
        Check(any) bool
        Reset()
        Union(*BloomDS) bool
}
```