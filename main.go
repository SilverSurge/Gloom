package main

import (
	"fmt"
	"gloom/bloom"
)

func main() {

	// create a bloom filter
	id := "default-bloom-filter-id"
	n_bits := uint32(2 * 1e5)
	n_hash := uint32(64)

	bf := bloom.NewBloomDefault(id, n_bits, n_hash)

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
