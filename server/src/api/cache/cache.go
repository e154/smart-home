package cache

import (
	"github.com/astaxie/beego/cache"
	"fmt"
	"time"
)

type Cache struct {
	bm			cache.Cache
	Cachetime	int64
	Name		string	"maincache"
}

func (c *Cache) init() (*Cache, error) {
	c.log("init")

	var err error
	c.bm, err = cache.NewCache("memory", fmt.Sprintf(`{"interval":%d}`, time.Duration(c.Cachetime) * time.Second))
	if err != nil {
		c.log("error %s", err.Error())
	}
	return c, err
}

func (c *Cache) ClearAll() (*Cache, error) {
	c.log("clear all")

	if c.bm == nil {
		c.init()
	}

	err := c.bm.ClearAll()

	return c, err
}

func (c *Cache) GetKey(key interface {}) string {
	return fmt.Sprintf("%s_%s", c.Name, key.(string))
}

func (c *Cache) Clear(key interface {}) (*Cache, error) {
	cacheKey := c.GetKey(key)
	c.log("clear %s", cacheKey)

	if c.bm == nil {
		c.init()
	}

	err := c.bm.Delete(cacheKey)

	return c, err
}

func (c *Cache) addToGroup(group, key string) (*Cache, error) {

	if c.bm == nil {
		c.init()
	}

	g := []string{}
	w := c.bm.Get(group)
	if w != nil {
		g = w.([]string)
	}

	exist := false
	for _,v := range g {
		if key == v {
			exist = true
		}
	}

	var err error
	if !exist {
		c.log("add to group %s", group)
		g = append(g, key)
		err = c.bm.Put(group, g, time.Duration(c.Cachetime) * time.Second)
	}

	return c, err
}

func (c *Cache) ClearGroup(group string) (*Cache, error) {
	c.log("clear group %s", group)

	if c.bm == nil {
		c.init()
	}

	g := []string{}
	w := c.bm.Get(group)
	if w == nil {
		return c, nil
	}

	g = w.([]string)
	if len(g) == 0 {
		return c, nil
	}

	for _, key := range g {
		c.bm.Delete(key)
	}

	_, err := c.Clear(group)

	return c, err
}

func (c *Cache) Put(group, key string, val interface {}) (*Cache, error) {
	c.log("put key %s", key)

	if c.bm == nil {
		c.init()
	}

	if err := c.bm.Put(key, val, time.Duration(c.Cachetime) * time.Second);err != nil {
		return c, err
	}

	return c.addToGroup(group, key)
}

func (c *Cache) IsExist(key string) bool {
	if c.bm == nil {
		c.init()
	}

	return c.bm.IsExist(key)
}

func (c *Cache) Get(key string) interface {} {
	c.log("get key %s", key)

	if c.bm == nil {
		c.init()
	}

	return c.bm.Get(key)
}

func (c *Cache) log(format string, a ...interface{}) {
	//log.Debug("Cache: ", fmt.Sprintf(format, a...))
}