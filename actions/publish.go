package actions

import (
	"strconv"

	"github.com/efrenbg1/gobroker/db"
)

func publish(req *sessionData) (bool, string) {
	if req.username != "" {
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
		if db.InAcls(&req.username, &topic) {
			payloadEnd, or := strconv.Atoi(req.data[topicEnd+1 : topicEnd+3])
			if err(or) {
				return false, ""
			}
			payloadEnd = topicEnd + payloadEnd + 3
			payload := req.data[topicEnd+3 : payloadEnd]
			db.SetTopic(&topic, &slot, &payload)
			watchSend(req, &topic, &slot, &payload)
			if req.qos {
				return true, "MQS1\n"
			}
			return true, ""
		}
	}
	return false, ""
}
