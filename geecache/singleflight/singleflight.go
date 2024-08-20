package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex // protects m
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	if c, ok := g.m[key]; ok {
		c.wg.Wait()         // wait if request processing
		return c.val, c.err // return after request
	}
	c := new(call)
	c.wg.Add(1)  // add lock before request
	g.m[key] = c // add to g.m，which means key has request process

	c.val, c.err = fn() // init request
	c.wg.Done()         // finish request

	delete(g.m, key) // update

	return c.val, c.err
}
