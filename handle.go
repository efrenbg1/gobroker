package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

func errH(e error) bool {
	if e != nil {
		log.Print(e)
		return true
	} else {
		return false
	}
}

func getlen(str *string) string {
	length := len(*str)
	if length < 10 {
		return "0" + strconv.Itoa(length)
	} else {
		return strconv.Itoa(length)
	}
}

func sina(a *string, list *[]string) bool { //checks if string is in array
	for _, b := range *list {
		if b == *a {
			return true
		}
	}
	return false
}

func handle(data *string, conn *net.Conn, timeout *time.Duration, last_will *string, last_will_s *string, last_will_p *string, username *string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Handle failed:", err)
		}
	}()
	if "MQS" == (*data)[0:3] {
		var act = (*data)[3]
		if act == '0' {
			done, resp := login(data, conn, timeout, username)
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
			done, resp := lastpublish(data, last_will, last_will_s, last_will_p, username)
			if resp != "" {
				_, err := (*conn).Write([]byte(resp))
				if err != nil {
				}
			}
			return done
		}
	}
	return false
}
