package publicTransport

import (
	"fmt"
)

type Airplane struct {
	Type             string
	Name             string
	Capacity         int
	NumberPassengers int
}

func (a *Airplane) AcceptPassenger(passengers int) bool {
	if a.NumberPassengers == a.Capacity {
		fmt.Printf("Passenger cann't travel by this airplane.\n The airplane is completely full.\n")
		return false
	}
	if a.NumberPassengers < a.Capacity {
		a.NumberPassengers += passengers
	}
	return true
}

func (a *Airplane) DropOffPassenger(passengers int) {
	if a.NumberPassengers >= passengers {
		a.NumberPassengers -= passengers
		//fmt.Printf("%d passengers dropped off from the airplane.\n", passengers)
	}
	if a.NumberPassengers <= passengers {
		a.NumberPassengers = 0
		fmt.Printf("There are %d passengers in the airplane. Everyone left.", a.NumberPassengers)
	}
}

func (a *Airplane) GetName() string {
	return a.Name
}

func (a *Airplane) GetType() string {
	return "airplane"
}
