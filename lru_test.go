package lru

import (
	"testing"
)

// TestNewLRUCache 測試創建緩存
func TestNewLRUCache(t *testing.T) {
	// 測試正常創建
	c := NewLRUCache(2)
	if c == nil {
		t.Fatal("NewLRUCache(2) returned nil")
	}
	if c.Cap() != 2 {
		t.Errorf("expected capacity 2, got %d", c.Cap())
	}
	if c.Len() != 0 {
		t.Errorf("expected length 0, got %d", c.Len())
	}

	// 測試無效容量
	if NewLRUCache(0) != nil {
		t.Error("NewLRUCache(0) should return nil")
	}
	if NewLRUCache(-1) != nil {
		t.Error("NewLRUCache(-1) should return nil")
	}
}

// TestBasicOperations 測試基本操作
func TestBasicOperations(t *testing.T) {
	c := NewLRUCache(2)

	// 測試 Put 和 Get
	c.Put("a", 1)
	if val, ok := c.Get("a"); !ok || val != 1 {
		t.Errorf("Get(a) = (%v, %v), want (1, true)", val, ok)
	}

	// 測試不存在的 key
	if val, ok := c.Get("b"); ok || val != nil {
		t.Errorf("Get(b) = (%v, %v), want (nil, false)", val, ok)
	}

	// 測試更新值
	c.Put("a", 100)
	if val, ok := c.Get("a"); !ok || val != 100 {
		t.Errorf("Get(a) after update = (%v, %v), want (100, true)", val, ok)
	}
}

// TestLRUEviction 測試 LRU 淘汰策略
func TestLRUEviction(t *testing.T) {
	c := NewLRUCache(2)

	// 填充緩存到容量上限
	c.Put("a", 1)
	c.Put("b", 2)

	if c.Len() != 2 {
		t.Errorf("expected length 2, got %d", c.Len())
	}

	// 添加第三個元素，應該淘汰最久未使用的 'a'
	c.Put("c", 3)

	if c.Len() != 2 {
		t.Errorf("expected length 2 after eviction, got %d", c.Len())
	}

	// 'a' 應該被淘汰了
	if _, ok := c.Get("a"); ok {
		t.Error("key 'a' should have been evicted")
	}

	// 'b' 和 'c' 應該還在
	if val, ok := c.Get("b"); !ok || val != 2 {
		t.Errorf("Get(b) = (%v, %v), want (2, true)", val, ok)
	}
	if val, ok := c.Get("c"); !ok || val != 3 {
		t.Errorf("Get(c) = (%v, %v), want (3, true)", val, ok)
	}
}

// TestLRUOrder 測試 LRU 順序更新
func TestLRUOrder(t *testing.T) {
	c := NewLRUCache(2)

	c.Put("a", 1)
	c.Put("b", 2)

	// 訪問 'a'，使其變成最近使用
	c.Get("a")

	// 添加 'c'，應該淘汰 'b'（因為 'a' 剛被訪問過）
	c.Put("c", 3)

	// 'a' 應該還在
	if val, ok := c.Get("a"); !ok || val != 1 {
		t.Errorf("Get(a) = (%v, %v), want (1, true)", val, ok)
	}

	// 'b' 應該被淘汰
	if _, ok := c.Get("b"); ok {
		t.Error("key 'b' should have been evicted")
	}

	// 'c' 應該還在
	if val, ok := c.Get("c"); !ok || val != 3 {
		t.Errorf("Get(c) = (%v, %v), want (3, true)", val, ok)
	}
}

// TestDelete 測試刪除操作
func TestDelete(t *testing.T) {
	c := NewLRUCache(2)

	c.Put("a", 1)
	c.Put("b", 2)

	// 刪除存在的 key
	c.Delete("a")
	if _, ok := c.Get("a"); ok {
		t.Error("key 'a' should have been deleted")
	}

	// 刪除不存在的 key 不應該出錯
	c.Delete("nonexistent")

	if c.Len() != 1 {
		t.Errorf("expected length 1, got %d", c.Len())
	}
}

// TestClear 測試清空操作
func TestClear(t *testing.T) {
	c := NewLRUCache(2)

	c.Put("a", 1)
	c.Put("b", 2)

	c.Clear()

	if c.Len() != 0 {
		t.Errorf("expected length 0 after Clear, got %d", c.Len())
	}

	if _, ok := c.Get("a"); ok {
		t.Error("key 'a' should have been cleared")
	}
	if _, ok := c.Get("b"); ok {
		t.Error("key 'b' should have been cleared")
	}
}

// TestUpdateMovesToFront 測試更新操作也會將節點移到前面
func TestUpdateMovesToFront(t *testing.T) {
	c := NewLRUCache(2)

	c.Put("a", 1)
	c.Put("b", 2)

	// 更新 'a' 的值，它應該變成最近使用
	c.Put("a", 100)

	// 添加 'c'，應該淘汰 'b'
	c.Put("c", 3)

	if _, ok := c.Get("b"); ok {
		t.Error("key 'b' should have been evicted")
	}
	if val, ok := c.Get("a"); !ok || val != 100 {
		t.Errorf("Get(a) = (%v, %v), want (100, true)", val, ok)
	}
}

// BenchmarkPut 基準測試 Put 操作
func BenchmarkPut(b *testing.B) {
	c := NewLRUCache(1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Put(string(rune(i%1000)), i)
	}
}

// BenchmarkGet 基準測試 Get 操作
func BenchmarkGet(b *testing.B) {
	c := NewLRUCache(1000)
	for i := 0; i < 1000; i++ {
		c.Put(string(rune(i)), i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(string(rune(i % 1000)))
	}
}
