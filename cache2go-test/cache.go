package my_cache2go

import (
	"time"
	"sync"
	"fmt"
)

type item struct {
	sync.RWMutex
	key interface{}
	data interface{}
	lifespan time.Duration
	accesson time.Time
}

type cachetable struct {
	sync.RWMutex
	name string
	items map[interface{}]*item
}

var (
	cache = make(map[string]*cachetable)
	mutex  sync.RWMutex
)

func Cache(name string) *cachetable  {
	mutex.Lock()
	defer mutex.Unlock()
	t, ok := cache[name]
	if !ok{
		t = &cachetable{
			name: name,
			items: make(map[interface{}]*item),
		}
		cache[name] = t
	}
	return t
}

func (t *cachetable) Add(key interface{}, lifespan time.Duration, data interface{})  {
	item := &item{
		key: key,
		data: data,
		accesson: time.Now(),
		lifespan: lifespan,
	}
	t.Lock()
	t.Addinternal(item)

}

func (t *cachetable) Addinternal(item *item)  {

	_, ok := t.items[item.key]
	if !ok{
		t.items[item.key] = item
		t.Unlock()
		t.expireCheck()
	}

}

func (t *cachetable) expireCheck()  {
	smallest := 0 * time.Second
	for k, i := range t.items{
		Now := time.Now()
		if Now.Sub(i.accesson) >= i.lifespan{
			fmt.Println("over time, stating to delete")
			t.delete(k)
		}
		if smallest == 0 || i.lifespan - Now.Sub(i.accesson) < smallest{
			smallest = i.lifespan - Now.Sub(i.accesson)
		}
	}

	if smallest > 0{
		time.AfterFunc(smallest, func() {
			go t.expireCheck()
		})
	}
}

func (t *cachetable) delete(key interface{})  {
	t.Lock()
	defer t.Unlock()
	delete(t.items, key)
	fmt.Println("deleting...")
}

func (t *cachetable) Value(key interface{}) interface{} {
	t.Lock()
	defer t.Unlock()
	r, ok := t.items[key]
	if !ok {
		fmt.Println("not in cache!")
		return nil
	}
	return r.data
}
