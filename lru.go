package lru

import(
	"container/list"
)

type entry struct{
	key string
	value any
}

type LRUCache struct{
	capacity int
	cache map[string]*list.Element
	ll *list.List
}

func NewLRUCache(capacity int) *LRUCache{
	if(capacity<=0){
		return nil
	}

	return &LRUCache{
		capacity: capacity,
		cache: make(map[string]*list.Element),
		ll: list.New(),
	}
} 

func (c *LRUCache)Get(key string)(value any, ok bool){
	if elem,hit:=c.cache[key]; hit{
		c.ll.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil,false
}

func (c *LRUCache)Put(key string,value any){
	if elem,hit:=c.cache[key]; hit {
		c.ll.MoveToFront(elem)
		elem.Value.(*entry).value=value
		return 
	}

	elem := c.ll.PushFront(&entry{key:key, value:value})
	c.cache[key]=elem

	if c.ll.Len() > c.capacity {
		c.removeOldest()
	}
}

func (c *LRUCache) Delete(key string){
	if elem, hit := c.cache[key]; hit{
		c.removeElement(elem)
	}
}

func (c *LRUCache) Clear(){
	c.cache= make(map[string]*list.Element)
	c.ll.Init()
}

func (c *LRUCache) Len() int {
	return c.ll.Len()
}

func (c *LRUCache) Cap() int{
	return c.capacity
}

func (c *LRUCache) removeOldest(){
	elem := c.ll.Back()
	if elem!=nil{
		c.removeElement(elem)
	}
}

func (c *LRUCache) removeElement(elem *list.Element){
	c.ll.Remove(elem)
	kv := elem.Value.(*entry)
	delete(c.cache, kv.key)
}