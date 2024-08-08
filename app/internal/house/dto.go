package house

type CreateHouseDTO struct {
	Address   string `json:"address" validate:"required"`
	Year      int    `json:"year" validate:"required"`
	Developer string `json:"developer"`
}
