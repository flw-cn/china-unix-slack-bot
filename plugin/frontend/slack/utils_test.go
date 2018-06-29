package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripDirectMention(t *testing.T) {
	cases := []struct {
		id   int
		text string
		want string
	}{
		{1, "<@U12345678> hello", "hello"},
		{2, "<@U12345678>   hello", "hello"},
		{3, "<@U12345678>hello", "hello"},
		{4, "<@U123oo678> hello", "<@U123oo678> hello"},
	}

	userID := "U12345678"
	for _, c := range cases {
		got := stripDirectMention(c.text, userID)
		assert.Equalf(t, c.want, got, "[case %d] they should equal", c.id)
	}
}
