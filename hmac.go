package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"time"
)

/*
 * Base on:
 * http://restcookbook.com/Basics/loggingin/
 * AWS v4 REST http://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html
 * http://www.reddit.com/r/golang/comments/2lu55c/restful_session_token/
 */

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func requestString(method, uri string, body []byte) string {
	k := method + "\n" + url.QueryEscape(uri) + "\n"
	h := hmac.New(sha256.New, []byte(k))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

func signString(method, uri string, body []byte, t time.Time) string {
	rs := requestString(method, uri, body)
	s := "Heketi" + "\n" + t.Format(time.RFC850) + "\n"
	h := hmac.New(sha256.New, []byte(s))
	h.Write([]byte(rs))
	return hex.EncodeToString(h.Sum(nil))
}

func createKey(key string, t time.Time) string {
	h := hmac.New(sha256.New, []byte("Heketi"))
	h.Write([]byte(key))
	h.Write([]byte(timeYYYYMMDD(t)))
	return hex.EncodeToString(h.Sum(nil))
}

func hmacSign(method, uri string, body []byte, t time.Time, key string) string {
	newkey := createKey(key, t)
	message := signString(method, uri, body, t)
	h := hmac.New(sha256.New, []byte(newkey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func timeYYYYMMDD(t time.Time) string {
	return t.Format("<20060102>")
}

func main() {
	t := time.Now()
	fmt.Println(createKey("MyKey", t))
	fmt.Println(requestString("GET", "/nodes", []byte(`{"size":23423}`)))
	fmt.Println(signString("GET", "/nodes", []byte(`{"size":23423}`), t))

	fmt.Println("Final")
	fmt.Printf("Date: %v\n", t.Format(time.RFC850))
	fmt.Println(hmacSign("GET", "/nodes", []byte(`{"size":23423}`), t, "MyKey"))
}
