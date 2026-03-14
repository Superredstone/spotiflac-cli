package lib

import "errors"

type AvailableServices int

const (
	AvailableServicesTidal = iota
)

func ParseAvailableServices(service string) (AvailableServices, error) {
	switch service {
	case "tidal":
		return AvailableServicesTidal, nil
	}
	return AvailableServicesTidal, errors.New("Invalid service.")
}

