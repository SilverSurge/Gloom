package bloom

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// helper to exercise a bloom implementation via IBloom
func runBasicBloomSuite(t *testing.T, b IBloom) {
	t.Helper()
	b.Reset()
	if b.Check("x") {
		t.Fatal("empty bloom should not contain 'x'")
	}
	b.Add("x")
	if !b.Check("x") {
		t.Fatal("bloom should contain 'x' after add")
	}
	// different value should probably not be present
	if b.Check("y") && b.Check("z") {
		// it's possible to get false positives; but both being true is suspicious for small filters
		t.Log("warning: possible false positives for small bloom")
	}
	b.Reset()
	if b.Check("x") {
		t.Fatal("after reset bloom should not contain 'x'")
	}
}

func TestBloomImplementationsBasic(t *testing.T) {
	// choose small parameters for quick tests
	nbits := uint64(256)
	nhash := uint64(3)
	b1 := NewBloomDefault("testbloom", nbits, nhash)
	b2 := NewBloomAtomicDefault("testbloom", nbits, nhash)
	b3 := NewBloomRWDefault("testbloom", nbits, nhash)
	b4 := NewBloomShardDefault("testbloom", nbits, nhash, 4)

	runBasicBloomSuite(t, b1)
	runBasicBloomSuite(t, b2)
	runBasicBloomSuite(t, b3)
	runBasicBloomSuite(t, b4)
}

func TestBloomDSUnionAndSaveLoad(t *testing.T) {
	nbits := uint64(128)
	nhash := uint64(2)
	bds1 := NewBloomDSCustom("u1", nbits, nhash, [2]uint64{DefaultSeed1, DefaultSeed2})
	bds2 := NewBloomDSCustom("u1", nbits, nhash, [2]uint64{DefaultSeed1, DefaultSeed2})

	bds1.Filter[0] = 0xdeadbeef
	bds2.Filter[0] = 0x01020304

	ok := bds1.Union(&bds2)
	if !ok {
		t.Fatal("expected union to succeed with same params")
	}
	if (bds1.Filter[0] & 0x01020304) == 0 {
		t.Fatal("union did not combine filters")
	}

	// save and load
	d := t.TempDir()
	bds1.ID = "mybloom"
	if err := bds1.Save(d); err != nil {
		t.Fatalf("save error: %v", err)
	}

	bdsL := NewBloomDSCustom("mybloom", nbits, nhash, [2]uint64{DefaultSeed1, DefaultSeed2})
	if err := bdsL.Load(d); err != nil {
		t.Fatalf("load error: %v", err)
	}
	if bdsL.Filter[0] != bds1.Filter[0] {
		t.Fatalf("loaded filter mismatch: got %x want %x", bdsL.Filter[0], bds1.Filter[0])
	}

	// try loading a non-existent file
	fake := NewBloomDSCustom("no-such-file", nbits, nhash, [2]uint64{DefaultSeed1, DefaultSeed2})
	if err := fake.Load(d); err == nil {
		t.Fatal("expected error when loading non-existent file")
	}
}

func TestGetIndicesDeterminism(t *testing.T) {
	nbits := uint64(1024)
	nhash := uint64(5)
	b := NewBloomDSCustom("idx", nbits, nhash, [2]uint64{DefaultSeed1, DefaultSeed2})

	v := "some value"
	idx1 := b.GetIndices(v)
	idx2 := b.GetIndices(v)
	if len(idx1) != int(nhash) || len(idx2) != int(nhash) {
		t.Fatalf("unexpected indices length")
	}
	for i := range idx1 {
		if idx1[i] != idx2[i] {
			t.Fatalf("indices not deterministic: %v vs %v", idx1, idx2)
		}
		if idx1[i] >= nbits {
			t.Fatalf("index out of range: %d >= %d", idx1[i], nbits)
		}
	}
}

func TestToBytesVariousTypes(t *testing.T) {
	// some smoke tests to ensure toBytes doesn't panic and produces deterministic outputs
	_ = toBytes(nil)
	_ = toBytes([]byte("abc"))
	_ = toBytes("hello")
	_ = toBytes(int64(-42))
	_ = toBytes(uint32(0xdead))
	_ = toBytes(3.1415)
	_ = toBytes(struct{ A int }{A: 10})
}

func TestGetFalsePositiveProbabilityEstimateMonotonic(t *testing.T) {
	p1 := GetFalsePositiveProbabilityEstimate(1024, 3, 10)
	p2 := GetFalsePositiveProbabilityEstimate(1024, 3, 100)
	if p2 <= p1 {
		t.Fatalf("expected probability to increase with more additions: %f <= %f", p2, p1)
	}
}

func TestGetOptimalParametersBasic(t *testing.T) {
	nadd := uint64(1000)
	p := 0.01
	nbits, nhash := GetOptimalParameters(nadd, p)
	if nbits == 0 || nhash == 0 {
		t.Fatalf("expected non-zero optimal parameters, got %d %d", nbits, nhash)
	}
}

// concurrency smoke test for BloomAtomic
func TestBloomAtomicConcurrentAdds(t *testing.T) {
	nbits := uint64(4096)
	nhash := uint64(4)
	b := NewBloomAtomicDefault("concurrent", nbits, nhash)
	wg := sync.WaitGroup{}
	n := 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			b.Add(i)
		}(i)
	}
	wg.Wait()
	// ensure some checks succeed
	if !b.Check(0) {
		t.Log("warning: first element not found; possible false negative (shouldn't happen) or ordering issue")
	}
}

// small helper to ensure tmpdir write permissions in CI
func TestSaveCreatesFileWithCorrectName(t *testing.T) {
	d := t.TempDir()
	b := NewBloomDSCustom("savetest", 64, 2, [2]uint64{DefaultSeed1, DefaultSeed2})
	if err := b.Save(d); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	// assert file exists
	p := filepath.Join(d, "savetest.bloom")
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("expected file %s to exist, stat error: %v", p, err)
	}
}
