package types

type Channel struct {
	ID   string
	Name string
}

func (c Channel) String() string {
	if c.Name != "" {
		return "#" + c.Name
	} else {
		return "#<" + c.ID + ">"
	}
}

type User struct {
	ID   string
	Name string
}

func (u User) String() string {
	if u.Name != "" {
		return u.Name
	} else {
		return "<" + u.ID + ">"
	}
}
