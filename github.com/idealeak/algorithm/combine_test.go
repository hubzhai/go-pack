package algorithm

import "testing"

func TestCombNumber(t *testing.T) {
	testcase := []struct {
		n uint64
		m uint64
		r uint64
	}{
		{2, 1, 2},
		{5, 3, 10},
		{7, 5, 21},
	}
	for i := 0; i < len(testcase); i++ {
		cn := CombNumber(testcase[i].n, testcase[i].m)
		if cn != testcase[i].r {
			t.Errorf("[fail] C(%d,%d)=%d, but result=%d", testcase[i].n, testcase[i].m, testcase[i].r, cn)
		} else {
			t.Logf("[ok] C(%d,%d)=%d", testcase[i].n, testcase[i].m, testcase[i].r)
		}
	}
}

func BenchmarkCombNumber(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CombNumber(52, 7)
	}
}

func TestCombinerSelectNoRecursion(t *testing.T) {
	testcase := []struct {
		atable []int32
		n      int
	}{
		{[]int32{0, 1, 2, 3, 4, 5, 6}, 5},
	}
	for i := 0; i < len(testcase); i++ {
		result := CombinerSelectNoRecursion(testcase[i].atable, testcase[i].n)
		for i := 0; i < len(result); i++ {
			t.Log(i, result[i])
		}
		t.Logf("Combiner(%v,%v)=cnt=%v \n", testcase[i].atable, testcase[i].n, len(result))
	}
}

func TestCombinerSelectUseRecursion(t *testing.T) {
	testcase := []struct {
		atable []int32
		n      int
	}{
		{[]int32{0, 1, 2, 3, 4, 5, 6}, 5},
	}
	var ints []int32
	for i := 0; i < 43; i++ {
		ints = append(ints, int32(i))
	}
	for i := 0; i < len(testcase); i++ {
		result := CombinerSelectUseRecursion(ints, 2)
		for i := 0; i < len(result); i++ {
			t.Log(i, result[i])
		}
		t.Logf("Combiner2(%v,%v)=cnt=%v \n", testcase[i].atable, testcase[i].n, len(result))
	}
}

func BenchmarkCombinerSelectNoRecursion(b *testing.B) {
	atable := make([]int32, 45)
	for i := 0; i < 45; i++ {
		atable[i] = int32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CombinerSelectNoRecursion(atable, 2)
	}
}

func BenchmarkCombinerSelectUseRecursion(b *testing.B) {
	atable := make([]int32, 45)
	for i := 0; i < 45; i++ {
		atable[i] = int32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CombinerSelectUseRecursion(atable, 2)
	}
}

func TestFullPermutation(t *testing.T) {
	var test = []int32{0, 1, 2, 3, 4}
	r := FullPermutation(test)
	for i := 0; i < len(r); i++ {
		t.Log("FullPermutation", i, "===", r[i])
	}
}

func BenchmarkFullPermutation(b *testing.B) {
	var test = []int32{0, 1, 2, 3, 4}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FullPermutation(test)
	}
}