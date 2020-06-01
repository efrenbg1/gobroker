package actions

import (
	"strconv"
)

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
