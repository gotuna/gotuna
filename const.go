package gotuna

type constError string

func (err constError) Error() string {
	return string(err)
}
