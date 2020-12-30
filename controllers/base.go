package controllers

import "fmt"

// HandleError = Handle all errors and print out the message
func HandleError(err error) string {
	var handleMsg string
	if err != nil {
		handleMsg = fmt.Sprintf("Error Occured : %v", err)
	}
	return handleMsg
}
