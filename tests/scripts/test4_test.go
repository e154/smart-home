package scripts

import (
	"testing"
	"github.com/e154/smart-home/common/debug"
	"fmt"
	"encoding/json"
)

type A struct {
	A string `json:"a"`
	C string `json:"c"`
	D string `json:"d"`
}

type B struct {
	B string `json:"b"`
	C string `json:"c"`
}

type D struct {
	*A
	*B
}

type F struct {
	C *D `json:"c"`
}

func (n D) MarshalJSON() (b []byte, err error) {
	switch {
	case n.A != nil:
		b, err = json.Marshal(n.A)
	case n.B != nil:
		b, err = json.Marshal(n.B)
	default:
		err = fmt.Errorf("empty prototype")
		return
	}
	return
}

func Test4(t *testing.T) {

	a := &A{
		A: "A",
		C: "C",
		D: "D",
	}

	//b := &B{
	//	B: "B",
	//	C: "G",
	//}

	c := &D{
		A: a,
		//B: b,
	}

	f := &F{
		C: c,
	}

	debug.Println(c)
	debug.Println(c.A)
	fmt.Println("---")
	fmt.Println(c)
	fmt.Println(c.A)
	fmt.Println("---")
	debug.Println(f)
}
