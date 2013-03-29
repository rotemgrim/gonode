// John Granström

// This package defines the basic functionality of gonode.
package gonode

import (
	"encoding/json"
	"fmt"
)

// Signal constants - internal
const signal_NOSIGNAL = -1
const signal_TERMINATION = 1

// Defines a raw command object transmitted over gonode - internal type
type command struct {
	Id     int
	Cmd    CommandData
	Signal int
}

// Defines a raw response object transmitted over gonode - internal type
type response struct {
	Id   int         `json:"id"`   // Make sure response is in lowercase
	Data CommandData `json:"data"` // Make sure response is in lowercase
}

// Processor is the provided function that handles each command and provides a result
type Processor func(cmd CommandData) (resp CommandData)

// CommandData is a wrapping type for representing the JSON objects transmitted
type CommandData map[string]interface{}

// Start the gonode listener which will enter an endless loop while waiting for commands
// Each command will be delegated to new go-routines and processed by the provided Processor.
// This function will return when gonode has been terminated.
func Start(proc Processor) {
	for { // Loop forever
		var s string
		fmt.Scanf("%s", &s) // Receive data from stdin		

		if len(s) < 1 { // Skip empty entries
			continue
		}

		// Parse JSON into Command struct
		var c command
		json.Unmarshal([]byte(s), &c)

		// Handle input
		switch c.Signal {
		case signal_NOSIGNAL:
			go handle(c, proc) // Handle commands on new go-routine
		case signal_TERMINATION:
			return // Abort loop on termination
		}
	}
}

// Handle a command by invoking processor and send result on stdout
func handle(c command, proc Processor) {
	// Create a response with the matching ID
	var r response
	r.Id = c.Id
	r.Data = proc(c.Cmd) // Set response data to processor result
	b, _ := json.Marshal(r)
	fmt.Println(string(b)) // Send JSON result on stdout
}