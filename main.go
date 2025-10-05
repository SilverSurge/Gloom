package main

import (
	"fmt"
	"gloom/bloom"
)

func main() {

	// create a bloom filter
	id := "default-bloom-filter-id"
	n_bits := uint32(64 * 1e5)
	n_hash := uint32(8)
	n_add := uint32(1e6)

	// estimate
	fmt.Printf("FalsePositiveProbabilityEstimate: %.2f %%\n", 100*bloom.FalsePositiveProbabilityEstimate(n_bits, n_hash, n_add))

	bf := bloom.NewBloomDefault(id, n_bits, n_hash)
	// bf := bloom.NewBloomCustom(id, n_bits, n_hash, [2]uint32{111,222})

	// adding data
	data1 := "silver"
	data2 := struct {
		name       string
		species    string
		isFriendly bool
	}{
		"Rex",
		"GoldenRetriever",
		true,
	}
	data3 := 1e9 + 7

	bf.Add(data1)
	bf.Add(data2)
	bf.Add(data3)

	// checking data
	fmt.Println(bf.Check("palladium"))
	fmt.Println(bf.Check(121))
	fmt.Println(bf.Check(data2))

}
