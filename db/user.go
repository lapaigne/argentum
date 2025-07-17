package db

type User struct {
	Id     int
	Worker int
	Login  string
	Hash   string
	Level  int
}

var AccessLevels = map[string]int{
	"admin":      100,
	"dispatcher": 50,
	"worker":     10,
}

func GetUser(login string) (User, error) {

	var u User
	query := "SELECT * FROM public.users WHERE login = $1"
	err := db.QueryRow(query, login).Scan(&u.Id, &u.Worker, &u.Login, &u.Hash, &u.Level)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func AddUser(u User) error {

	// query := "INSERT INTO public.users"

	return nil
}
