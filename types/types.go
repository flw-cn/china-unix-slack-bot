package types

type Channel struct {
	ID   string
	Name string
}

func (c Channel) String() string {
	var str string

	if c.Name != "" {
		str = "#" + c.Name
	} else {
		str = "#<" + c.ID + ">"
	}

	return str
}

type User struct {
	ID   string
	Name string
}

func (u User) String() string {
	var str string

	if u.Name != "" {
		str = u.Name
	} else {
		str = "<" + u.ID + ">"
	}

	return str
}
