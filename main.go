package main

/*
TODO: Logout
TODO: Encrypt Password
*/

import (
	"PassGo/dialog"
	imaalum "PassGo/imaalum"
	userstruct "PassGo/model/user"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"syscall"

	"golang.org/x/term"
)

func main() {
	fmt.Println(`
***IIUM WIFI LOGIN***
AUTOMATE LOGIN TO IIUM-STUDENT WITH 2 EASY STEPS
1. STORE YOUR CREDENTIALS
2. CONNECT TO NETWORK

***CREDENTIALS ARE SAVED IN YOUR PC***
	`)
	ws := new(sync.WaitGroup)
	fmt.Println("Select Function\n1.Store Creds\n2.Connect To Network\n3.Logout\n4.Login To iMaalum\n5.Results")
	var user_input int = 0
	_, err := fmt.Scan(&user_input)
	if err != nil {
		log.Fatal(err)
	}
	eval_choice(user_input,ws)

}

func eval_choice(choice int,ws *sync.WaitGroup){
	switch choice {
		case 1:
			store_pass()
		case 2:
			connect_network()
		case 3:
			logout_network()
		case 4:
			dialog.XPlatMessageBox("TEST", "TEST")
			var client = imaalum.Imaalum_login()
			ws.Add(3)
			go imaalum.GetFinance(ws, &client)
			go imaalum.GetConfimationSlip(ws, &client)
			go imaalum.GetGeneralExamTimeTable(ws, &client)
			ws.Wait()
		case 5: 
			checkResult()
	}
}

func checkResult(){
	var UserStruct userstruct.User
	file, err := os.ReadFile("indexxx.json")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	
	var session_val string = ""
	fmt.Println("Session i.e 2021/2022")
	_, err = fmt.Scan(&session_val)
	if err != nil {
		log.Fatal(err)
	}
	var semester_val string = ""
	fmt.Println("Semester")
	_, err = fmt.Scan(&semester_val)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(file, &UserStruct)
	decodeUser, _ := base64.StdEncoding.DecodeString(UserStruct.User)
	decodePass, _ := base64.StdEncoding.DecodeString(UserStruct.Password)
	form_val := url.Values{
		"mat_no" : {string(decodeUser)},
		"pin_no" : {string(decodePass)},
		"sessi"  : {session_val},
		"semester" : {semester_val},
		"login" :{"Login"},
	}
	resp,err := http.PostForm("https://myapps.iium.edu.my/anrapps/viewResult.php", form_val)
	if resp.StatusCode == 200 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			dialog.XPlatMessageBox("ERROR", err.Error())
			os.Exit(1)
		}

		_ = os.WriteFile("result.html", bodyBytes, 0644)

		dialog.XPlatMessageBox("Done", "Download Result Complete")
	}
}

func store_pass() {
	userVal := ""
	passwordVal := ""
	fmt.Println("Insert Username")
	fmt.Scan(&userVal)
	encodeUser := base64.StdEncoding.EncodeToString([]byte(userVal))
	fmt.Println("Insert Password")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	passwordVal = string(bytepw)
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
		panic(err)
	}
	dialog.XPlatMessageBox("SUCCESS", "YOU ARE LOGGED OUT OF IIUM-STUDENT")

}
func connect_network() {
	var UserStruct userstruct.User
	file, err := os.ReadFile("indexxx.json")
	if err != nil {
		log.Fatal(err)
		panic(err)

	}

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
	formVal = nil
	if err != nil {
		dialog.XPlatMessageBox("ERROR", err.Error())
		log.Fatal(err)
		panic(err)
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	dialog.XPlatMessageBox("SUCCESS", "You are now connected to IIUM-Student")
	fmt.Println(res)

}

// MessageBox of Win32 API.
