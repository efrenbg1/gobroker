package actions

import (
	"gobroker/db"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	conns  = make(map[string][]*net.Conn)
	lconns sync.RWMutex
)

// watchStart - Add topic to watch
func watchStart(req *sessionData) (bool, string) {
	if req.username != "" {
		topicEnd, or := strconv.Atoi(req.data[4:6])
		if err(or) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = req.data[6:topicEnd]
		if db.InAcls(&req.username, &topic) {
			if req.subscribe != "" {
				watchKill(req)
			}
			lconns.RLock()
			defer lconns.RUnlock()
			req.subscribe = topic
			routines := conns[topic]
			if len(routines) > 4 {
				routines = routines[1:]
				routines = append(routines, req.conn)
			} else {
				routines = append(routines, req.conn)
			}
			conns[topic] = routines
			return true, "MQS4\n"
		}
		return false, "MQS8\n"
	}
	return false, ""
}

// watchKill - Remove current watch
func watchKill(req *sessionData) {
	lconns.RLock()
	defer lconns.RUnlock()
	routines := conns[req.subscribe]
	for i, n := range conns[req.subscribe] {
		if n == req.conn {
			routines[i] = routines[len(routines)-1]
			routines = routines[:len(routines)-1]
			conns[req.subscribe] = routines
			return
		}
	}
}

// watchSend - Send updated message to all the watch
func watchSend(req *sessionData, topic *string, slot *int, payload *string) {
	msg := "MQS5" + strconv.Itoa(*slot) + len2(payload) + *payload + "\n"
	lconns.RLock()
	defer lconns.RUnlock()
	routines := conns[*topic]
	for _, n := range routines {
		if req.conn == n {
			continue
		}
		(*n).SetDeadline(time.Now().Add(time.Duration(req.timeout) * time.Second))
		(*n).Write([]byte(msg))
	}
}
