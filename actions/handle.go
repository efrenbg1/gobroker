package actions

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"gobroker/db"
)

// sessionData contains all the data associated with the session of a client
type sessionData struct {
	conn      *net.Conn
	data      string
	lwTopic   string
	lwSlot    int
	lwPayload string
	username  string
	subscribe string
	timeout   int
	qos       bool
}

// len2 - Returns the length of a string with a padding of 2 characters
func len2(str *string) string {
	length := len(*str)
	if length < 10 {
		return "0" + strconv.Itoa(length)
	}
	return strconv.Itoa(length)
}

// err - Function to handle errors and print them while debugging
func err(e error) bool {
	if e != nil {
		if db.Conf.Debug {
			log.Print(e)
		}
		return true
	}
	return false
}

// Handle - main loop function to handle a single connection
func Handle(conn net.Conn) {
	var (
		buf = make([]byte, 310)
		r   = bufio.NewReaderSize(conn, 310)
		req = sessionData{&conn, "", "", 0, "", "", "", 10, false}
	)

	defer func() {
		recover()
		log.Printf("Client from %s disconnected", conn.RemoteAddr())
	}()
	defer func() {
		conn.Close()
		if req.subscribe != "" {
			watchKill(&req)
		}
		if req.lwTopic != "" {
			db.SetTopic(&req.lwTopic, &req.lwSlot, &req.lwPayload)
			watchSend(&req, &req.lwTopic, &req.lwSlot, &req.lwPayload)
		}
	}()

	conn.SetDeadline(time.Now().Add(time.Duration(req.timeout) * time.Second))

	for {
		len, or := r.Read(buf)
		if or != nil || len > 310 {
			return
		}
		req.data = string(buf[:len])

		if strings.HasSuffix(req.data, "\n") {
			conn.SetDeadline(time.Now().Add(time.Duration(req.timeout) * time.Second))

			if "MQS" != req.data[0:3] {
				return
			}
			var action = req.data[3]
			done, resp := false, ""
			switch action {
			case '0':
				done, resp = login(&req)
			case '1':
				done, resp = publish(&req)
			case '2':
				done, resp = retrieve(&req)
			case '3':
				done, resp = lastPublish(&req)
			case '4':
				done, resp = watchStart(&req)
			case '6':
				if strings.Index(conn.RemoteAddr().String(), db.Conf.Master+":") != 0 {
					return
				}
				done, resp = masterPublish(&req)
			case '7':
				if strings.Index(conn.RemoteAddr().String(), db.Conf.Master+":") != 0 {
					return
				}
				done, resp = masterRetrieve(&req)
			case '8':
				if strings.Index(conn.RemoteAddr().String(), db.Conf.Master+":") != 0 {
					return
				}
				done, resp = masterUser(&req)
			case '9':
				if strings.Index(conn.RemoteAddr().String(), db.Conf.Master+":") != 0 {
					return
				}
				done, resp = masterAcls(&req)
			}
			if !done {
				return
			} else if resp != "" {
				_, _ = conn.Write([]byte(resp))
			}
		}
	}
}
