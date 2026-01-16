package main

import (
	"fmt"
	"net"

	"github.com/niudevelop/httpfromtcp/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("Connection has been accepted")
		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		// ch := getLinesChannel(conn)
		// for line := range ch {
		// 	fmt.Println(line)
		// }

		// fmt.Println("Connection has been closed")
	}
}

// func getLinesChannel(r io.ReadCloser) <-chan string {
// 	out := make(chan string)

// 	go func() {
// 		defer close(out)
// 		defer r.Close()

// 		line := ""
// 		b := make([]byte, 8)

// 		for {
// 			n, err := r.Read(b)

// 			if n > 0 {
// 				line += string(b[:n])
// 				parts := strings.Split(line, "\n")

// 				for i := 0; i < len(parts)-1; i++ {
// 					out <- strings.TrimSuffix(parts[i], "\r") // handle \r\n
// 				}
// 				line = parts[len(parts)-1]
// 			}

// 			if err != nil {
// 				break
// 			}
// 		}

// 		if line != "" {
// 			out <- strings.TrimSuffix(line, "\r")
// 		}
// 	}()

// 	return out
// }

// func readFile() {
// 	file, _ := os.Open("messages.txt")
// 	line := ""
// 	b := make([]byte, 8)
// 	for {
// 		n, err := file.Read(b)
// 		if err != nil {
// 			break
// 		}
// 		line += string(b[:n])
// 		parts := strings.Split(line, "\n")
// 		for i := 0; i < len(parts)-1; i++ {
// 			fmt.Printf("read: %s\n", parts[i])
// 		}
// 		line = parts[len(parts)-1]
// 	}
// 	if line != "" {
// 		fmt.Printf("read: %s\n", line)
// 	}
// }
