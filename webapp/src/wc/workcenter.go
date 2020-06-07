package wc

import "time"

const (
	timeToEscalate  = 10 * time.Minute
	monitorInterval = 5 * time.Second
)

var escalationLevels = [...]string{
	0: "None",
	1: "Immediate Supervisors",
	2: "Managers",
	3: "Directors",
	4: "Executives",
}

var statusLevels = [...]string{
	0: "Green",
	1: "Yellow",
	2: "Red",
}

// Workcenter is a type that holds the information relevant to the Andon status
// of a single workcenter.
type Workcenter struct {
	ID              int
	Name            string
	CurrentProduct  string
	Status          int
	EscalationLevel int
	StatusSetAt     time.Time
}

// TimeAtStatus returns the amount of time that the workcenter has been at
// the current status.
func (wc Workcenter) TimeAtStatus() time.Duration {
	return time.Now().Sub(wc.StatusSetAt)
}

// TimeTillEscalation returns the amount of time before the next escalation level
// is set.
func (wc Workcenter) TimeTillEscalation() time.Duration {
	if wc.Status == 0 {
		return 0
	}
	return timeToEscalate - wc.TimeAtStatus()
}

// StatusDescription returns the description of the work center's status.
func (wc Workcenter) StatusDescription() string {
	return statusLevels[wc.Status]
}

// EscalationLevelDescription returns the description of the work center's esclation level.
func (wc Workcenter) EscalationLevelDescription() string {
	return escalationLevels[wc.EscalationLevel]
}

//Get WorkCenter by ID
func GetWorkCenter(id int) (Workcenter, error) {
	if id > len(workCenters){
		id = 1
	}
	return workCenters[id-1], nil
}

var workCenters = []Workcenter{
	{
		ID:              1,
		Name:            "Assembly Line 1",
		CurrentProduct:  "Widgets",
		Status:          2,
		EscalationLevel: 0,
		StatusSetAt:     time.Now(),
	},
	{
		ID:              2,
		Name:            "Assembly Line 2",
		CurrentProduct:  "Widgets",
		Status:          2,
		EscalationLevel: 0,
		StatusSetAt:     time.Now(),
	},
}
