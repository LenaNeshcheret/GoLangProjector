package publicTransport

type PublicTransport interface {
	AcceptPassenger(passengers int) bool
	DropOffPassenger(passengers int)
	GetName() string
	GetType() string
}
