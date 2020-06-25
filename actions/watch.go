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
func WatchStart(data *string, conn *net.Conn, username *string, subscribe *string) (bool, string) {
	if *username != "" {
		topicEnd, err := strconv.Atoi((*data)[4:6])
		if Error(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = (*data)[6:topicEnd]
		if InAcls(username, &topic) {
			if *subscribe != "" {
				WatchKill(conn, subscribe)
			}
			lconns.RLock()
			defer lconns.RUnlock()
			*subscribe = topic
			routines := conns[topic]
			routines = append(routines, conn)
			conns[topic] = routines
			return true, "MQS4"
		}
		return false, "MQS8\n"
	}
	return false, ""
}

// WatchKill - Remove current watch
func WatchKill(conn *net.Conn, topic *string) {
	lconns.RLock()
	defer lconns.RUnlock()
	routines := conns[*topic]
	for i, n := range conns[*topic] {
		if n == conn {
			routines[i] = routines[len(routines)-1]
			routines = routines[:len(routines)-1]
			conns[*topic] = routines
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
