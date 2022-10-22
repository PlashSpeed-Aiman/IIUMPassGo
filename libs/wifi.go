package libs

import (
	"PassGo/dialog"
	userstruct "PassGo/model/user"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
)

func store_pass(username string, password string) {
	userVal := username
	passwordVal := password
	encodeUser := base64.StdEncoding.EncodeToString([]byte(userVal))
	encodePass := base64.StdEncoding.EncodeToString([]byte(passwordVal))
	data := userstruct.User{
		User:     encodeUser,
		Password: encodePass,
	}
	// fmt.Println(data)
	file, _ := json.MarshalIndent(data, "", " ")
	err2 := os.WriteFile("indexxx.json", file, fs.ModePerm)
	if err2 != nil {
		log.Fatal(err2)
	}
}
func logout_network() {
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	_, err := client.Get("https://captiveportalmahallahgombak.iium.edu.my/auth/logout.html")
	if err != nil {
		dialog.XPlatMessageBox("ERROR", "UNABLE TO LOGOUT (I GUESS)")
	}
	dialog.XPlatMessageBox("SUCCESS", "YOU ARE LOGGED OUT OF IIUM-STUDENT")

}
func connect_network() {
	file, err := os.ReadFile("indexxx.json")
	if err != nil {
		log.Fatal(err)
	}
	var UserStruct userstruct.User
	json.Unmarshal(file, &UserStruct)
	decodeUser, _ := base64.StdEncoding.DecodeString(UserStruct.User)
	decodePass, _ := base64.StdEncoding.DecodeString(UserStruct.Password)
	formVal := url.Values{
		"user":     {string(decodeUser)},
		"password": {string(decodePass)},
		"url":      {"http://www.iium.edu.my/"},
		"cmd":      {"authenticate"},
		"Login":    {"Log In"},
	}
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}
	resp, err := client.PostForm("https://captiveportalmahallahgombak.iium.edu.my/cgi-bin/login", formVal)

	if err != nil {
		dialog.XPlatMessageBox("ERROR", err.Error())
		log.Fatal(err)
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	dialog.XPlatMessageBox("SUCCESS", "You are now connected to IIUM-Student")
	fmt.Println(res)

}
