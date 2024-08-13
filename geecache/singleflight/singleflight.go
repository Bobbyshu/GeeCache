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
