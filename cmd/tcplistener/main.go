package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		str := ""
		for {
			data := make([]byte, 8)

			n, err := f.Read(data)

			if err != nil {
				break
			}

			data = data[:n]

			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				out <- str
				str = ""
			}

			str += string(data)

		}
		if len(str) != 0 {
			out <- str
		}
	}()

	return out

}

func main() {

	port := 42069
	listener, err := net.Listen("tcp", fmt.Sprint(":", port))

	if err != nil {
		log.Fatal("error", "error", err)
	}

	fmt.Printf("Listening on port: %d\n", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println("\nShutdown signal received, shutting down")
		listener.Close()
	}()

	info, err := io.ReadAll(&bytes.Reader{})

	if err != nil {
		log.Fatal("error", "error", err)
	}

	for i := range info {
		fmt.Print(i)
	}

	for {
		conn, err := listener.Accept()

		if err != nil {

			select {
			case <-ctx.Done():
				fmt.Println("Server shutting down gracefully")
				return
			default:
				log.Println("accept error:", err)
				continue
			}

		}

		for line := range getLinesChannel(conn) {

			fmt.Printf("read: %s\n", line)
		}

	}

}
