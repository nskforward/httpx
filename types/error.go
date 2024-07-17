package types

type Error struct {
	Status     int
	Text       string
	StackTrace string
}

func (e Error) Error() string {
	return e.Text
}
