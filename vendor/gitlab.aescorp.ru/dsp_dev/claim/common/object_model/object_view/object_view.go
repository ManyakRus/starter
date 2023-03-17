package object_view

import "sort"

// ViewItem -
type ViewItem struct {
	Hide bool
	Key  string
	Val  interface{}
}

type ViewMap []*ViewItem

// Len Получение размера
func (t ViewMap) Len() int {
	return len(t)
}

// Swap Обмен значений
func (t ViewMap) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Less Меньше
func (t ViewMap) Less(i, j int) bool {
	return t[i].Key < t[j].Key
}

// Sort Сортировка
func (t *ViewMap) Sort() {
	sort.Sort(t)
}

// Set Добавление с сортировкой данных
func (t *ViewMap) Set(key string, val interface{}) {
	if t == nil {
		return
	}

	i := t.Search(key)

	// Update
	if i < len(*t) && (*t)[i].Key == key {
		(*t)[i].Val = val
		return
	}

	// Append to end
	if i >= len(*t) {
		*t = append(*t, &ViewItem{
			Key: key,
			Val: val,
		})
		return
	}

	// Insert
	*t = append(*t, &ViewItem{})
	copy((*t)[i+1:], (*t)[i:])
	(*t)[i] = &ViewItem{
		Key: key,
		Val: val,
	}
}

// Append Добавление с сортировкой данных
func (t *ViewMap) Append(key string, val interface{}) {
	if t == nil {
		return
	}

	i := t.Search(key)

	// Update
	if i < len(*t) && (*t)[i].Key == key {
		(*t)[i].Val = val
		return
	}

	// Append to end
	*t = append(*t, &ViewItem{
		Key: key,
		Val: val,
	})
}

// Get Получение элемента по ключу
func (t *ViewMap) Get(key string) (out interface{}) {
	if t == nil {
		return
	}

	i := t.Search(key)

	// Get
	if i < len(*t) && (*t)[i].Key == key {
		out = (*t)[i].Val
		return
	}
	return
}

// GetToIndex Получение элемента по индексу
func (t *ViewMap) GetToIndex(index int) (out interface{}) {
	if t == nil {
		return
	}
	if t.Len() <= index {
		return nil
	}
	return (*t)[index].Val
}

// Unset Очистить значение по ключу
func (t *ViewMap) Unset(key string) {
	if t == nil {
		return
	}

	i := t.Search(key)
	if i < len(*t) && (*t)[i].Key == key {
		copy((*t)[i:], (*t)[i+1:])
		(*t)[len(*t)-1] = nil // or the zero value of T
		*t = (*t)[:len(*t)-1]
	}
}

// Search Поиск значения по ключу
func (t *ViewMap) Search(key string) int {
	if t == nil {
		return 0
	}

	return sort.Search(len(*t), func(i int) bool {
		return (*t)[i].Key >= key
	})
}
