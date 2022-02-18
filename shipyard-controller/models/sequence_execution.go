package models

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

type SequenceExecution struct {
	ID       string                  `json:"_id" bson:"_id"`
	Sequence keptnv2.Sequence        `json:"sequence" bson:"sequence"`
	Status   SequenceExecutionStatus `json:"status" bson:"status"`
	Scope    EventScope              `json:"scope" bson:"scope"`
	// InputProperties contains properties of the event which triggered the task sequence
	InputProperties interface{} `json:"inputProperties" bson:"inputProperties"`
}

type SequenceExecutionStatus struct {
	State string `json:"state" bson:"state"` // triggered, waiting, suspended (approval in progress), paused, finished, cancelled, timedOut
	// StateBeforePause is needed to keep track of the state before a sequence has been paused. Example: when a sequence has been paused while being queued, and then resumed, it should not be set to started immediately, but to the state it had before
	StateBeforePause string                `json:"stateBeforePause" bson:"stateBeforePause"`
	PreviousTasks    []TaskExecutionResult `json:"previousTasks" bson:"previousTasks"`
	CurrentTask      TaskExecutionState    `json:"currentTask" bson:"currentTask"`
}

type TaskExecutionResult struct {
	Name        string                 `json:"name" bson:"name"`
	TriggeredID string                 `json:"triggeredID" bson:"triggeredID"`
	Result      string                 `json:"result" bson:"result"`
	Status      string                 `json:"status" bson:"status"`
	Properties  map[string]interface{} `json:"properties" bson:"properties"`
}

type TaskExecutionState struct {
	Name        string      `json:"name" bson:"name"`
	TriggeredID string      `json:"triggeredID" bson:"triggeredID"`
	Events      []TaskEvent `json:"events" bson:"events"`
}

func (e *SequenceExecution) GetNextTaskOfSequence() *keptnv2.Task {
	if e.Status.CurrentTask.IsFailed() {
		return nil
	}
	nextTaskIndex := 0
	if e.Status.PreviousTasks != nil && len(e.Status.PreviousTasks) > 0 {
		nextTaskIndex = len(e.Status.PreviousTasks)
	}

	if len(e.Sequence.Tasks) > nextTaskIndex {
		return &e.Sequence.Tasks[nextTaskIndex]
	}
	return nil
}

func (e *SequenceExecution) GetLastTaskExecutionResult() TaskExecutionResult {
	if len(e.Status.PreviousTasks) == 0 {
		return TaskExecutionResult{}
	}
	return e.Status.PreviousTasks[len(e.Status.PreviousTasks)-1]
}

func (e *SequenceExecution) CompleteCurrentTask() {
	var result keptnv2.ResultType
	var status keptnv2.StatusType
	if e.Status.CurrentTask.IsFailed() {
		result = keptnv2.ResultFailed
	} else if e.Status.CurrentTask.IsWarning() {
		result = keptnv2.ResultWarning
	} else {
		result = keptnv2.ResultPass
	}
	if e.Status.CurrentTask.IsErrored() {
		status = keptnv2.StatusErrored
	} else {
		status = keptnv2.StatusSucceeded
	}

	e.Status.PreviousTasks = append(
		e.Status.PreviousTasks,
		TaskExecutionResult{
			Name:        e.Status.CurrentTask.Name,
			TriggeredID: e.Status.CurrentTask.TriggeredID,
			Result:      string(result),
			Status:      string(status),
			Properties:  nil,
		},
	)
	e.Status.CurrentTask = TaskExecutionState{}
}

func (e *SequenceExecution) GetNextTriggeredEventData() map[string]interface{} {
	eventPayload := map[string]interface{}{}

	if inputMap, ok := e.InputProperties.(map[string]interface{}); ok {
		eventPayload = inputMap
	}

	eventPayload["project"] = e.Scope.Project
	eventPayload["stage"] = e.Scope.Stage
	eventPayload["service"] = e.Scope.Service

	nextTask := e.GetNextTaskOfSequence()
	if nextTask != nil && nextTask.Properties != nil {
		eventPayload[nextTask.Name] = nextTask.Properties
	}

	if len(e.Status.PreviousTasks) > 0 {
		for _, previousTask := range e.Status.PreviousTasks {
			eventPayload[previousTask.Name] = previousTask.Properties
		}
		lastTaskIndex := len(e.Status.PreviousTasks) - 1
		eventPayload["result"] = e.Status.PreviousTasks[lastTaskIndex].Result
		eventPayload["status"] = e.Status.PreviousTasks[lastTaskIndex].Status
	}

	return eventPayload
}

func (e *SequenceExecution) IsPaused() bool {
	return e.Status.State == SequencePaused
}

// CanBePaused determines whether a sequence can be paused, based on its current state. E.g. a finished sequence cannot be paused
func (e *SequenceExecution) CanBePaused() bool {
	return e.Status.State == SequenceStartedState || e.Status.State == SequenceWaitingState || e.Status.State == SequenceTriggeredState
}

// Pause tries to pause the sequence execution, based on its current state. If it was successful, returns true. If it could not be paused, false is returned
func (e *SequenceExecution) Pause() bool {
	if !e.CanBePaused() {
		return false
	}
	e.Status.StateBeforePause = e.Status.State
	return true
}

// Resume tries to resume the sequence execution, based on its current state. If it was successful, returns true. If it could not be paused, false is returned
func (e *SequenceExecution) Resume() bool {
	if !e.IsPaused() {
		return false
	}
	e.Status.State = e.Status.StateBeforePause
	return true
}

// IsFinished indicates if a task is finished, i.e. the number of task.started and task.finished events line up
func (e *TaskExecutionState) IsFinished() bool {
	if len(e.Events) == 0 {
		return false
	}
	nrStartedEvents := 0
	nrFinishedEvents := 0
	for _, event := range e.Events {
		if keptnv2.IsStartedEventType(event.EventType) {
			nrStartedEvents++
		} else if keptnv2.IsFinishedEventType(event.EventType) {
			nrFinishedEvents++
		}
	}

	if nrFinishedEvents == nrStartedEvents && nrFinishedEvents > 0 {
		return true
	}
	return false
}

func (e *TaskExecutionState) IsFailed() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Result == string(keptnv2.ResultFailed) {
				return true
			}
		}
	}
	return false
}

func (e *TaskExecutionState) IsWarning() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Result == string(keptnv2.ResultWarning) {
				return true
			}
		}
	}
	return false
}

func (e *TaskExecutionState) IsPassed() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Result == string(keptnv2.ResultFailed) || event.Result == string(keptnv2.ResultWarning) {
				return false
			}
		}
	}
	return true
}

func (e *TaskExecutionState) IsErrored() bool {
	for _, event := range e.Events {
		if keptnv2.IsFinishedEventType(event.EventType) {
			if event.Status == string(keptnv2.StatusErrored) {
				return true
			}
		}
	}
	return false
}

type TaskEvent struct {
	EventType  string                 `json:"eventType" bson:"eventType"`
	Source     string                 `json:"source" bson:"source"`
	Result     string                 `json:"result" bson:"result"`
	Status     string                 `json:"status" bson:"status"`
	Time       string                 `json:"time" bson:"time"`
	Properties map[string]interface{} `json:"properties" bson:"properties"`
}

type SequenceExecutionFilter struct {
	Scope              EventScope
	Status             []string
	Name               string
	CurrentTriggeredID string
}

type SequenceExecutionUpsertOptions struct {
	CheckUniqueTriggeredID bool
}
