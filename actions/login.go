package actions

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/efrenbg1/gobroker/db"
)

// login user into system
func login(req *sessionData) (bool, string) {
	userEnd, or := strconv.Atoi(req.data[4:6])
	if err(or) {
		return false, ""
	}
	userEnd = userEnd + 6
	var user = req.data[6:userEnd]
	get := db.GetPw(&user)
	if err(or) {
		return false, "MQS9\n"
	}
	if get != "" {
		pwEnd, or := strconv.Atoi(req.data[userEnd : userEnd+2])
		if err(or) {
			return false, ""
		}
		pwEnd = pwEnd + userEnd + 2
		pw := req.data[userEnd+2 : pwEnd]

		hash := sha256.New()
		_, or = io.Copy(hash, strings.NewReader(get[64:]+pw))
		if err(or) {
			return false, ""
		}
		sum := hash.Sum(nil)
		if get[:64] == hex.EncodeToString(sum) {
			req.username = user
			req.timeout = 100
			req.qos = req.data[pwEnd] == '1'
			log.Println("New connection from " + (*(req.conn)).RemoteAddr().String() + " of " + user)
			return true, "MQS0\n"
		}
		return false, "MQS9\n"
	}
	return false, ""

}
