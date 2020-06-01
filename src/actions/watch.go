package actions

import (
	"net"
	"strconv"
	"sync"
)

func watch(data *string, lconn *sync.RWMutex, conn *net.Conn, username *string, subscribe *string) (bool, string) {
	if *username != "" {
		topicEnd, err := strconv.Atoi((*data)[4:6])
		if errH(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = (*data)[6:topicEnd]
		if inAcls(username, &topic) {
			lconns.RLock()
			defer lconns.RUnlock()
			*subscribe = topic
			routines := conns[topic]
			routines = append(routines, routine{conn: conn, lconn: lconn})
			conns[topic] = routines
			return true, "MQS4"
		}
		return false, "MQS8\n"
	}
	return false, ""
}
