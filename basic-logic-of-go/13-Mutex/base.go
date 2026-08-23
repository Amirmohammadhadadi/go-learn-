package main

import "sync"

func protected() {
	mux.Lock()
	defer mux.Unlock()
}

type car struct {
	num int
	sync.RWMutex
}

func (c *car) setColorP(num int) {
	c.Lock()
	defer c.Unlock()
	c.num += num
}

func main() {
	c := car{num: 1}
	go setColor(&c, 10)
	go setColor(&c, 5)
}

func setColor(c *car, num int) {
	c.setColorP(num)
}
