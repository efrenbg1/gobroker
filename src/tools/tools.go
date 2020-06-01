package tools

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"strconv"
	"strings"
)

// Error handler for debugging
func Error(e error) bool {
	if e != nil {
		log.Print(e)
		return true
	}
	return false
}

// Len returns the length of a string with a padding of 2 characters
func Len(str *string) string {
	length := len(*str)
	if length < 10 {
		return "0" + strconv.Itoa(length)
	}
	return strconv.Itoa(length)
}

// SinA checks if string is in array
func SinA(a *string, list *[]string) bool {
	for _, b := range *list {
		if b == *a {
			return true
		}
	}
	return false
}

// CheckPw verify that string matched with its hash
func CheckPw(hs *string, pw *string) bool {
	hash := sha256.New()
	if _, err := io.Copy(hash, strings.NewReader((*hs)[64:]+*pw)); err != nil {
		return false
	}
	sum := hash.Sum(nil)
	if (*hs)[:64] == hex.EncodeToString(sum) {
		return true
	}
	return false
}
