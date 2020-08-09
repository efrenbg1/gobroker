package actions

import (
	"gobroker/db"
	"strconv"
)

// masterPublish allows anyone comming from 127.0.0.1 to publish to a topic
func masterPublish(req *sessionData) (bool, string) {
	topicEnd, or := strconv.Atoi(req.data[4:6])
	if err(or) {
		return false, ""
	}
	topicEnd = topicEnd + 6
	var topic = req.data[6:topicEnd]
	slot, or := strconv.Atoi(req.data[topicEnd : topicEnd+1])
	if err(or) || slot > 9 || slot < 0 {
		return false, ""
	}
	payloadEnd, or := strconv.Atoi(req.data[topicEnd+1 : topicEnd+3])
	if err(or) {
		return false, ""
	}
	payloadEnd = topicEnd + payloadEnd + 3
	payload := req.data[topicEnd+3 : payloadEnd]
	db.SetTopic(&topic, &slot, &payload)
	watchSend(req, &topic, &slot, &payload)
	return true, ""
}

// masterRetrieve allows anyone comming from 127.0.0.1 to read data from a topic
func masterRetrieve(req *sessionData) (bool, string) {
	topicEnd, or := strconv.Atoi(req.data[4:6])
	if err(or) {
		return false, ""
	}
	topicEnd = topicEnd + 6
	var topic = req.data[6:topicEnd]
	slot, or := strconv.Atoi(req.data[topicEnd : topicEnd+1])
	if err(or) || slot > 9 || slot < 0 {
		return false, ""
	}
	payload := db.GetTopic(&topic, &slot)
	if payload == "" {
		return true, "MQS7\n"
	}
	return true, string("MQS2" + len2(&payload) + payload + "\n")
}

// masterUser handle the request from master when some user settings have been changed
func masterUser(req *sessionData) (bool, string) {
	userEnd, or := strconv.Atoi((req.data)[4:6])
	if err(or) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = req.data[6:userEnd]
	db.DelUser(&user)
	return true, "MQS8\n"
}

// masterAcls handle the request from master when acls of a user have been changed
func masterAcls(req *sessionData) (bool, string) {
	userEnd, or := strconv.Atoi((req.data)[4:6])
	if err(or) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = req.data[6:userEnd]
	db.DelAcls(&user)
	return true, "MQS9\n"
}
