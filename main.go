package main

/*
TODO: Logout
TODO: Encrypt Password
*/

import (
	imaalum "PassGo/imaalum"
	_ "PassGo/model/user"
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
	"sync"
	"syscall"
	"unsafe"

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
	fmt.Println("Select Function\n1.Store Creds\n2.Connect To Network\n3.Logout\n4.Login To iMaalum")
	var user_input int = 0
	_, err := fmt.Scan(&user_input)
	if err != nil {
		log.Fatal(err)
	}
	switch user_input {
	case 1:
		store_pass()
	case 2:
		connect_network()
	case 3:
		logout_network()
	case 4:
		var client = imaalum.Imaalum_login()
		ws.Add(3)
		go imaalum.GetFinance(ws, client)
		go imaalum.GetConfimationSlip(ws, client)
		go imaalum.GetGeneralExamTimeTable(ws, client)
		ws.Wait()
	}

}

// func read_pass() {

// }

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
		MessageBoxPlain("ERROR", "UNABLE TO LOGOUT (I GUESS)")
	}
	MessageBoxPlain("SUCCESS", "YOU ARE LOGGED OUT OF IIUM-STUDENT")

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
		MessageBoxPlain("ERROR", err.Error())
		log.Fatal(err)
	}
	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)
	MessageBoxPlain("SUCCESS", "You are now connected to IIUM-Student")
	fmt.Println(res)

}

// MessageBox of Win32 API.
func MessageBox(hwnd uintptr, caption, title string, flags uint) int {
	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		uintptr(flags))

	return int(ret)
}

// MessageBoxPlain of Win32 API.
func MessageBoxPlain(title, caption string) int {
	const (
		NULL  = 0
		MB_OK = 0
	)
	return MessageBox(NULL, caption, title, MB_OK)
}
