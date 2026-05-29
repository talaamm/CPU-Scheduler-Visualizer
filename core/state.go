package core

// defines process state

type ProcessState string

const (
	StateNew        ProcessState = "NEW"
	StateReady      ProcessState = "READY"   // to be scheduled to enter cpu
	StateRunning    ProcessState = "RUNNING" // on cpu
	StateWaiting    ProcessState = "WAITING" // for io
	StateTerminated ProcessState = "TERMINATED"
)
