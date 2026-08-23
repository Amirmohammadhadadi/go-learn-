package main

import "fmt"

type car struct {
	color string
}

func (c car) setColor(color string) {
	c.color = color
}
func (c *car) setColorP(color string) {
	c.color = color
}

func Rmain1() {
	c := car{color: "white"}
	c.setColor("blue")
	fmt.Println(c.color)
	// prints "white"
}
func Rmain2() {
	c := car{color: "white"}
	c.setColorP("blue")
	fmt.Println(c.color)
	// prints "blue"
}
