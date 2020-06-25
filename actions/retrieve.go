package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"strconv"
)

func Retrieve(req *SessionData) (bool, string) {
	if req.Username != "" {
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
		if InAcls(&req.Username, &topic) {
			payload := GetTopic(&topic, &slot)
			if payload == "" {
				return true, "MQS7\n"
			}
			return true, string("MQS2" + Len(&payload) + payload + "\n")
		}
		return false, "MQS8\n"
	}
	return false, ""
}
