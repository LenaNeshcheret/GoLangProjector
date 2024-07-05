package passengers

import (
	"GoLangProjector/hw6/hw6/route"
	"fmt"
)

type Passenger struct {
	Name string
}

func (p *Passenger) CompleteRoute(r *route.Route) {
	fmt.Printf("%s is traveling by route: %s\n", p.Name, r.Name)
	for _, vehicle := range r.Vehicles {
		fmt.Printf(" %s - %s.\n", vehicle.GetType(), vehicle.GetName())
		if vehicle.AcceptPassenger(1) != true {
			fmt.Printf("You cann't drive this route. \n Because the vehicle %s %s has no free seats. Choose another route.\n",
				vehicle.GetType(), vehicle.GetName())
			return
		} else {
			vehicle.DropOffPassenger(1)
		}
	}
}
