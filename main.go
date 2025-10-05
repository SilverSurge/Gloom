package main

import (
	"fmt"
	"gloom/bloom"
)

func main() {

	// create a bloom filter
	id := "default-bloom-filter-id"

	n_add := uint64(1e6)
	prob_fp := float64(0.01)

	n_bits, n_hash := bloom.GetOptimalParameters(n_add, prob_fp)

	// estimate
	fmt.Printf("n_add: %d\n", n_add)
	fmt.Printf("prob_fp: %.2f\n", prob_fp)
	fmt.Printf("optimal n_bits: %d\n", n_bits)
	fmt.Printf("optimal n_hash: %d\n", n_hash)

	fmt.Printf("false positive probability estimate: %.4f %%\n", 100*bloom.GetFalsePositiveProbabilityEstimate(n_bits, n_hash, n_add))

	bf := bloom.NewBloomShardDefault(id, n_bits, n_hash, 64)
	// bf := bloom.NewBloomCustom(id, n_bits, n_hash, [2]uint64{111,222})

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
	fmt.Println("not found 1:", bf.Check("palladium"))
	fmt.Println("not found 2:", bf.Check(121))
	fmt.Println("already present:", bf.Check(data2))
}
