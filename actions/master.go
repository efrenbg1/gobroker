package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"strconv"
)

// MasterPublish allows anyone comming from 127.0.0.1 to publish to a topic
func MasterPublish(req *SessionData) (bool, string) {
	topicEnd, err := strconv.Atoi(req.Data[4:6])
	if Error(err) {
		return false, ""
	}
	topicEnd = topicEnd + 6
	var topic = req.Data[6:topicEnd]
	slot, err := strconv.Atoi(req.Data[topicEnd : topicEnd+1])
	if Error(err) || slot > 9 || slot < 0 {
		return false, ""
	}
	payloadEnd, err := strconv.Atoi(req.Data[topicEnd+1 : topicEnd+3])
	if Error(err) {
		return false, ""
	}
	payloadEnd = topicEnd + payloadEnd + 3
	payload := req.Data[topicEnd+3 : payloadEnd]
	SetTopic(&topic, &slot, &payload)
	WatchSend(&topic, &slot, &payload)
	return true, ""
}

// MasterRetrieve allows anyone comming from 127.0.0.1 to read data from a topic
func MasterRetrieve(req *SessionData) (bool, string) {
	topicEnd, err := strconv.Atoi(req.Data[4:6])
	if Error(err) {
		return false, ""
	}
	topicEnd = topicEnd + 6
	var topic = req.Data[6:topicEnd]
	slot, err := strconv.Atoi(req.Data[topicEnd : topicEnd+1])
	if Error(err) || slot > 9 || slot < 0 {
		return false, ""
	}
	payload := GetTopic(&topic, &slot)
	if payload == "" {
		return true, "MQS7\n"
	}
	return true, string("MQS2" + Len(&payload) + payload + "\n")
}

// MasterUser handle the request from master when some user settings have been changed
func MasterUser(req *SessionData) (bool, string) {
	userEnd, err := strconv.Atoi((req.Data)[4:6])
	if Error(err) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = req.Data[6:userEnd]
	DelUser(&user)
	return true, "MQS8\n"
}

// MasterAcls handle the request from master when acls of a user have been changed
func MasterAcls(req *SessionData) (bool, string) {
	userEnd, err := strconv.Atoi((req.Data)[4:6])
	if Error(err) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = req.Data[6:userEnd]
	DelAcls(&user)
	return true, "MQS9\n"
}
