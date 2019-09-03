package main

import (
	"log"
	"net"
	"strconv"
	"time"
)

////////// LOGIN //////////
func login(data *string, conn *net.Conn, timeout *time.Duration, username *string) (bool, string) {
	user_end, err := strconv.Atoi((*data)[4:6])
	if errH(err) {
		return false, ""
	}
	user_end = user_end + 6
	var user = (*data)[6:user_end]
	get := get_pw(&user)
	if errH(err) {
		return false, "MQS9\n"
	}
	if get != "" {
		pw_end, err := strconv.Atoi((*data)[user_end : user_end+2])
		if errH(err) {
			return false, ""
		}
		pw_end = pw_end + user_end + 2
		pw := (*data)[user_end+2 : pw_end]
		if check_pw(&get, &pw) == true {
			timer, err := strconv.Atoi((*data)[pw_end : len(*data)-1])
			if errH(err) {
				return false, "MQS9\n"
			}
			if timer > 99 {
				return false, ""
			}
			*timeout = time.Duration(timer)
			err = (*conn).SetDeadline(time.Now().Add(*timeout * time.Second))
			if errH(err) {
				return false, ""
			}
			*username = user
			log.Println("New connection from " + (*conn).RemoteAddr().String() + " of " + user)
			return true, "MQS0\n"
		} else {
			return false, "MQS9\n"
		}
	}
	return false, ""

}

////////// PUBLISH //////////
func publish(data *string, username *string) (bool, string) {
	if *username != "" {
		topic_end, err := strconv.Atoi((*data)[4:6])
		if errH(err) {
			return false, ""
		}
		topic_end = topic_end + 6
		var topic = (*data)[6:topic_end]
		slot, err := strconv.ParseInt((*data)[topic_end:topic_end+1], 10, 64)
		if errH(err) {
			return false, ""
		}
		if in_acls(username, &topic) {
			payload_end, err := strconv.Atoi((*data)[topic_end+1 : topic_end+3])
			if errH(err) {
				return false, ""
			}
			payload_end = topic_end + payload_end + 3
			payload := (*data)[topic_end+3 : payload_end]
			err = topics.HSet(topic, string(slot), payload).Err()
			if errH(err) {
				return false, ""
			}
			return true, ""
		}
	}
	return false, ""
}

////////// RETRIEVE //////////
func retrieve(data *string, username *string) (bool, string) {
	if *username != "" {
		topic_end, err := strconv.Atoi((*data)[4:6])
		if err != nil {
			return false, ""
		}
		topic_end = topic_end + 6
		var topic = (*data)[6:topic_end]
		slot, err := strconv.ParseInt((*data)[topic_end:topic_end+1], 10, 64)
		if in_acls(username, &topic) {
			payload, err := topics.HGet(topic, string(slot)).Result()
			if err != nil {
				return true, "MQS7\n"
			}
			return true, string("MQS2" + getlen(&payload) + payload + "\n")
		}
		return false, "MQS8\n"
	}
	return false, ""
}

////////// LAST-WILL //////////
func lastpublish(data *string, last_will *string, last_will_s *int64, last_will_p *string, username *string) (bool, string) {
	if *username != "" {
		var topic_end, err = strconv.Atoi((*data)[4:6])
		if err != nil {
			return false, ""
		}
		topic_end = topic_end + 6
		var topic = (*data)[6:topic_end]
		slot, err := strconv.ParseInt((*data)[topic_end:topic_end+1], 10, 64)
		if in_acls(username, &topic) {
			payload_end, err := strconv.Atoi((*data)[topic_end+1 : topic_end+3])
			if err != nil {
				return false, ""
			}
			payload_end = topic_end + payload_end + 3
			payload := (*data)[topic_end+2 : payload_end]
			*last_will = topic
			*last_will_s = slot
			*last_will_p = payload
			return true, "MQS3\n"
		}
		return false, "MQS8\n"
	}
	return false, ""
}
