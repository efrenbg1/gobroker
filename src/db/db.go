package db

import (
	"database/sql"
	"log"
	"sync"

	_ "../tools/tools.go"
	_ "github.com/go-sql-driver/mysql"
)

var (
	users   = make(map[string]string)
	lusers  sync.RWMutex
	acls    = make(map[string][]string)
	lacls   sync.RWMutex
	topics  = make(map[string][10]string)
	ltopics sync.RWMutex
	db      = dbStart()
)

func GetUser(key *string) string {
	lusers.RLock()
	defer lusers.RUnlock()
	return users[*key]

}

func SetUser(key *string, value *string) {
	lusers.Lock()
	defer lusers.Unlock()
	users[*key] = *value

}

func GetAcls(key *string) []string {
	lacls.RLock()
	defer lacls.RUnlock()
	return acls[*key]
}

func SetAcls(key *string, value *[]string) {
	lacls.Lock()
	defer lacls.Unlock()
	acls[*key] = *value
}

func GetTopic(key *string, slot *int) string {
	ltopics.RLock()
	defer ltopics.RUnlock()
	return topics[*key][*slot]
}

func SetTopic(key *string, slot *int, value *string) {
	ltopics.Lock()
	defer ltopics.Unlock()
	data := topics[*key]
	data[*slot] = *value
	topics[*key] = data
}

func dbStart() *sql.DB {
	//db, err := sql.Open("mysql", "web:SuperPowers4All@tcp(127.0.0.1:3306)/rmote")
	db, err := sql.Open("mysql", "web:Edilizia5!@tcp(192.168.0.4:3306)/rmote")
	if err != nil {
		log.Println("Error connecting to mysql server")
	}
	err = db.Ping()
	if err != nil {
		log.Println("Error pinging mysql server")
	}
	return db
}

func getPw(user *string) string {
	pw := getUser(user)
	if pw != "" {
		return pw
	}
	err := db.QueryRow("SELECT pw FROM user WHERE username=? LIMIT 1", *user).Scan(&pw)
	if err != nil {
		return ""
	}
	setUser(user, &pw)
	return pw
}

func inAcls(user *string, topic *string) bool {
	macs := getAcls(user)
	if SinA(topic, &macs) {
		return true
	}
	macs = nil
	q, err := db.Query("SELECT a.mac FROM acls AS a LEFT JOIN share AS s ON a.mac=s.mac WHERE a.user=(SELECT id FROM user WHERE username=?) OR s.user=(SELECT id FROM user WHERE username=?)", *user, *user)
	if err != nil {
		return false
	}
	var mac string
	for q.Next() {
		err = q.Scan(&mac)
		if err != nil {
			return false
		}
		macs = append(macs, mac)
	}
	if len(macs) > 0 {
		setAcls(user, &macs)
	}
	return sina(topic, &macs)
}
