package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"log"
	"strconv"
)

// Login user into system
func Login(req *SessionData) (bool, string) {
	userEnd, err := strconv.Atoi(req.Data[4:6])
	if Error(err) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = req.Data[6:userEnd]
	get := GetPw(&user)
	if Error(err) {
		return false, "MQS9\n"
	}
	if get != "" {
		pwEnd, err := strconv.Atoi(req.Data[userEnd : userEnd+2])
		if Error(err) {
			return false, ""
		}
		pwEnd = pwEnd + userEnd + 2
		pw := req.Data[userEnd+2 : pwEnd]
		if CheckPw(&get, &pw) == true {
			req.Username = user
			log.Println("New connection from " + (*(req.Conn)).RemoteAddr().String() + " of " + user)
			return true, "MQS0\n"
		}
		return false, "MQS9\n"
	}
	return false, ""

}
