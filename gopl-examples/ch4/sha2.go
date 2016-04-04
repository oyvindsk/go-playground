package main

import (
	"bufio"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	//shafunc := stdinSha256

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "384":
			sha := stdinSha384()
			fmt.Println("Got SHA:'", sha, "'")
		case "512":
			sha := stdinSha512()
			fmt.Println("Got SHA:'", sha, "'")
		default:
			fmt.Println("huh?")
		}
	} else {
		sha := stdinSha256()
		fmt.Println("Got SHA:'", sha, "'")
	}
}

func readStdio() []byte {
	r := bufio.NewReader(os.Stdin)
	buf, err := ioutil.ReadAll(r)
	fmt.Printf("Read %v bytes, err: %v\n", len(buf), err)
	fmt.Println(buf)
	return buf
}

func stdinSha256() *[32]byte {
	buf := readStdio()
	sha := sha256.Sum256(buf)
	return &sha
}

func stdinSha384() *[48]byte {
	buf := readStdio()
	sha := sha512.Sum384(buf)
	return &sha
}

func stdinSha512() *[64]byte {
	buf := readStdio()
	sha := sha512.Sum512(buf)
	return &sha
}
