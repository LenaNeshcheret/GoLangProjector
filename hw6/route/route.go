package route

import (
	"GoLangProjector/hw6/hw6/publicTransport"
	"fmt"
)

type Route struct {
	Name     string
	Vehicles []publicTransport.PublicTransport
}

func (r *Route) AddVehicleToRoute(vehicle publicTransport.PublicTransport) {
	r.Vehicles = append(r.Vehicles, vehicle)
}

func (r *Route) ShowListOfVehiclesOnRoute() {
	fmt.Printf("All vehicle in our route %s:\n", r.Name)
	for i, vehicle := range r.Vehicles {
		fmt.Printf("%d. %s: %s \n", i+1, vehicle.GetType(), vehicle.GetName())
	}
}
