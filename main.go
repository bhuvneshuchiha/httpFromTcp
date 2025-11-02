package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(ch)
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			i := bytes.IndexByte(data, '\n')
			if i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				ch <- str
				str = ""
			}
			str += string(data)
		}
		if len(str) != 0 {
			ch <- str
		}
	}()
	return ch
}

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("error", "error", err)
	}
	output := getLinesChannel(f)
	for val := range output {
		fmt.Printf("read: %s\n", val)
	}
}
