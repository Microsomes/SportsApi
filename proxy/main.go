package main

import (
	"fmt"
	"io"
	"net"
)

func readLoop(conn net.Conn, forward net.Conn) {
	for {
		buf := make([]byte, 2048)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("end of file reached")
			break
		}
		println(string(buf))
	}
}

func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "whthekiller.com:80")

	if err != nil {
		fmt.Println("error connecting")
	}
	defer dst.Close()

	go func() {
		io.Copy(dst, src)
	}()

	io.Copy(src, dst)

	defer src.Close()

}

func main() {

	//listen tcp
	n1, err := net.Listen("tcp", ":5002")

	if err != nil {
		fmt.Println("error listening")
	}

	for {

		conn2, _ := n1.Accept()

		go handle(conn2)

	}

}
