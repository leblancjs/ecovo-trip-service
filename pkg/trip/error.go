package trip

// A NotFoundError is an error that represents that no trip was found.
type NotFoundError struct {
	msg string
}

func (e NotFoundError) Error() string {
	return e.msg
}
