package main

import (
	"log"
	"net"
	"strconv"
	"sync"
)

////////// LOGIN //////////
func login(data *string, conn *net.Conn, username *string) (bool, string) {
	userEnd, err := strconv.Atoi((*data)[4:6])
	if errH(err) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = (*data)[6:userEnd]
	get := getPw(&user)
	if errH(err) {
		return false, "MQS9\n"
	}
	if get != "" {
		pwEnd, err := strconv.Atoi((*data)[userEnd : userEnd+2])
		if errH(err) {
			return false, ""
		}
		pwEnd = pwEnd + userEnd + 2
		pw := (*data)[userEnd+2 : pwEnd]
		if checkPw(&get, &pw) == true {
			*username = user
			log.Println("New connection from " + (*conn).RemoteAddr().String() + " of " + user)
			return true, "MQS0\n"
		}
		return false, "MQS9\n"
	}
	return false, ""

}

////////// PUBLISH //////////
func publish(data *string, username *string) (bool, string) {
	if *username != "" {
		topicEnd, err := strconv.Atoi((*data)[4:6])
		if errH(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = (*data)[6:topicEnd]
		slot, err := strconv.Atoi((*data)[topicEnd : topicEnd+1])
		if errH(err) || slot > 9 || slot < 0 {
			return false, ""
		}
		if inAcls(username, &topic) {
			payloadEnd, err := strconv.Atoi((*data)[topicEnd+1 : topicEnd+3])
			if errH(err) {

			}
			payloadEnd = topicEnd + payloadEnd + 3
			payload := (*data)[topicEnd+3 : payloadEnd]
			setTopic(&topic, &slot, &payload)
			return true, ""
		}
	}
	return false, ""
}

////////// RETRIEVE //////////
func retrieve(data *string, username *string) (bool, string) {
	if *username != "" {
		topicEnd, err := strconv.Atoi((*data)[4:6])
		if errH(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = (*data)[6:topicEnd]
		slot, err := strconv.Atoi((*data)[topicEnd : topicEnd+1])
		if errH(err) || slot > 9 || slot < 0 {
			return false, ""
		}
		if inAcls(username, &topic) {
			payload := getTopic(&topic, &slot)
			if payload == "" {
				return true, "MQS7\n"
			}
			return true, string("MQS2" + getlen(&payload) + payload + "\n")
		}
		return false, "MQS8\n"
	}
	return false, ""
}

////////// LAST-WILL //////////
func lastpublish(data *string, lastWill *string, lastWillS *int, lastWillP *string, username *string) (bool, string) {
	if *username != "" {
		topicEnd, err := strconv.Atoi((*data)[4:6])
		if errH(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = (*data)[6:topicEnd]
		slot, err := strconv.Atoi((*data)[topicEnd : topicEnd+1])
		if errH(err) || slot > 9 || slot < 0 {
			return false, ""
		}
		if inAcls(username, &topic) {
			payloadEnd, err := strconv.Atoi((*data)[topicEnd+1 : topicEnd+3])
			if errH(err) {
				return false, ""
			}
			payloadEnd = topicEnd + payloadEnd + 3
			payload := (*data)[topicEnd+3 : payloadEnd]
			*lastWill = topic
			*lastWillS = slot
			*lastWillP = payload
			return true, "MQS3\n"
		}
		return false, "MQS8\n"
	}
	return false, ""
}

////////// SUBSCRIBE //////////
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
