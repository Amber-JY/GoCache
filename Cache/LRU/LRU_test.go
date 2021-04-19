package LRU

import (
	"testing"
)

//使用String作为值类型， 实现Value接口，
type String string

func (s String) Len() int {
	return len(s)
}

func TestLruCache_Add(t *testing.T) {
	cache := New(int64(11))
	cache.Add("testKey", String("1"))
	cache.Add("testKey", String("4567")) //替换

	if _, ok := cache.Get("testKey"); ok {
		//t.Log(cache.schList.Front().Value)
	}

	if cache.curByte != int64(len("testKey")+String("4567").Len()) {
		t.Errorf("byte error. expected %d, but %d", 11, cache.curByte)
	}

}
