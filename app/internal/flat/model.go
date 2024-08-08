package flat

import "github.com/Polyrom/houses_api/internal/house"

type Flat struct {
	ID     int         `json:"id"`
	House  house.House `json:"house_id"`
	Price  int         `json:"price"`
	Rooms  int         `json:"rooms"`
	Status string      `json:"status"`
}
