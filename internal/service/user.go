package service

type userCred struct {
	Id       int64
	Email    string
	Password string
}

type User struct {
	UserCred *userCred
	jwt      *Claims
}

func NewUser() *User {
	return &User{}
}

func (u *User) CreateUser() error {
	return nil
}

func (u *User) UserAuthentication() error {
	//db, err := u.storage.Connect()
	//if err != nil {
	//	return
	//}
	//
	//token, err := u.jwt.CreateNewToken(u.email)
	//if err != nil {
	//	return
	//}
	return nil
}
