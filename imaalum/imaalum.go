package imaalum

import (
	"PassGo/dialog"
	_ "PassGo/dialog"
	_ "PassGo/model/user"
	userstruct "PassGo/model/user"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sync"
)

type ImaalumClient struct {
	client *http.Client
}

func Imaalum_login() ImaalumClient {

	file, err := os.ReadFile("indexxx.json")
	if err != nil {
		log.Fatal(err)
	}
	var UserStruct userstruct.User
	json.Unmarshal(file, &UserStruct)
	decodeUser, _ := base64.StdEncoding.DecodeString(UserStruct.User)
	decodePass, _ := base64.StdEncoding.DecodeString(UserStruct.Password)

	formVal := url.Values{
		"username":    {string(decodeUser)},
		"password":    {string(decodePass)},
		"execution":   {"e1s1"},
		"_eventId":    {"submit"},
		"geolocation": {""},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		// error handling
	}

	client := &http.Client{
		Jar: jar,
	}

	urlObj, _ := url.Parse("https://imaluum.iium.edu.my/")
	resp_first, _ := client.Get("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome")
	client.Jar.SetCookies(urlObj, resp_first.Cookies())
	cookies1 := resp_first.Cookies()
	resp, _ := client.PostForm("https://cas.iium.edu.my:8448/cas/login?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome?service=https%3a%2f%2fimaluum.iium.edu.my%2fhome", formVal)
	newCook := append(cookies1, resp.Cookies()...)
	client.Jar.SetCookies(urlObj, newCook)
	// resp_get,_ :=client.Get("https://imaluum.iium.edu.my/MyFinancial")

	//  for _, element := range newCook{
	// 	fmt.Println(element)
	// }

	// if resp_get.StatusCode == 200{
	// 	bodyBytes, err := io.ReadAll(resp_get.Body)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// _ = os.WriteFile("test.html",bodyBytes,0644)
	// // bodyString := string(bodyBytes)
	// // fmt.Println(bodyString)
	// }
	// MessageBoxPlain("Response",string(resp_get.Status))
	// client.Get("https://cas.iium.edu.my:8448/cas/logout?service=http://imaluum.iium.edu.my/")
	var MClient = ImaalumClient{
		client: client,
	}
	return MClient
}
func GetGeneralExamTimeTable(ws *sync.WaitGroup, client *ImaalumClient) {
	response, _ := client.client.Get("https://imaluum.iium.edu.my/MyAcademic/course_timetable")
	if response.StatusCode == 200 {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			dialog.XPlatMessageBox("ERROR", err.Error())
			os.Exit(1)
		}
		_ = os.WriteFile("timetable.pdf", bodyBytes, 0644)

		dialog.XPlatMessageBox("Done", "course_timetable Downloaded")
		ws.Done()

	}

	client.client.Get("https://cas.iium.edu.my:8448/cas/logout?service=http://imaluum.iium.edu.my/")
	ws.Done()
}
func GetConfimationSlip(ws *sync.WaitGroup, client *ImaalumClient) {
	var session_val string = ""
	fmt.Println("Session i.e 2021/2022")
	_, err := fmt.Scan(&session_val)
	if err != nil {
		log.Fatal(err)
	}
	var semester_val string = ""
	fmt.Println("Semester")
	_, err = fmt.Scan(&semester_val)
	if err != nil {
		log.Fatal(err)
	}
	response, _ := client.client.Get(fmt.Sprintf("https://imaluum.iium.edu.my/confirmationslip?ses=%s&sem=%s",session_val,semester_val))
	if response.StatusCode == 200 {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			dialog.XPlatMessageBox("ERROR", err.Error())
			os.Exit(1)
		}

		_ = os.WriteFile("cs.html", bodyBytes, 0644)

		dialog.XPlatMessageBox("Done", "Download Complete")
	}

	client.client.Get("https://cas.iium.edu.my:8448/cas/logout?service=http://imaluum.iium.edu.my/")
	ws.Done()
}
func GetFinance(ws *sync.WaitGroup, client *ImaalumClient) {
	response, _ := client.client.Get("https://imaluum.iium.edu.my/MyFinancial")
	if response.StatusCode == 200 {
		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			dialog.XPlatMessageBox("ERROR", err.Error())
			os.Exit(1)
		}
		_ = os.WriteFile("Finance.pdf", bodyBytes, 0644)

		dialog.XPlatMessageBox("Done", "Finance.PDF Download Complete")
	}

	client.client.Get("https://cas.iium.edu.my:8448/cas/logout?service=http://imaluum.iium.edu.my/")
	ws.Done()
}
