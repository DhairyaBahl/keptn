package v2

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type TaskSequence struct {
	ID       string             `json:"id" bson:"_id"`
	Sequence keptnv2.Sequence   `json:"sequence" bson:"sequence"`
	Status   TaskSequenceStatus `json:"status" bson:"status"`
	Scope    EventScope         `json:"scope" bson:"scope"`
}

type EventScope struct {
	KeptnContext string `json:"keptnContext" bson:"keptnContext"`
	Project      string `json:"project" bson:"project"`
	Stage        string `json:"stage" bson:"stage"`
	Service      string `json:"service" bson:"service"`
}

type TaskSequenceStatus struct {
	State         string                `json:"state" bson:"state"` // triggered, waiting, suspended (approval in progress), paused, finished, cancelled, timedOut
	PreviousTasks []TaskExecutionResult `json:"previousTasks" bson:"previousTasks"`
	CurrentTask   TaskExecution         `json:"currentTask" bson:"currentTask"`
}

type TaskExecutionResult struct {
	Name        string                 `json:"name" bson:"name"`
	TriggeredID string                 `json:"triggeredID" bson:"triggeredID"`
	Result      string                 `json:"result" bson:"result"`
	Status      string                 `json:"status" bson:"status"`
	Properties  map[string]interface{} `json:"properties" bson:"properties"`
}

type TaskExecution struct {
	Name        string      `json:"name" bson:"name"`
	TriggeredID string      `json:"triggeredID" bson:"triggeredID"`
	Events      []TaskEvent `json:"events" bson:"events"`
}

type TaskEvent struct {
	EventType  string                 `json:"eventType" bson:"eventType"`
	Source     string                 `json:"source" bson:"source"`
	Result     string                 `json:"result" bson:"result"`
	Status     string                 `json:"status" bson:"status"`
	Time       string                 `json:"time" bson:"time"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
}
