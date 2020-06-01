package actions

import (
	"strconv"
)

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
