package domain

type Role int

const (
	Guest     Role = iota // 0
	Usr                   // 1
	Moderator             // 2
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	//Password  []byte             `json:"-"`
	Password  []byte `json:"password"`
	Email     string `json:"email"`
	ImagePath string `json:"imagePath"`
	ImageData []byte `json:"imageData"`
	Role      Role
}

type ProfileUsecase interface {
	GetUserData(userID int) (User, error)
	UpdateUser(newUser User) (User, error)
	UploadImage(userID int, imageData []byte) (string, error)
}

type ProfileRepository interface {
	GetUser(userID int) (User, error)
	UpdateUser(newUser User) (User, error)
}
