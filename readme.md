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

Here’s a simple example:

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

### Internal Design

```
Add(value)
 ├─→ getIndices(value)
 │    ├─→ toBytes(value)
 │    ├─→ hash(seed1, data)
 │    └─→ hash(seed2, data)
 └─→ mark bits as true

Check(value)
 ├─→ getIndices(value)
 └─→ verify all bits set
```

## 🧩 Dependencies

* [`github.com/twmb/murmur3`](https://pkg.go.dev/github.com/twmb/murmur3) — fast Murmur3 hashing library.

## 💡 Future Ideas

* Add serialization/deserialization support
* ~~Introduce concurrency-safe version~~ ✅
* ~~Provide false-positive probability estimator~~

## 🌱 Inspiration

This project was built as a simple exploration of Bloom filters and hashing concepts in Go.
It’s intentionally minimal for clarity and educational purposes.

## ⚖️ License

Distributed under the MIT License. See `LICENSE` for more information.

## 📜 Change Log

1. Replace `[]bool` with `[]uint64` for memory optimization.
2. Add concurrency safe versions: `BloomRW`, `BloomAtomic`, and `BloomShard`.