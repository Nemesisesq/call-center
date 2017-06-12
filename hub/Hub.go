package hub

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Hub struct {
	Agents           map[string]bool
	ReadyAgents      chan Agent
	WaitingProspects chan Prospect
	ActiveCalls      chan Call
}

func NewHub() Hub {
	//make buffering dynamic,
	return Hub{
		Agents:           map[string]bool{},
		ReadyAgents:      make(chan Agent, 100),
		WaitingProspects: make(chan Prospect, 100),
		ActiveCalls:      make(chan Call, 100),
	}
}

type Call struct {
	Agent    Agent
	Prospect Prospect
}

type Person struct {
	Name        string
	PhoneNumber string
}

type Agent struct {
	Person
}

type Prospect struct {
	Person
}

func (h Hub) Status(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Query())
	fmt.Println("Hello", "World")

	b, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(b))
}
