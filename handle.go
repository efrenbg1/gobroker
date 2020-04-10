package main

import (
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

func errH(e error) bool {
	if e != nil {
		log.Print(e)
		return true
	}
	return false
}

func getlen(str *string) string {
	length := len(*str)
	if length < 10 {
		return "0" + strconv.Itoa(length)
	}
	return strconv.Itoa(length)
}

func sina(a *string, list *[]string) bool { //checks if string is in array
	for _, b := range *list {
		if b == *a {
			return true
		}
	}
	return false
}

func handle(data *string, lconn *sync.RWMutex, conn *net.Conn, lastWill *string, lastWillS *int, lastWillP *string, username *string, subscribe *string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Handle failed:", err)
		}
	}()
	if "MQS" == (*data)[0:3] {
		var act = (*data)[3]
		if act == '0' {
			done, resp := login(data, conn, username)
			if resp != "" {
				_, err := (*conn).Write([]byte(resp))
				if err != nil {
				}
			}
			return done
		} else if act == '1' {
			done, resp := publish(data, username)
			if resp != "" {
				_, err := (*conn).Write([]byte(resp))
				if err != nil {
				}
			}
			return done
		} else if act == '2' {
			done, resp := retrieve(data, username)
			if resp != "" {
				_, err := (*conn).Write([]byte(resp))
				if err != nil {
				}
			}
			return done
		} else if act == '3' {
			done, resp := lastpublish(data, lastWill, lastWillS, lastWillP, username)
			if resp != "" {
				_, err := (*conn).Write([]byte(resp))
				if err != nil {
				}
			}
			return done
		} else if act == '4' {
			done, resp := watch(data, lconn, conn, username, subscribe)
			if resp != "" {
				_, err := (*conn).Write([]byte(resp))
				if err != nil {
				}
			}
			return done
		} else if strings.Index((*conn).RemoteAddr().String(), "127.0.0.1:") == 0 {
			userEnd, err := strconv.Atoi((*data)[4:6])
			if !errH(err) {
				userEnd = userEnd + 6
				var user = (*data)[6:userEnd]
				if act == '8' {
					lusers.Lock()
					defer lusers.Unlock()
					delete(users, user)
					_, _ = (*conn).Write([]byte("MQS8\n"))
					//_, _ = (*conns[user][0].conn).Write([]byte("hackeado\n"))
					return true
				} else if act == '9' {
					lacls.Lock()
					defer lacls.Unlock()
					delete(acls, user)
					_, _ = (*conn).Write([]byte("MQS9\n"))
					return true
				}
			}
		}
	}
	return false
}
