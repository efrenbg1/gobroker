package main

import (
	"crypto/tls"
	"log"
	"os"
	"syscall"

	"gobroker/actions"
	"gobroker/db"
)

func setRLimit() {
	var rLimit syscall.Rlimit
	or := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err(or) {
		os.Exit(3)
	}
	rLimit.Max = 20000
	rLimit.Cur = 20000
	or = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err(or) {
		os.Exit(3)
	}
	or = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err(or) {
		os.Exit(3)
	}
	log.Println("TCP limit set to:", rLimit.Max)
}

func err(e error) bool {
	if e != nil {
		log.Print(e)
		return true
	}
	return false
}

func main() {
	setRLimit()
	cert, or := tls.LoadX509KeyPair("cert/cert.pem", "cert/key.pem")
	if err(or) {
		os.Exit(1)
	}
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	listen, or := tls.Listen("tcp4", db.Conf.Host, cfg)
	if err(or) {
		os.Exit(1)
	}
	log.Printf("Listening on %s", listen.Addr())
	for {
		conn, or := listen.Accept()
		if or != nil {
			if db.Conf.Debug {
				log.Println(or)
			}
			continue
		}
		go actions.Handle(conn)
	}
}
