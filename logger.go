package goodidea

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

var (
	Logr *slog.Logger
)

// ControllerError - An error with nested response
type ControllerError struct {
	//Msg is a string field to capture a custom error message
	Msg string `json:"msg"`
	//Func is the name of the function the error occured in
	Func string `json:"func"`
	//Reason is a longer explication of what happened, could be a message from library
	Reason string `json:"reason"`
	//Cause is a pointer to why this error is being raised, e.g, error from sub functions
	Cause *ControllerError `json:"cause"`
}

func (e *ControllerError) Error() string {
	s, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"msg":"%s","func":"%s", "reason":"%s","aside":"could not marshal this error properly"}`, e.Msg, e.Func, e.Reason)
	}
	return string(s)
}

func SetupLogger() {
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	Logr = l
}
