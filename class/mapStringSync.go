package class

import (
	"database/sql/driver"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/mapkit"
	"sync"
)

/**
针对PG的jsonb
*/

// 同时继承scan和value方法
type MapStringSync struct {
	sync.RWMutex
	Map   map[string]interface{}
	Valid bool
}

// todo 序列化时暂无加锁

func (th *MapStringSync) MarshalJSON() ([]byte, error) {
	if th.Valid {
		return jsonkit.JSON().Marshal(th.Map)
	}
	return []byte("null"), nil
}

func (th *MapStringSync) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		th.Valid = false
		return nil
	}
	var s map[string]interface{}
	if err := jsonkit.JSON().Unmarshal(data, &s); err != nil {
		return err
	}
	th.Valid = true
	th.Map = s
	return nil
}

// Scan implements the Scanner interface.
func (th *MapStringSync) Scan(value interface{}) error {
	if value == nil {
		th.Map, th.Valid = nil, false
		return nil
	}
	th.Valid = true
	th.Map = jsonkit.ParseMap(string(value.([]byte)))
	return nil
}

// Value implements the driver Valuer interface.
func (th *MapStringSync) Value() (driver.Value, error) {
	if !th.Valid || th.Map == nil {
		return nil, nil
	}
	return jsonkit.ToString(th.Map), nil
}

func (th *MapStringSync) IsValid() bool {
	return th.Valid
}

func (th *MapStringSync) Set(val interface{}) {
	th.Lock()
	defer th.Unlock()
	if v, ok := val.(MapStringSync); ok {
		if v.Map == nil {
			th.Map = map[string]interface{}{}
		} else {
			th.Map = v.Map
		}
		th.Valid = true
	} else if v, ok := val.(map[string]interface{}); ok {
		th.Map = v
		th.Valid = true
	} else {
		panic(exception.New("class.MapStringSync set error"))
	}
}

func (th *MapStringSync) PutAll(val map[string]interface{}) {
	th.Lock()
	defer th.Unlock()
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	mapkit.PutAll(th.Map, val)
	th.Valid = true
}

func (th *MapStringSync) PutIfAbsent(key string, val interface{}) {
	th.Lock()
	defer th.Unlock()
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	if _, ok := th.Map[key]; !ok {
		th.Map[key] = val
	}
	th.Valid = true
}

func (th *MapStringSync) Put(key string, val interface{}) {
	th.Lock()
	defer th.Unlock()
	if th.Map == nil {
		th.Map = map[string]interface{}{}
	}
	th.Map[key] = val
	th.Valid = true
}

func (th *MapStringSync) Remove() {
	th.Lock()
	th.Valid = false
	th.Map = map[string]interface{}{}
	th.Unlock()
}

func (th *MapStringSync) Delete(key string) {
	th.Lock()
	th.Valid = false
	delete(th.Map, key)
	th.Unlock()
}

func (th *MapStringSync) IsEmpty() bool {
	if !th.Valid {
		return true
	}
	if len(th.Map) == 0 {
		return true
	}
	return false
}

func (th *MapStringSync) Contains(key string) bool {
	v, ok := th.Map[key]
	if ok {
		return v != nil
	}
	return ok
}
