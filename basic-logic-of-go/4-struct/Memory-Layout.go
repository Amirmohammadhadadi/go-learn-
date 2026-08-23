package __struct

import (
	"fmt"
	"reflect"
)

type stats1 struct {
	Reach    uint16
	NumPosts uint8
	NumLikes uint8
}
type stats2 struct {
	NumPosts uint8
	Reach    uint16
	NumLikes uint8
}

func _main() {
	typ1 := reflect.TypeOf(stats1{})
	typ2 := reflect.TypeOf(stats2{})
	fmt.Printf("Struct is %d bytes\n", typ1.Size())
	fmt.Printf("Struct is %d bytes\n", typ2.Size())
}
