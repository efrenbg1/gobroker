package actions

import (
	. "gobroker/db"
	. "gobroker/tools"
	"log"
	"net"
	"strconv"
)

func Login(data *string, conn *net.Conn, username *string) (bool, string) {
	userEnd, err := strconv.Atoi((*data)[4:6])
	if Error(err) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = (*data)[6:userEnd]
	get := GetPw(&user)
	if Error(err) {
		return false, "MQS9\n"
	}
	if get != "" {
		pwEnd, err := strconv.Atoi((*data)[userEnd : userEnd+2])
		if Error(err) {
			return false, ""
		}
		pwEnd = pwEnd + userEnd + 2
		pw := (*data)[userEnd+2 : pwEnd]
		if CheckPw(&get, &pw) == true {
			*username = user
			log.Println("New connection from " + (*conn).RemoteAddr().String() + " of " + user)
			return true, "MQS0\n"
		}
		return false, "MQS9\n"
	}
	return false, ""

}
