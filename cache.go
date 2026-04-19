package lru

type Cache interface {
	Get(key string) (value any, ok bool)
	Put(key string, value any)
	Delete(key string)
	Clear()

	Len() int
	Cap() int
}
