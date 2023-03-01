package models

type Admin struct {
	Uid               string `json:"uid"`
	Name              string `json:"name"`
	EncryptedPassword string `json:"encrypted_password"`
}

func (Admin) TableName() string {
	return "admins"
}
