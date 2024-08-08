package flat

type FlatID int

type FlatDTO struct {
	ID        int    `json:"id" validate:"required"`
	HouseID   int    `json:"house_id" validate:"required"`
	Price     int    `json:"price" validate:"required"`
	Rooms     int    `json:"rooms" validate:"required"`
	Moderator string `json:"-"`
	Status    string `json:"status" validate:"required"`
}

type UpdateFlatStatusDTO struct {
	ID      int    `json:"id" validate:"required"`
	HouseID int    `json:"house_id" validate:"required"`
	Status  string `json:"status" validate:"required,oneof_modstat"`
}

type CreateFlatDTO struct {
	HouseID int `json:"house_id" validate:"required"`
	Price   int `json:"price" validate:"required,min=1"`
	Rooms   int `json:"rooms" validate:"required,min=0"`
}

type GetFlatByIDDTO struct {
	ID      int `json:"id" validate:"required"`
	HouseID int `json:"house_id" validate:"required"`
}
