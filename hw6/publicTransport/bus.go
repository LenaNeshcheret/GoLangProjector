package publicTransport

import (
	"fmt"
)

type Bus struct {
	Type             string
	Name             string
	Capacity         int
	NumberPassengers int
}

func (b *Bus) AcceptPassenger(passengers int) bool {
	if b.NumberPassengers == b.Capacity {
		fmt.Printf("Passenger cann't travel by this bus.\n The bus is completely full.\n")
		return false
	}
	if b.NumberPassengers < b.Capacity {
		b.NumberPassengers += passengers
	}
	return true
}

func (b *Bus) DropOffPassenger(passengers int) {
	if b.NumberPassengers >= passengers {
		b.NumberPassengers -= passengers
	}
	if b.NumberPassengers <= passengers {
		b.NumberPassengers = 0
		fmt.Printf("There are %d passengers in the bus. Everyone left.", b.NumberPassengers)
	}
}

func (b *Bus) GetName() string {
	return b.Name
}

func (b *Bus) GetType() string {
	return "bus"
}
