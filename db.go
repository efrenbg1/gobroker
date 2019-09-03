package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
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
	topics = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       3,  // use default DB
	})
	db = db_start()
)

//////// STATUS CHECKERS ////////
func db_start() *sql.DB {
	db, err := sql.Open("mysql", "web:SuperPowers4All@tcp(127.0.0.1:3306)/rmote")
	if err != nil {
		log.Println("Error connecting to mysql server")
	}
	err = db.Ping()
	if err != nil {
		log.Println("Error pinging mysql server")
	}
	return db
}

/////////////////////////////////

func get_pw(user *string) string {
	get, err := users.Get(*user).Result()
	if err == nil {
		return get
	} else {
		results, err := db.Query(string("SELECT pw FROM user WHERE username=?"), *user)
		if err != nil {
			return ""
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
	}
	return ""
}

func in_acls(user *string, topic *string) bool {
	if len(*topic) < 100 && len(*user) < 100 {
		var get []string
		get, err := acls.SMembers(*user).Result() //action checker in sina function
		if err != nil {
			return false
		} else if sina(topic, &get) == false {
			results, err := db.Query(string("(SELECT mac FROM acls WHERE user=(SELECT id FROM user WHERE username=%s)) UNION (SELECT acls.mac FROM acls, share WHERE share.user=(SELECT id FROM user WHERE username=%s) AND share.mac=acls.mac)"), *user, *user)
			if err != nil {
				return false
			}
			type Tag struct {
				mac string `json:"mac"`
			}
			get = nil
			for results.Next() {
				var tag Tag
				err = results.Scan(&tag.mac)
				if err != nil {
					return false
				}
				get = append(get, tag.mac)
			}
			if len(get) > 0 {
				err := acls.SAdd(*user, get).Err()
				if err != nil {
					return false
				}
			}
		}
		return sina(topic, &get)
	}
	return false
}

func check_pw(hs *string, pw *string) bool {
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
