package publicTransport

import (
	"fmt"
)

type Train struct {
	Type             string
	Name             string
	Capacity         int
	NumberPassengers int
}

func (t *Train) AcceptPassenger(passengers int) bool {
	if t.NumberPassengers == t.Capacity {
		fmt.Printf("Passenger cann't travel by this train.\n The train is completely full.\n")
		return false
	}
	if t.NumberPassengers < t.Capacity {
		t.NumberPassengers += passengers
	}
	return true
}

func (t *Train) DropOffPassenger(passengers int) {
	if t.NumberPassengers >= passengers {
		t.NumberPassengers -= passengers
		//fmt.Printf("%d passengers dropped off from the train.\n", passengers)
	}
	if t.NumberPassengers <= passengers {
		t.NumberPassengers = 0
		fmt.Printf("There are %d passengers in the train. Everyone left.", t.NumberPassengers)
	}
}

func (t *Train) GetName() string {
	return t.Name
}

func (t *Train) GetType() string {
	return "train"
}
