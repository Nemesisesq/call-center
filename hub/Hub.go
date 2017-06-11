package hub


type Hub struct {

	Agents map[string]bool

	ReadyAgents chan Agent

	WaitingProspects chan Prospect

	ActiveCalls chan Call


}

type Call struct {
	Agent  Agent
	Prospect Prospect
}

type Person struct {
	Name string
	PhoneNumber string
}


type Agent struct {
	Person

}

type Prospect struct {
	Person
}