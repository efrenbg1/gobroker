package main

import (
	. "gobroker/actions"
	. "gobroker/db"
	. "gobroker/tools"
	"log"
	"net"
	"strconv"
	"strings"
)

func send(conn *net.Conn, resp string) {
	if resp == "" {
		return
	}
	_, _ = (*conn).Write([]byte(resp))
}

func handle(data *string, conn *net.Conn, lastWill *string, lastWillS *int, lastWillP *string, username *string, subscribe *string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Handle failed:", err)
		}
	}()
	if "MQS" == (*data)[0:3] {
		var act = (*data)[3]
		if act == '0' {
			done, resp := Login(data, conn, username)
			send(conn, resp)
			return done
		} else if act == '1' {
			done, resp := Publish(data, username)
			send(conn, resp)
			return done
		} else if act == '2' {
			done, resp := Retrieve(data, username)
			send(conn, resp)
			return done
		} else if act == '3' {
			done, resp := LastPublish(data, lastWill, lastWillS, lastWillP, username)
			send(conn, resp)
			return done
		} else if act == '4' {
			done, resp := WatchStart(data, conn, username, subscribe)
			send(conn, resp)
			return done
		} else if strings.Index((*conn).RemoteAddr().String(), "127.0.0.1:") == 0 {
			userEnd, err := strconv.Atoi((*data)[4:6])
			if !Error(err) {
				userEnd = userEnd + 6
				var user = (*data)[6:userEnd]
				if act == '8' {
					DelUser(&user)
					_, _ = (*conn).Write([]byte("MQS8\n"))
					return true
				} else if act == '9' {
					DelAcls(&user)
					_, _ = (*conn).Write([]byte("MQS9\n"))
					return true
				}
			}
		}
	}
	return false
}
