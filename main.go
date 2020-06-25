package main

import (
	"bufio"
	"crypto/tls"
	. "gobroker/actions"
	. "gobroker/db"
	v "gobroker/windows"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}

func listen(port int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Fatal error:", err)
		}
	}()
	certPem, err := ioutil.ReadFile("cert/cert.pem")
	check(err)
	keyPem, err := ioutil.ReadFile("cert/key.pem")
	check(err)
	cert, err := tls.X509KeyPair(certPem, keyPem)
	check(err)
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	listen, err := tls.Listen("tcp4", ":"+strconv.Itoa(port), cfg)
	defer func() {
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
	log.Printf("Listening on %s", listen.Addr())
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		} else {
			go handler(conn)
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()
	var (
		buf       = make([]byte, 210)
		r         = bufio.NewReader(conn)
		w         = bufio.NewWriter(conn)
		lwTopic   = ""
		lwSlot    = 0
		lwPayload = ""
		username  = ""
		subscribe = ""
	)
	conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
LOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])
		w.Flush()
		switch err {
		case io.EOF:
			break LOOP
		case nil:
			if strings.HasSuffix(data, "\n") {
				conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
				if len(data) > 310 {
					break LOOP
				}
				if handle(&data, &conn, &lwTopic, &lwSlot, &lwPayload, &username, &subscribe) == false {
					break LOOP
				}
			}
		default:
			break LOOP
		}
	}
	if subscribe != "" {
		WatchKill(&conn, &subscribe)
	}
	if lwTopic != "" {
		SetTopic(&lwTopic, &lwSlot, &lwPayload)
	}
	log.Printf("Client from %s disconnected", conn.RemoteAddr())
}

func main() {
	v.TCP()
	listen(2443)
}
