package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"strconv"
)

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
