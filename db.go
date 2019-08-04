package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"regexp"
	"strings"
)

var (
	users = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	acls = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       2,  // use default DB
	})
	status = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       3,  // use default DB
	})
	actions = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       4,  // use default DB
	})
	db = db_start()
)



//////// STATUS CHECKERS ////////
func db_start() *sql.DB {
	db, err := sql.Open("mysql", "web:1Q2w3e4r@tcp(127.0.0.1:3306)/rmote")
	if err != nil {log.Println("Error connecting to mysql server")}
	err = db.Ping()
	if err != nil {log.Println("Error pinging mysql server")}
	return db
}

/////////////////////////////////



func get_pw(user *string) string {
	var safe = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	if safe(*user) == true {
		get, err := users.Get(*user).Result()
		if err == nil {
			return get
		} else {
			results, err := db.Query(string("SELECT pw FROM user WHERE username='" + *user + "' limit 1"))
			if err != nil {
				db_start()
				results, err = db.Query(string("SELECT pw FROM user WHERE username='" + *user + "' limit 1"))
				if err != nil {
					return ""
				}
			}
			type Tag struct {
				pw string `json:"pw"`
			}
			for results.Next() {
				var tag Tag
				err = results.Scan(&tag.pw)
				if err != nil {
					return ""
				}
				err = users.Set(*user, tag.pw, 0).Err()
				if err != nil {
					return ""
				}
				return tag.pw
			}
			return ""
		}
	}
	return ""
}




func in_acls(user *string, topic *string) bool {
	var safe = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString
	var topic_safe = regexp.MustCompile(`^([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})$`).MatchString
	if safe(*user) == true && topic_safe(*topic) == true {
		var get []string
		get, err := acls.SMembers(*user).Result() //action checker in sina function
		if err != nil {
			return false
		} else if sina(topic, &get) == false {
			results, err := db.Query(string("SELECT mac FROM acls WHERE username='" + *user + "'"))
			if err != nil {
				db_start()
				results, err = db.Query(string("SELECT mac FROM acls WHERE username='" + *user + "'"))
				if err != nil {
					return false
				}
			}
			type Tag struct {
				mac string `json:"mac"`
			}
			for results.Next() {
				var tag Tag
				err = results.Scan(&tag.mac)
				if err != nil {
					return false
				}
				get = append(get, tag.mac)
			}
			if len(get) > 0 {
				err := acls.SAdd(*user,get).Err()
				if err != nil {return false}
			}
		}
		return sina(topic, &get)
	}
	return false
}



func check_pw (hs *string, pw *string) bool{
	hash := sha256.New()
	if _, err := io.Copy(hash, strings.NewReader(*pw)); err != nil {
		return false
	}
	sum := hash.Sum(nil)
	if *hs == hex.EncodeToString(sum) {
		return true
	} else {
		return false
	}
}