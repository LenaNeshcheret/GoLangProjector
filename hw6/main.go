package main

import (
	"GoLangProjector/hw6/hw6/passengers"
	"GoLangProjector/hw6/hw6/publicTransport"
	"GoLangProjector/hw6/hw6/route"
)

func main() {

	bus32 := &publicTransport.Bus{Name: "№ 32", Capacity: 30, NumberPassengers: 25}
	bus66 := &publicTransport.Bus{Name: "№ 66", Capacity: 30, NumberPassengers: 5}
	train := &publicTransport.Train{Name: "Kyiv-Lviv", Capacity: 100, NumberPassengers: 99}
	airplane := &publicTransport.Airplane{Name: "Lviv-Dresden", Capacity: 200, NumberPassengers: 200}

	routeA := &route.Route{Name: "A"}
	routeA.AddVehicleToRoute(bus32)
	routeA.AddVehicleToRoute(train)
	routeA.AddVehicleToRoute(bus66)

	routeB := &route.Route{Name: "B"}
	routeB.AddVehicleToRoute(bus66)
	routeB.AddVehicleToRoute(airplane)
	routeB.AddVehicleToRoute(train)

	routeA.ShowListOfVehiclesOnRoute()

	passenger1 := &passengers.Passenger{Name: "John Doe"}
	passenger1.CompleteRoute(routeA)

	routeB.ShowListOfVehiclesOnRoute()

	passenger2 := &passengers.Passenger{Name: "Tom Lee"}
	passenger2.CompleteRoute(routeB)
}
