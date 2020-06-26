package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"net"
	"strconv"
	"sync"
	"time"
)

var (
	conns  = make(map[string][]*net.Conn)
	lconns sync.RWMutex
)

// WatchStart - Add topic to watch
func WatchStart(req *SessionData) (bool, string) {
	if req.Username != "" {
		topicEnd, err := strconv.Atoi(req.Data[4:6])
		if Error(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = req.Data[6:topicEnd]
		if InAcls(&req.Username, &topic) {
			if req.Subscribe != "" {
				WatchKill(req)
			}
			lconns.RLock()
			defer lconns.RUnlock()
			req.Subscribe = topic
			routines := conns[topic]
			if len(routines) > 4 {
				routines = routines[1:]
				routines = append(routines, req.Conn)
			} else {
				routines = append(routines, req.Conn)
			}
			conns[topic] = routines
			return true, "MQS4"
		}
		return false, "MQS8\n"
	}
	return false, ""
}

// WatchKill - Remove current watch
func WatchKill(req *SessionData) {
	lconns.RLock()
	defer lconns.RUnlock()
	routines := conns[req.Subscribe]
	for i, n := range conns[req.Subscribe] {
		if n == req.Conn {
			routines[i] = routines[len(routines)-1]
			routines = routines[:len(routines)-1]
			conns[req.Subscribe] = routines
			return
		}
	}
}

// WatchSend - Send updated message to all the watch
func WatchSend(topic *string, slot *int, payload *string) {
	msg := "MQS5" + Len(topic) + *topic + strconv.Itoa(*slot) + Len(payload) + *payload + "\n"
	lconns.RLock()
	defer lconns.RUnlock()
	routines := conns[*topic]
	for _, n := range routines {
		(*n).SetDeadline(time.Now().Add(time.Duration(10) * time.Second))
		(*n).Write([]byte(msg))
	}
}
