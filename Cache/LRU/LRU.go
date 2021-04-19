package LRU

import (
	"container/list"
)

const (
	ELEM_ADDED int = iota
	CACHE_EMPTY
	DELETE_OK
	DATA_ERROR
)

//存储结构
type lruCache struct {
	curByte int64
	capByte int64
	schList *list.List               //使用list双向链表
	cache   map[string]*list.Element //哈希表的值是对应的链表节点
}

//定义链表节点的数据结构
type entry struct {
	key   string
	value Value
}

//value可以是任何实现了Len方法的类型
type Value interface {
	Len() int
}

//
// @Description: 工厂函数
// @param cap cache容量
// @return *lruCache
//
func New(capByte int64) *lruCache {
	return &lruCache{
		capByte: capByte,
		schList: new(list.List),
		cache:   make(map[string]*list.Element, capByte),
	}
}

//
// @Description: 新增元素
// @receiver c cache
// @param key
// @param value
// @return 增加元素后是否发生淘汰
//
func (c *lruCache) Add(key string, value Value) int {
	if elem, ok := c.cache[key]; ok { //已存在则更新值
		c.curByte += int64(value.Len() - elem.Value.(*entry).value.Len()) //更新值后占用的空间可能改变
		elem.Value.(*entry).value = value
		c.schList.MoveToFront(elem)
	} else { //不存在加入，判断是否淘汰
		elem := c.schList.PushFront(&entry{key, value})
		c.curByte += elem.Value.(*entry).Size()

		c.cache[key] = elem
	}
	if c.curByte > c.capByte { //超容量需进行淘汰
		c.Delete()
	}
	return ELEM_ADDED
}

//
// @Description: 删除最旧元素
// @receiver c
// @param key
// @return int 是否存在
//
func (c *lruCache) Delete() int {
	if c.curByte == 0 {
		return CACHE_EMPTY
	}
	item := c.schList.Back()
	if _, ok := c.cache[item.Value.(entry).key]; ok {
		c.schList.Remove(item)
		delete(c.cache, item.Value.(*entry).key)
		c.curByte -= item.Value.(*entry).Size()
		return DELETE_OK
	}
	return DATA_ERROR
}

//
// @Description: 查询元素
// @receiver c
// @param key
// @return Value 元素值
// @return bool 是否找到
//
func (c *lruCache) Get(key string) (Value, bool) {
	if elem, ok := c.cache[key]; ok {
		c.schList.MoveToFront(elem)
		entry := elem.Value.(*entry) //获取节点Value（内部实现为空接口），转化为entry
		return entry.value, true
	}
	return nil, false
}

func (e *entry) Size() int64 {
	return int64(len(e.key)) + int64(e.value.Len())
}
