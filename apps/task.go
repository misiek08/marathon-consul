package apps

import (
	"encoding/json"
	"net"
	"fmt"
)

type Task struct {
	ID                 TaskId              `json:"id"`
	TaskStatus         string              `json:"taskStatus"`
	AppID              AppId               `json:"appId"`
	Host               string              `json:"host"`
	Ports              []int               `json:"ports"`
	HealthCheckResults []HealthCheckResult `json:"healthCheckResults"`
}

// Marathon Task Id
// Usually in the form of AppId.uuid with '/' replaced with '_'
type TaskId string

func (id TaskId) String() string {
	return string(id)
}

type HealthCheckResult struct {
	Alive bool `json:"alive"`
}

type TasksResponse struct {
	Tasks []*Task `json:"tasks"`
}

func ParseTasks(jsonBlob []byte) ([]*Task, error) {
	tasks := &TasksResponse{}
	err := json.Unmarshal(jsonBlob, tasks)

	return tasks.Tasks, err
}

func ParseTask(event []byte) (*Task, error) {
	task := &Task{}
	err := json.Unmarshal(event, task)
	return task, err
}

func (t *Task) IsHealthy() bool {
	if len(t.HealthCheckResults) < 1 {
		return false
	}
	register := true
	for _, healthCheckResult := range t.HealthCheckResults {
		register = register && healthCheckResult.Alive
	}
	return register
}

func (t *Task) HostToIPv4() (net.IP, error) {
	IPs, err := net.LookupIP(t.Host)
	if err != nil {
		return nil, err
	}
	for _, IP := range IPs {
		if IP.To4() != nil {
			return IP, nil
		}
	}
	return nil, fmt.Errorf("Could not resolve host to IPv4")
}