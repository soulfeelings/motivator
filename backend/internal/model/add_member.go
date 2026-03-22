package model

type AddMemberRequest struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=6"`
	Role        Role    `json:"role"`
	DisplayName *string `json:"display_name,omitempty"`
}

type AddMemberResponse struct {
	UserID     string     `json:"user_id"`
	Email      string     `json:"email"`
	Membership Membership `json:"membership"`
}
