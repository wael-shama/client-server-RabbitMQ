package orderedmap

import "sync"

type node struct {
	key   string
	value string
	prev  *node
	next  *node
}

type OrderedMap struct {
	mu    sync.RWMutex
	store map[string]*node
	head  *node
	tail  *node
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		store: make(map[string]*node),
	}
}

func (om *OrderedMap) Add(key, value string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	if n, exists := om.store[key]; exists {
		n.value = value
		return
	}

	newNode := &node{key: key, value: value}
	om.store[key] = newNode

	if om.tail == nil {
		om.head = newNode
		om.tail = newNode
		return
	}

	om.tail.next = newNode
	newNode.prev = om.tail
	om.tail = newNode
}

func (om *OrderedMap) Delete(key string) {
	om.mu.Lock()
	defer om.mu.Unlock()

	n, exists := om.store[key]
	if !exists {
		return
	}

	if n.prev != nil {
		n.prev.next = n.next
	} else {
		om.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		om.tail = n.prev
	}

	delete(om.store, key)
}

func (om *OrderedMap) Get(key string) (string, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	n, exists := om.store[key]
	if !exists {
		return "", false
	}
	return n.value, true
}

func (om *OrderedMap) GetAll() []map[string]string {
	om.mu.RLock()
	defer om.mu.RUnlock()

	result := make([]map[string]string, 0, len(om.store))
	current := om.head
	for current != nil {
		result = append(result, map[string]string{
			"key":   current.key,
			"value": current.value,
		})
		current = current.next
	}
	return result
}
