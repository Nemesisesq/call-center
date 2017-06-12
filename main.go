package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/urfave/negroni"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()
	n := negroni.Classic()

	r.HandleFunc("/twiml", twiml)
	r.HandleFunc("/call", call)
	//OneOff()

	n.UseHandler(r)
	logrus.Info("Listening on :" + port)
	http.ListenAndServe(":"+port, n)
}

func twiml(w http.ResponseWriter, r *http.Request) {
	//twiml := TwiML{Play: "https://s3.us-east-2.amazonaws.com/sounds4nem/gary_v_rant_60_mins.mp3"}
	tsay := &TwiMLSay{
		Voice: "alice",
		Value: "Please Wait we are connecting you with the prospect",
	}

	tDial := &TwiMLDial{
		Value: "+19145579235",
	}

	twiml := TwiML{
		Say:  *tsay,
		Dial: *tDial,
	}

	x, err := xml.MarshalIndent(twiml, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

type TwiMLSay struct {
	XMLName  xml.Name `xml:"Say"`
	Voice    string   `xml:"voice,attr"`
	Language string   `xml:"language,attr"`
	Value    string   `xml:",chardata"`
}

type TwiMLDial struct {
	XMLName xml.Name `xml:"Dial"`

	Value string `xml:",chardata"`

	Action                        string `xml:"action,attr,omitempty"`                           //relative or absolute URL	no default action for Dial
	Method                        string `xml:"method,attr,omitempty"`                           //GET, POST	POST
	Timeout                       string `xml:"timeout,attr,omitempty"`                          //positive integer	30 seconds
	HangupOnStar                  string `xml:"hangupOnStar,attr,omitempty"`                   //true, false	false
	TimeLimit                     string `xml:"timeLimit,attr,omitempty"`                       //positive integer (seconds)	14400 seconds (4 hours)
	CallerId                      string `xml:"callerId,attr,omitempty"`                        //a valid phone number, or client identifier if you are dialing a <Client>.	Caller's callerId
	Record                        string `xml:"record,attr,omitempty"`                           //do-not-record, record-from-answer, record-from-ringing, record-from-answer-dual, record-from-ringing-dual.For backward compatibility, true is an alias for record-from-answer and false is an alias for do-not-record. do-notrecord
	Trim                          string `xml:"trim,attr,omitempty"`                             //trim-silence, do-not-trim	do-not-trim
	RecordingStatusCallback       string `xml:"recordingStatusCallback,attr,omitempty"`        //relative or absolute URL	none
	RecordingStatusCallbackMethod string `xml:"recordingStatusCallbackMethod,attr,omitempty"` //GET, POST	POST
	RingTone                      string `xml:"ringTone,attr,omitempty"`                        //ISO 3166-1 alpha-2 country code	automatic
}

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Say  TwiMLSay  `xml:",omitempty"`
	Play string `xml:",omitempty"`
	Dial TwiMLDial    `xml:",omitempty"`
}

func call(w http.ResponseWriter, r *http.Request) {
	caller()

	//resp, err := CallAgent("+2165346715")

	//if err != nil {
	//	panic(err)
	//}
	//if resp.StatusCode >= 200 && resp.StatusCode < 300 {
	//	var data map[string]interface{}
	//	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	//	err := json.Unmarshal(bodyBytes, &data)
	//	if err == nil {
	//		fmt.Println(data["sid"])
	//	}
	//} else {
	//	fmt.Println(resp.Status)
	//	w.Write([]byte("Go Royals!"))
	//}

}

func CallAgent(toNum string) (*http.Response, error) {
	accountSid := "AC8babac161b27ec214bed203884635819"
	authToken := "5c575b32cf3208e7a86e849fd0cd697b"
	//callSid := "PNbf2d127871ca9856d3d06e700edbf3a1"
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%v/Calls.json", accountSid)
	v := url.Values{}
	v.Set("To", toNum)
	logrus.Info(toNum)
	v.Set("From", "+19726468378")
	call_in_number := fmt.Sprintf("%v/twiml", os.Getenv("SELF_URL"))
	logrus.Info(call_in_number)
	v.Set("Url", call_in_number)
	rb := *strings.NewReader(v.Encode())
	// Create Client
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest("POST", urlStr, &rb)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// make request
	resp, err := client.Do(req)
	return resp, err
}

func caller() {
	numbers := []string{
		//"+12163466385",
		"+12165346715",
	}
	for _, v := range numbers {
		resp, _ := CallAgent(v)
		logrus.Info(resp)
	}

}

func OneOff() {

	tz, err := time.LoadLocation("America/New_York")

	if err != nil {
		panic(err)
	}

	c := cron.NewWithLocation(tz)

	CallAgent("+12165346715")

	//c.AddFunc("0 0 4 * * 1-5", func() { CallAgent("+12165346715") })
	//c.AddFunc("@every 2h", func() { CallAgent("+12165346715") })
	//c.AddFunc("@every 5s", func() { logrus.Info("making call") })
	//c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
	//c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") })
	c.Start()
	//..
	// Funcs are invoked in their own goroutine, asynchronously.
	//...
	// Funcs may also be added to a running Cron
	//..
	// Inspect the cron job entries' next and previous run times.
	inspect(c.Entries())
	//..
	//c.Stop()  // Stop the scheduler (does not stop any jobs already running).
}
func inspect(entries []*cron.Entry) {
	for _, value := range entries {
		logrus.Info(*value)

	}
}
