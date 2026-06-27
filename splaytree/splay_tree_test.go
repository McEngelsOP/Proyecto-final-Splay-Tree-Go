package splaytree

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func TestInsertAndSearch(t *testing.T) {
	tree := New[int, string]()
	tree.Insert(10, "ten")
	tree.Insert(5, "five")
	tree.Insert(15, "fifteen")
	tree.Insert(3, "three")
	tree.Insert(7, "seven")

	cases := []struct {
		key   int
		value string
		found bool
	}{
		{10, "ten", true},
		{5, "five", true},
		{15, "fifteen", true},
		{3, "three", true},
		{7, "seven", true},
		{99, "", false},
		{0, "", false},
	}

	for _, tc := range cases {
		v, ok := tree.Search(tc.key)
		if ok != tc.found {
			t.Errorf("Search(%d): found=%v, want %v", tc.key, ok, tc.found)
		}
		if ok && v != tc.value {
			t.Errorf("Search(%d): value=%q, want %q", tc.key, v, tc.value)
		}
	}
}

func TestSplayToRoot(t *testing.T) {
	tree := New[int, string]()
	for _, k := range []int{10, 5, 15, 3, 7, 12, 20} {
		tree.Insert(k, fmt.Sprintf("v%d", k))
	}

	targets := []int{7, 15, 3, 12}
	for _, key := range targets {
		tree.Search(key)
		root, _ := tree.Root()
		if root != key {
			t.Errorf("tras Search(%d), raíz=%d, se esperaba %d", key, root, key)
		}
	}
}

func TestInsertSplaysToRoot(t *testing.T) {
	tree := New[int, string]()
	keys := []int{10, 5, 15, 3, 7}
	for _, k := range keys {
		tree.Insert(k, "")
		root, _ := tree.Root()
		if root != k {
			t.Errorf("tras Insert(%d), raíz=%d, se esperaba %d", k, root, k)
		}
	}
}

func TestUpdateExistingKey(t *testing.T) {
	tree := New[string, int]()
	tree.Insert("85123A", 10)
	tree.Insert("85123A", 999)

	v, ok := tree.Search("85123A")
	if !ok {
		t.Fatal("clave '85123A' no encontrada tras actualización")
	}
	if v != 999 {
		t.Errorf("valor=%d, want 999", v)
	}
	if tree.Size() != 1 {
		t.Errorf("size=%d, want 1", tree.Size())
	}
}

func TestDelete(t *testing.T) {
	tree := New[int, string]()
	keys := []int{10, 5, 15, 3, 7, 12, 20}
	for _, k := range keys {
		tree.Insert(k, fmt.Sprintf("v%d", k))
	}

	if !tree.Delete(3) {
		t.Error("Delete(3) retornó false")
	}
	if _, ok := tree.Search(3); ok {
		t.Error("3 encontrado tras eliminación")
	}

	if !tree.Delete(5) {
		t.Error("Delete(5) retornó false")
	}
	if _, ok := tree.Search(5); ok {
		t.Error("5 encontrado tras eliminación")
	}

	tree.Search(15)
	if !tree.Delete(15) {
		t.Error("Delete(15) retornó false")
	}
	if _, ok := tree.Search(15); ok {
		t.Error("15 encontrado tras eliminación")
	}

	if tree.Delete(999) {
		t.Error("Delete(999) retornó true, se esperaba false")
	}

	if tree.Size() != len(keys)-3 {
		t.Errorf("size=%d, want %d", tree.Size(), len(keys)-3)
	}
}

func TestDeleteAll(t *testing.T) {
	tree := New[int, string]()
	keys := []int{5, 3, 7, 1, 9, 4, 6}
	for _, k := range keys {
		tree.Insert(k, "")
	}
	for _, k := range keys {
		if !tree.Delete(k) {
			t.Errorf("Delete(%d) retornó false", k)
		}
	}
	if tree.Size() != 0 {
		t.Errorf("size=%d tras eliminar todo, want 0", tree.Size())
	}
	if _, ok := tree.Search(5); ok {
		t.Error("Search en árbol vacío retornó true")
	}
}

func TestInOrder(t *testing.T) {
	tree := New[int, string]()
	keys := []int{10, 3, 15, 7, 1, 20, 5}
	for _, k := range keys {
		tree.Insert(k, "")
	}

	got := tree.InOrder()
	want := make([]int, len(keys))
	copy(want, keys)
	sort.Ints(want)

	if len(got) != len(want) {
		t.Fatalf("InOrder len=%d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("InOrder[%d]=%d, want %d", i, got[i], want[i])
		}
	}

	tree.Delete(3)
	tree.Delete(15)
	got = tree.InOrder()
	want = []int{1, 5, 7, 10, 20}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("InOrder tras delete [%d]=%d, want %d", i, got[i], want[i])
		}
	}
}

