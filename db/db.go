package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	. "gobroker/tools"

	_ "github.com/go-sql-driver/mysql"
)

// SessionData contains all the data associated with the session of a client
type SessionData struct {
	Conn      *net.Conn
	Data      string
	LwTopic   string
	LwSlot    int
	LwPayload string
	Username  string
	Subscribe string
}

type SettingsDB struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"pw"`
	Database string `json:"db"`
}

type Settings struct {
	Host   string     `json:"host"`
	Master string     `json:"master"`
	Mysql  SettingsDB `json:"mysql"`
}

var (
	Conf    Settings
	users   = make(map[string]string)
	lusers  sync.RWMutex
	acls    = make(map[string][]string)
	lacls   sync.RWMutex
	topics  = make(map[string][10]string)
	ltopics sync.RWMutex
	db      = dbStart()
)

func getUser(key *string) string {
	lusers.RLock()
	defer lusers.RUnlock()
	return users[*key]
}

func setUser(key *string, value *string) {
	lusers.Lock()
	defer lusers.Unlock()
	users[*key] = *value
}

// DelUser - Clear local cache for user password's hash
func DelUser(key *string) {
	lusers.Lock()
	defer lusers.Unlock()
	delete(users, *key)
}

func getAcls(key *string) []string {
	lacls.RLock()
	defer lacls.RUnlock()
	return acls[*key]
}

func setAcls(key *string, value *[]string) {
	lacls.Lock()
	defer lacls.Unlock()
	acls[*key] = *value
}

// DelAcls - Clear local cache for user acls
func DelAcls(key *string) {
	lusers.Lock()
	defer lusers.Unlock()
	delete(acls, *key)
}

// GetTopic - Get value from slot of topic
func GetTopic(key *string, slot *int) string {
	ltopics.RLock()
	defer ltopics.RUnlock()
	return topics[*key][*slot]
}

// SetTopic - Update and slot of a topic with a value
func SetTopic(key *string, slot *int, value *string) {
	ltopics.Lock()
	defer ltopics.Unlock()
	data := topics[*key]
	data[*slot] = *value
	topics[*key] = data
}

func dbStart() *sql.DB {
	configFile, err := os.Open("settings.json")
	if err != nil {
		log.Println("Error while openning config file. Are permissions right?")
		os.Exit(1)
	}
	defer configFile.Close()
	bytes, _ := ioutil.ReadAll(configFile)
	json.Unmarshal(bytes, &Conf)
	log.Println("Config loaded")
	log.Println("Connecting to MySQL server...")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", Conf.Mysql.User, Conf.Mysql.Password, Conf.Mysql.Host, Conf.Mysql.Database))
	if err != nil {
		log.Println(err)
		log.Println("Error connecting to mysql server")
	}
	err = db.Ping()
	if err != nil {
		log.Println(err)
		log.Println("Error pinging mysql server")
	}
	return db
}

// GetPw - Returns hash of password of the passed username. It looks first into local cache and if not found gets it from MySQL server
func GetPw(user *string) string {
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

// InAcls -  Check if user has access to topic name. It looks first into local cache and if not found gets it from MySQL server
func InAcls(user *string, topic *string) bool {
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
	return SinA(topic, &macs)
}
