package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"strconv"
)

// LastPublish data to update the topic if client disconnects
func LastPublish(req *SessionData) (bool, string) {
	if req.Username != "" {
		topicEnd, err := strconv.Atoi((req.Data)[4:6])
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
			payloadEnd, err := strconv.Atoi(req.Data[topicEnd+1 : topicEnd+3])
			if Error(err) {
				return false, ""
			}
			payloadEnd = topicEnd + payloadEnd + 3
			payload := req.Data[topicEnd+3 : payloadEnd]
			req.LwTopic = topic
			req.LwSlot = slot
			req.LwPayload = payload
			return true, "MQS3\n"
		}
		return false, "MQS8\n"
	}
	return false, ""
}
