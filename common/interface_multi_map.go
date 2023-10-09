package common

type MultiMap[K comparable, V any] interface {
	Put(key K, value V) MultiMap[K, V]
	Get(key K) []V
	Remove(key K) MultiMap[K, V]
	Clean() MultiMap[K, V]
	Keys() []K
	Values() [][]V
	Entries() []MultiMapEntry[K, V]
}

type MultiMapEntry[K comparable, V any] struct {
	Key    K
	Values []V
}
