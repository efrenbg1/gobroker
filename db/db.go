package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql" // Import all functions from mysql package
)

// SettingsDB - Contains the DB settings placed inside settings.json
type SettingsDB struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"pw"`
	Database string `json:"db"`
}

// Settings - Contains the general settings placed inside settings.json
type Settings struct {
	Host   string     `json:"host"`
	Master string     `json:"master"`
	Debug  bool       `json:"debug"`
	Mysql  SettingsDB `json:"mysql"`
}

var (
	Conf    Settings // Make global settings accesible to other packages
	users   = make(map[string]string)
	lusers  sync.RWMutex
	acls    = make(map[string][]string)
	lacls   sync.RWMutex
	topics  = make(map[string][10]string)
	ltopics sync.RWMutex
	db      = dbStart()
)

// SinA checks if string is in array
func SinA(a *string, list *[]string) bool {
	for _, b := range *list {
		if b == *a {
			return true
		}
	}
	return false
}

// err - Function to handle errors and print them while debugging
func err(e error) bool {
	if e != nil {
		log.Print(e)
		return true
	}
	return false
}

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
	configFile, or := os.Open("settings.json")
	if err(or) {
		log.Println("Error while openning config file. Are permissions right?")
		os.Exit(1)
	}
	defer configFile.Close()
	bytes, _ := ioutil.ReadAll(configFile)
	or = json.Unmarshal(bytes, &Conf)
	if err(or) {
		log.Println("Error loading config file!")
		os.Exit(1)
	}
	log.Println("Config loaded")
	log.Println("Connecting to MySQL server...")
	db, or := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", Conf.Mysql.User, Conf.Mysql.Password, Conf.Mysql.Host, Conf.Mysql.Database))
	if err(or) {
		log.Println(or)
		log.Println("Error connecting to mysql server")
	}
	or = db.Ping()
	if err(or) {
		log.Println(or)
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
	or := db.QueryRow("SELECT pw FROM user WHERE username=? LIMIT 1", *user).Scan(&pw)
	if err(or) {
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
	q, or := db.Query("SELECT a.mac FROM acls AS a LEFT JOIN share AS s ON a.mac=s.mac WHERE a.user=(SELECT id FROM user WHERE username=?) OR s.user=(SELECT id FROM user WHERE username=?)", *user, *user)
	if err(or) {
		return false
	}
	var mac string
	for q.Next() {
		or = q.Scan(&mac)
		if err(or) {
			return false
		}
		macs = append(macs, mac)
	}
	if len(macs) > 0 {
		setAcls(user, &macs)
	}
	return SinA(topic, &macs)
}
