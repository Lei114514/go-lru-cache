package lru

import (
	"sync"
	"testing"
)

// TestConcurrentBasic 測試基本的並發操作
func TestConcurrentBasic(t *testing.T) {
	c := NewConcurrentCache(100)
	var wg sync.WaitGroup

	// 並發寫入
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put(string(rune(n)), n)
		}(i)
	}

	// 並發讀取
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Get(string(rune(n)))
		}(i)
	}

	wg.Wait()

	// 驗證容量不超過限制
	if c.Len() > c.Cap() {
		t.Errorf("cache length %d exceeds capacity %d", c.Len(), c.Cap())
	}
}

// TestConcurrentReadWrite 測試並發讀寫混合場景
func TestConcurrentReadWrite(t *testing.T) {
	c := NewConcurrentCache(50)
	var wg sync.WaitGroup

	// 預填充一些數據
	for i := 0; i < 50; i++ {
		c.Put(string(rune(i)), i)
	}

	// 並發讀寫
	for i := 0; i < 100; i++ {
		wg.Add(2)

		// 寫goroutine
		go func(n int) {
			defer wg.Done()
			c.Put(string(rune(n)), n*2)
		}(i)

		// 讀goroutine
		go func(n int) {
			defer wg.Done()
			c.Get(string(rune(n % 50)))
		}(i)
	}

	wg.Wait()
}

// TestConcurrentDelete 測試並發刪除操作
func TestConcurrentDelete(t *testing.T) {
	c := NewConcurrentCache(100)
	var wg sync.WaitGroup

	// 預填充
	for i := 0; i < 100; i++ {
		c.Put(string(rune(i)), i)
	}

	// 並發刪除和讀取
	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func(n int) {
			defer wg.Done()
			c.Delete(string(rune(n)))
		}(i)

		go func(n int) {
			defer wg.Done()
			c.Get(string(rune(n)))
		}(i)
	}

	wg.Wait()
}

// TestConcurrentClear 測試並發Clear操作
func TestConcurrentClear(t *testing.T) {
	c := NewConcurrentCache(100)
	var wg sync.WaitGroup

	// 並發寫入和清空
	for i := 0; i < 10; i++ {
		wg.Add(2)

		// 寫入goroutine
		go func(n int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				c.Put(string(rune(n*10+j)), n*10+j)
			}
		}(i)

		// 清空goroutine
		go func() {
			defer wg.Done()
			c.Clear()
		}()
	}

	wg.Wait()

	// 最終狀態應該是一致的（不會panic）
	_ = c.Len()
}

// TestConcurrentEviction 測試並發場景下的LRU淘汰
func TestConcurrentEviction(t *testing.T) {
	c := NewConcurrentCache(10) // 小容量，強制淘汰
	var wg sync.WaitGroup

	// 並發寫入大量數據，觸發頻繁淘汰
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			c.Put(string(rune(n)), n)
		}(i)
	}

	wg.Wait()

	// 驗證容量始終不超過限制
	if c.Len() > c.Cap() {
		t.Errorf("cache length %d exceeds capacity %d", c.Len(), c.Cap())
	}
}

// TestConcurrentCacheImplementsInterface 驗證 ConcurrentCache 實現 Cache 接口
func TestConcurrentCacheImplementsInterface(t *testing.T) {
	var _ Cache = NewConcurrentCache(10)
}

// BenchmarkConcurrentGet 並發讀取基準測試
func BenchmarkConcurrentGet(b *testing.B) {
	c := NewConcurrentCache(1000)
	for i := 0; i < 1000; i++ {
		c.Put(string(rune(i)), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c.Get(string(rune(i % 1000)))
			i++
		}
	})
}

// BenchmarkConcurrentPut 並發寫入基準測試
func BenchmarkConcurrentPut(b *testing.B) {
	c := NewConcurrentCache(1000)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			c.Put(string(rune(i%1000)), i)
			i++
		}
	})
}

// BenchmarkConcurrentMixed 並發混合讀寫基準測試
func BenchmarkConcurrentMixed(b *testing.B) {
	c := NewConcurrentCache(1000)
	for i := 0; i < 1000; i++ {
		c.Put(string(rune(i)), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				c.Get(string(rune(i % 1000)))
			} else {
				c.Put(string(rune(i%1000)), i)
			}
			i++
		}
	})
}
