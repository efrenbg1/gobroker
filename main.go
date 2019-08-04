package main

import (
	v "./windows"
	"bufio"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func SocketServer(port int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Fatal error:", err)
		}
	}()
	certPem, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		panic(err)
	}
	keyPem, err := ioutil.ReadFile("key.pem")
	if err != nil {
		panic(err)
	}
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	listen, err := tls.Listen("tcp4", ":"+strconv.Itoa(port),cfg)
	defer func(){
		log.Println("Reloading tcp server...")
		err = listen.Close()
		if err != nil {
			log.Println(err)
		}
		recover()
	}()
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		time.Sleep(time.Second)
	}
	log.Printf("Listening on %s",listen.Addr())
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		} else {
			go handler(conn)
		}
		time.Sleep(time.Millisecond*50)
	}

}

func handler(conn net.Conn) {
	defer conn.Close()
	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
		timeout = time.Duration(10)
		last_will = ""
		last_will_p = ""
		username = ""
	)
	conn.SetDeadline(time.Now().Add(timeout*time.Second))
	username = ""
LOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])
		w.Flush()
		switch err {
		case io.EOF:
			break LOOP
		case nil:
			if strings.HasSuffix(data, "\n")  {
				if len(data)>310{
					break LOOP
				}
				if handle(&data, &conn, &timeout, &last_will, &last_will_p, &username) == false{
					break LOOP
				}
				conn.SetDeadline(time.Now().Add(timeout*time.Second))
			}
		default:
			break LOOP
		}

	}
	if last_will != "" {
		err := status.Set(last_will, last_will_p, 0).Err()
		if err != nil {
			log.Println("Error communicating with redis")
			panic(err)
		}
	}
	username = ""
	log.Printf("Client from %s disconnected", conn.RemoteAddr())
}


func main() {
	v.TCP()
	for {
		SocketServer(2443)
	}
}

