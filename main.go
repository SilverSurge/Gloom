package main

import (
	"fmt"
	"os"

	"github.com/SilverSurge/Gloom/bloom"
)

func main() {

	// find bloom filter params
	id := "default-bloom-filter-id"

	n_add := uint64(1e6)
	prob_fp := float64(0.01)

	n_bits, n_hash := bloom.GetOptimalParameters(n_add, prob_fp)

	n_shards := n_bits / 16

	// find the false positive probability estimate
	fmt.Printf("n_add: %d\n", n_add)
	fmt.Printf("prob_fp: %.2f\n", prob_fp)
	fmt.Printf("optimal n_bits: %d\n", n_bits)
	fmt.Printf("optimal n_hash: %d\n", n_hash)
	fmt.Printf("false positive probability estimate: %.4f %%\n", 100*bloom.GetFalsePositiveProbabilityEstimate(n_bits, n_hash, n_add))

	// create bloom filters
	bf1 := bloom.NewBloomDefault(id, n_bits, n_hash)
	bf2 := bloom.NewBloomRWDefault(id, n_bits, n_hash)
	bf3 := bloom.NewBloomAtomicDefault(id, n_bits, n_hash)
	bf4 := bloom.NewBloomShardDefault(id, n_bits, n_hash, n_shards)
	// bf := bloom.NewBloomCustom(id, n_bits, n_hash, [2]uint64{111,222})

	// create data
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
	data4 := 3.14159265359

	// add data
	bf1.Add(data1)
	bf2.Add(data2)
	bf3.Add(data3)
	bf4.Add(data4)

	// apply union
	bf1.Union(&bf2.State)
	bf1.Union(&bf3.State)

	// save and load bloom_ds

	// save the bloom filter state
	err := bf1.State.Save("./save_dir")
	if err != nil {
		panic(err)
	}

	// load the bloom filter state
	bds := bloom.BloomDS{ID: id} // make sure that you have the same id
	err = bds.Load("./save_dir")
	if err != nil {
		panic(err)
	}

	// simple clean up
	defer os.RemoveAll("./save_dir")

	// create a bloom filter using the loaded state
	bf := bloom.NewBloomFromBloomDS(&bds)

	fmt.Println("not found 1:", bf.Check("palladium"))
	fmt.Println("not found 2:", bf.Check(121))
	fmt.Println("already present 1:", bf.Check(data1))
	fmt.Println("already present 2:", bf.Check(data2))
}
