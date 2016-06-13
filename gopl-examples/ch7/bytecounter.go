package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {

	ww, cnt := CountingWriter(os.Stdout)
	ww.Write([]byte("hello!"))
	ww.Write([]byte("1234"))

	fmt.Println("\tcnt:", *cnt)

	//var c ByteCounter
	//c.Write([]byte("hello"))
	//fmt.Println(c) // "5", = len("hello")

	//var w WordCounter
	//w.Write([]byte("en to tre"))
	//fmt.Println(w)
}

type CWriter struct {
	io.Writer
	cnt  *int64
	orgW io.Writer
}

func (c CWriter) Write(p []byte) (int, error) {
	*c.cnt += int64(len(p))
	return c.orgW.Write(p)

}

func CountingWriter(w io.Writer) (io.Writer, *int64) {

	var cw CWriter
	cw.orgW = w
	count := int64(0)
	cw.cnt = &count
	return cw, cw.cnt

}

type WordCounter int

func (c *WordCounter) Write(p []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		*c++
	}
	return len(p), nil
}

type ByteCounter int

func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p)) // convert int to ByteCounter
	return len(p), nil
}
