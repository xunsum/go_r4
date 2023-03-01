package models

type User struct {
	Uid                  string `json:"uid"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	EncryptedPassword    string `json:"encrypted_password"`
	Slogan               string `json:"slogan"`
	ProfileImageLocation string `json:"profile_image_location"`
	IsBlocked            int    `json:"is_blocked"`
}

func (User) TableName() string {
	return "users"
}
