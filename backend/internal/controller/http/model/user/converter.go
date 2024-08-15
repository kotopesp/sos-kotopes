package user

import "github.com/kotopesp/sos-kotopes/internal/core"

func (u *User) ToCoreUser() core.User {
	return core.User{
		Username:     u.Username,
		PasswordHash: u.Password,
		Description:  u.Description,
		Photo:        u.Photo,
		Firstname:    u.Firstname,
		Lastname:     u.Lastname,
	}
}
func (u *UpdateUser) ToCoreUpdateUser() core.UpdateUser {
	if u == nil {
		return core.UpdateUser{}
	}
	return core.UpdateUser{
		Username:     u.Username,
		PasswordHash: u.Password,
		Description:  u.Description,
		Photo:        u.Photo,
		Firstname:    u.Firstname,
		Lastname:     u.Lastname,
	}
}

func ToResponseUser(user *core.User) ResponseUser {
	if user == nil {
		return ResponseUser{}
	}
	return ResponseUser{
		ID:          user.ID,
		Username:    user.Username,
		Lastname:    user.Lastname,
		Firstname:   user.Firstname,
		Photo:       user.Photo,
		Description: user.Description,
	}
}