func TestLocalityEffect(t *testing.T) {
	const N = 10_000
	tree := New[int, string]()
	for i := 0; i < N; i++ {
		tree.Insert(i, fmt.Sprintf("producto_%d", i))
	}

	hotKey := 4242

	start := time.Now()
	tree.Search(hotKey)
	firstDuration := time.Since(start)

	const repetitions = 1_000
	start = time.Now()
	for i := 0; i < repetitions; i++ {
		tree.Search(hotKey)
	}
	avgRepeated := time.Since(start) / repetitions

	t.Logf("Primer acceso a clave %d:       %v", hotKey, firstDuration)
	t.Logf("Promedio %d accesos repetidos: %v", repetitions, avgRepeated)
	t.Logf("Altura del árbol tras splay:    %d", tree.Height())

	root, _ := tree.Root()
	if root != hotKey {
		t.Errorf("tras accesos repetidos, raíz=%d, want %d", root, hotKey)
	}
}

func TestMultipleHotKeys(t *testing.T) {
	const N = 5_000
	tree := New[string, int]()

	for i := 0; i < N; i++ {
		code := fmt.Sprintf("SKU%05d", i)
		tree.Insert(code, i)
	}

	hotKeys := []string{"SKU00010", "SKU00042", "SKU00100", "SKU00250", "SKU00999"}

	const accesses = 500
	for _, key := range hotKeys {
		for i := 0; i < accesses; i++ {
			if _, ok := tree.Search(key); !ok {
				t.Errorf("hot key %q no encontrada", key)
			}
		}
		root, _ := tree.Root()
		if root != key {
			t.Errorf("tras %d accesos a %q, raíz=%q", accesses, key, root)
		}
	}
	t.Logf("Altura final con %d nodos y accesos sesgados: %d", N, tree.Height())
}

func TestStringKeys(t *testing.T) {
	tree := New[string, string]()
	products := map[string]string{
		"85123A": "WHITE HANGING HEART T-LIGHT HOLDER",
		"22423":  "REGENCY CAKESTAND 3 TIER",
		"85099B": "JUMBO BAG RED RETROSPOT",
		"21212":  "PACK OF 72 RETROSPOT CAKE CASES",
		"20725":  "LUNCH BAG RED RETROSPOT",
	}
	for code, desc := range products {
		tree.Insert(code, desc)
	}
	for code, desc := range products {
		v, ok := tree.Search(code)
		if !ok {
			t.Errorf("producto %q no encontrado", code)
		}
		if v != desc {
			t.Errorf("producto %q: descripción=%q, want %q", code, v, desc)
		}
	}
	if tree.Size() != len(products) {
		t.Errorf("size=%d, want %d", tree.Size(), len(products))
	}
}

func BenchmarkInsert(b *testing.B) {
	tree := New[int, string]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(i, "value")
	}
}

func BenchmarkInsertRandom(b *testing.B) {
	tree := New[int, string]()
	rng := rand.New(rand.NewSource(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Insert(rng.Int(), "value")
	}
}

func BenchmarkSearchRandom(b *testing.B) {
	const N = 10_000
	tree := New[int, string]()
	for i := 0; i < N; i++ {
		tree.Insert(i, "value")
	}
	rng := rand.New(rand.NewSource(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Search(rng.Intn(N))
	}
}

func BenchmarkSearchLocality(b *testing.B) {
	const N = 10_000
	tree := New[int, string]()
	for i := 0; i < N; i++ {
		tree.Insert(i, "value")
	}
	hotKey := N / 2
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Search(hotKey)
	}
}

func BenchmarkSearchSkewed(b *testing.B) {
	const N = 10_000
	tree := New[int, string]()
	for i := 0; i < N; i++ {
		tree.Insert(i, "value")
	}
	// Top 20% de claves
	hotKeys := make([]int, N/5)
	for i := range hotKeys {
		hotKeys[i] = i
	}
	rng := rand.New(rand.NewSource(42))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if rng.Float64() < 0.8 {
			// 80% de accesos van al 20% de claves
			tree.Search(hotKeys[rng.Intn(len(hotKeys))])
		} else {
			tree.Search(rng.Intn(N))
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	tree := New[int, string]()
	for i := 0; i < b.N; i++ {
		tree.Insert(i, "value")
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		tree.Delete(i)
	}
}
