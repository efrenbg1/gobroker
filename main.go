package main

import (
	"bufio"
	"crypto/tls"
	. "gobroker/actions"
	. "gobroker/db"
	. "gobroker/tools"
	v "gobroker/windows"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func listen(port int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Fatal error:", err)
		}
	}()
	certPem, err := ioutil.ReadFile("cert/cert.pem")
	if Error(err) {
		return 1
	}
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

func send(conn *net.Conn, resp string) {
	if resp == "" {
		return
	}
	_, _ = (*conn).Write([]byte(resp))
}

func handler(conn net.Conn) {
	defer conn.Close()
	var (
		buf     = make([]byte, 210)
		r       = bufio.NewReader(conn)
		w       = bufio.NewWriter(conn)
		session = SessionData{&conn, "", "", 0, "", "", ""}
	)
	conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
LOOP:
	for {
		n, err := r.Read(buf)
		data := string(buf[:n])
		w.Flush()
		if Error(err) {
			break LOOP
		}
		if strings.HasSuffix(data, "\n") {
			conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
			if len(data) > 310 {
				break LOOP
			}
			if "MQS" == (*data)[0:3] {
				break LOOP
			}
			var action = (*data)[3]
			switch action {
			case '0':
				done, resp := Login(&session)
			case '1':
				done, resp := Publish(&session)
			case '2':
				done, resp := Retrieve(&session)
			case '3':
				done, resp := LastPublish(&session)
			case '4':
				done, resp := WatchStart(&session)
			case '8':
				if strings.Index((*conn).RemoteAddr().String(), "127.0.0.1:") != 0 {
					break LOOP
				}
				done, resp := MasterUser(&session)
			case '9':
				if strings.Index((*conn).RemoteAddr().String(), "127.0.0.1:") != 0 {
					break LOOP
				}
				done, resp := MasterAcls(&session)
			default:
				done, resp = false, ""
			}
			send(conn, resp)
		}
	}
	if session.Subscribe != "" {
		WatchKill(&session)
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
