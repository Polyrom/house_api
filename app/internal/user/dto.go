package user

type UserRegisterDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"user_type" validate:"required,oneof=client moderator"`
}

type UserIDDTO struct {
	UserID UserID `json:"user_id"`
}

type TokenDTO struct {
	Token Token `json:"token"`
}

type UserLoginDTO struct {
	UserID   UserID `json:"user_id" validate:"required"`
	Password string `json:"password" validate:"required"`
}
