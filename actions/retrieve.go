package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"strconv"
)

func Retrieve(data *string, username *string) (bool, string) {
	if *username != "" {
		topicEnd, err := strconv.Atoi((*data)[4:6])
		if Error(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = (*data)[6:topicEnd]
		slot, err := strconv.Atoi((*data)[topicEnd : topicEnd+1])
		if Error(err) || slot > 9 || slot < 0 {
			return false, ""
		}
		if InAcls(username, &topic) {
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
