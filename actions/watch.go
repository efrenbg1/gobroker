package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"net"
	"strconv"
	"sync"
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
			routines = append(routines, req.Conn)
			conns[topic] = routines
			return true, "MQS4"
		}
		return false, "MQS8\n"
	}
	return false, ""
}

// WatchKill - Remove current watch
func WatchKill(s *SessionData) {
	lconns.RLock()
	defer lconns.RUnlock()
	routines := conns[s.Subscribe]
	for i, n := range conns[s.Subscribe] {
		if n == s.Conn {
			routines[i] = routines[len(routines)-1]
			routines = routines[:len(routines)-1]
			conns[s.Subscribe] = routines
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
		(*n).Write([]byte(msg))
	}
}
