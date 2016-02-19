package pkg


import (
"sync"
)

type Filter interface {
	Action(key string, value interface{}, format string, meta interface{}) interface{}
}

type FilterFunc func(key string, value interface{}, format string, meta interface{})  interface{}

func (f FilterFunc) Action(key string, value interface{}, format string, meta interface{}) interface{} {
	return f(key, value, format, meta)
}

var (
	filters = []Filter{}
	mutex sync.Mutex
)

func AddFilter(filter FilterFunc) {
	mutex.Lock()
	defer mutex.Unlock()

	filters = append(filters, filter)
}

func init() {

}