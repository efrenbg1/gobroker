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
	"os"
	"strings"
	"time"
)

func handler(conn net.Conn) {
	defer func() {
		recover()
		log.Printf("Client from %s disconnected", conn.RemoteAddr())
	}()
	defer conn.Close()
	var (
		buf = make([]byte, 210)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
		req = SessionData{&conn, "", "", 0, "", "", ""}
	)
	defer func() {
		if req.Subscribe != "" {
			WatchKill(&req)
		}
		if req.LwTopic != "" {
			SetTopic(&req.LwTopic, &req.LwSlot, &req.LwPayload)
		}
	}()
	conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
LOOP:
	for {
		n, err := r.Read(buf)
		req.Data = string(buf[:n])
		w.Flush()
		if err != nil {
			break LOOP
		}
		if strings.HasSuffix(req.Data, "\n") {
			conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
			if len(req.Data) > 310 {
				break LOOP
			}
			if "MQS" != req.Data[0:3] {
				break LOOP
			}
			var action = req.Data[3]
			done, resp := false, ""
			switch action {
			case '0':
				done, resp = Login(&req)
			case '1':
				done, resp = Publish(&req)
			case '2':
				done, resp = Retrieve(&req)
			case '3':
				done, resp = LastPublish(&req)
			case '4':
				done, resp = WatchStart(&req)
			case '6':
				if strings.Index(conn.RemoteAddr().String(), "127.0.0.1:") != 0 {
					break LOOP
				}
				done, resp = MasterPublish(&req)
			case '7':
				if strings.Index(conn.RemoteAddr().String(), "127.0.0.1:") != 0 {
					break LOOP
				}
				done, resp = MasterRetrieve(&req)
			case '8':
				if strings.Index(conn.RemoteAddr().String(), "127.0.0.1:") != 0 {
					break LOOP
				}
				done, resp = MasterUser(&req)
			case '9':
				if strings.Index(conn.RemoteAddr().String(), "127.0.0.1:") != 0 {
					break LOOP
				}
				done, resp = MasterAcls(&req)
			}
			if !done {
				break LOOP
			} else if resp != "" {
				_, _ = conn.Write([]byte(resp))
			}
		}
	}
}

func main() {
	v.TCP()
	certPem, err := ioutil.ReadFile("cert/cert.pem")
	if Error(err) {
		os.Exit(1)
	}
	keyPem, err := ioutil.ReadFile("cert/key.pem")
	if Error(err) {
		os.Exit(1)
	}
	cert, err := tls.X509KeyPair(certPem, keyPem)
	if Error(err) {
		os.Exit(1)
	}
	port := "2443"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	listen, err := tls.Listen("tcp4", ":"+port, cfg)
	if Error(err) {
		os.Exit(1)
	}
	log.Printf("Listening on %s", listen.Addr())
	defer func() {
		log.Println("Error while creating new handler! Recovering...")
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()
	for {
		conn, err := listen.Accept()
		if Error(err) {
			continue
		}
		go handler(conn)
		time.Sleep(time.Millisecond * 50)
	}
}
