package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		curLine := ""

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)

			if err != nil {
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				curLine += string(data[:i])
				ch <- curLine
				data = data[i+1:]
				curLine = ""
			}
			curLine += string(data)
		}
		if len(curLine) != 0 {
			ch <- curLine
		}
	}()

	return ch
}

func main() {
	l, err := net.Listen("tcp", ":42069")

	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		c, err := l.Accept()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Connection accepted")

		for l := range getLinesChannel(c) {
			fmt.Printf("read: %s\n", l)
		}
	}

}
