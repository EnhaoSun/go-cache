package main

import (
	"fmt"
	"gocache/lru"
)

type String string

func (d String) Len() int {
	return len(d)
}

func main() {
	lru := lru.New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		fmt.Println("cache hit key1=1234 failed")
	} else {
		fmt.Printf("key: key1, value: %s\n", v)
	}
	if _, ok := lru.Get("key2"); !ok {
		fmt.Println("cache miss key2 failed")
	}
}
