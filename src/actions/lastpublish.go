package actions

import (
	"strconv"

	t "../tools"
)

// LastPublish data to update the topic if client disconnects
func LastPublish(data *string, lastWill *string, lastWillS *int, lastWillP *string, username *string) (bool, string) {
	if *username != "" {

		topicEnd, err := strconv.Atoi((*data)[4:6])
		if t.Error(err) {
			return false, ""
		}
		topicEnd = topicEnd + 6
		var topic = (*data)[6:topicEnd]
		slot, err := strconv.Atoi((*data)[topicEnd : topicEnd+1])
		if t.Error(err) || slot > 9 || slot < 0 {
			return false, ""
		}
		if inAcls(username, &topic) {
			payloadEnd, err := strconv.Atoi((*data)[topicEnd+1 : topicEnd+3])
			if t.Error(err) {
				return false, ""
			}
			payloadEnd = topicEnd + payloadEnd + 3
			payload := (*data)[topicEnd+3 : payloadEnd]
			*lastWill = topic
			*lastWillS = slot
			*lastWillP = payload
			return true, "MQS3\n"
		}
		return false, "MQS8\n"
	}
	return false, ""
}
