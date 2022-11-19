package lru

import (
	"container/list"
)

// Cache is a LRU cache. concurrent unsafe
type Cache struct {
	maxBytes  int64                         // maximum memory allowed
	nBytes    int64                         // used memory
	ll        *list.List                    // double linkedlist
	cache     map[string]*list.Element      // key: string, value: pointer to the node in list
	OnEvicted func(key string, value Value) // callback function: executed when an entry is deleted
}

// New is Constructure of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get look ups a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToBack(ele)     // Tail stores most recently used
		kv := ele.Value.(*entry) // Value is an empty interface, casting the type to `&entry` here
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.ll.Front() // front stores oldest item
	if ele != nil {
		c.ll.Remove(ele) // remove item from list
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)                                // delete entry
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) //update current memory usge
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// key exists, update the corresponding value
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len()) // update current memory usage
		kv.value = value
	} else {
		ele := c.ll.PushBack(&entry{key, value})         // add entry to list
		c.cache[key] = ele                               // add key and entry to cache
		c.nBytes += int64(len(key)) + int64(value.Len()) // update current memory usage
	}
	// remove oldest item if current memory usage exceeds the maximum memory
	if c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}
