package response

type Error struct {
	Status int
	Text   string
}

func (e Error) Error() string {
	return e.Text
}
