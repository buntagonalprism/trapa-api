package common

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

type cache struct {
	data map[string]interface{}
}

func NewCache() Cache {
	return &cache{
		data: make(map[string]interface{}),
	}
}

func (c *cache) Set(key string, value interface{}) error {
	c.data[key] = value
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	return c.data[key], nil
}

func (c *cache) Delete(key string) error {
	delete(c.data, key)
	return nil
}
