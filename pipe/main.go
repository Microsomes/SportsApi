package main

import (
	"io"
	"os"
	"time"
)

func main() {

	pr, pw := io.Pipe()

	go func() {
		pw.Write([]byte("hello"))
		time.Sleep(time.Second)
		pw.Write([]byte("nigger"))

		pw.Close()
	}()

	io.Copy(os.Stdout, pr)

}
