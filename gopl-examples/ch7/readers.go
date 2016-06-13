package main

import (
	"fmt"
	"io"
	"os"
)

type Foo struct {
	s   string
	pos int
}

func main() {
	s := "Foo!"
	rdr := NewFoo(s)
	lr := LimitReader(rdr, 30)
	res := make([]byte, 10)
	cnt, err := lr.Read(res)
	fmt.Println("top levelread returned:", cnt, err, string(res))
	if err != nil && err != io.EOF {
		fmt.Println("err:", err)
		os.Exit(1)
	}
	fmt.Println(string(res[0:cnt]))
}

func NewFoo(s string) io.Reader {
	r := &Foo{s: s, pos: 0}
	return r
}

type Bar struct {
	readFrom io.Reader
	left     int64
}

func (r *Bar) Read(p []byte) (n int, err error) {
	fmt.Println("Bar Read")
	var i int
	for r.left > 0 {
		r.left--
		pp := make([]byte, 1)
		cnt, err := r.readFrom.Read(pp)
		if err != nil {
			break
		}
		fmt.Println("org read returned:", cnt, err, pp)
		p[i] = pp[0]
		i++
	}

	return n, nil
}

func LimitReader(r io.Reader, n int64) io.Reader {
	b := Bar{r, n}
	return &b
}

func (r *Foo) Read(p []byte) (n int, err error) {
	fmt.Println("Foo Read")

	n = copy(p, r.s[r.pos:])
	r.pos += n
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil

	// for i, c := range *r {
	// 	fmt.Println("char:", string(c))
	// 	//p = append(p, byte(c))
	// 	p[i] = byte(c)
	// }
	// return len(*r), nil
}
