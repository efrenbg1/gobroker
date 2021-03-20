package actions

import (
	"strconv"

	"gobroker/db"
)

func retrieve(req *sessionData) (bool, string) {
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
			payload := db.GetTopic(&topic, &slot)
			if payload == "" {
				return true, "MQS7\n"
			}
			return true, string("MQS2" + len2(&payload) + payload + "\n")
		}
		return false, "MQS8\n"
	}
	return false, ""
}
