package slack

import (
	"regexp"
	"strings"
)

// There are some basic CONCEPT need to be determined first:
//
// We said a message is X message means:
//
// DirectMessage(aka DM, IM) means someone talks with you alone.
// DirectMention means someone mentions you at the very beginning of the message.
// Mentioned means someone mentions you in the message.
// Message means someone said something in the message but did not mention you.

// isDirectMessage returns true if this message is in a direct message conversation
func isDirectMessage(channel string) bool {
	return strings.HasPrefix(channel, "D")
}

// isDirectMention returns true is message is a Direct Mention that mentions a specific user.
// A direct mention is a mention at the very beginning of the message.
func isDirectMention(text string, userID string) bool {
	return strings.HasPrefix(text, `<@`+userID+`>`)
}

// isMentioned returns true if this message contains a mention of a specific user
func isMentioned(text string, userID string) bool {
	return strings.Contains(text, `<@`+userID+`>`)
}

// whoMentioned returns a list of userIDs mentioned in the message
func whoMentioned(text string) []string {
	r, rErr := regexp.Compile(`<@(U[a-zA-Z0-9]+)>`)
	if rErr != nil {
		return []string{}
	}

	allMatches := r.FindAllStringSubmatch(text, -1)
	dict := make([]string, len(allMatches))
	for i, r := range allMatches {
		dict[i] = r[1]
	}

	return dict
}

// stripDirectMention removes a leading mention (aka direct mention) from a message string
func stripDirectMention(text string, userID string) string {
	r, rErr := regexp.Compile(`(?s)^(<@` + userID + `>[\:]*[\s]*)?`)
	if rErr != nil {
		return ""
	}

	return r.ReplaceAllString(text, "")
}

// stripAllMention removes all mention from a message string
func stripAllMention(text string) string {
	r, rErr := regexp.Compile(`<@U[a-zA-Z0-9]+>`)
	if rErr != nil {
		return ""
	}

	return strings.TrimSpace(r.ReplaceAllString(text, ""))
}
