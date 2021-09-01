package kardbot

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"

	"github.com/bwmarrin/discordgo"
)

var (
	// Possible greetings the bot will respond to
	// TODO: make these configurable
	greetings = []string{
		"Hello",
		"Hi",
		"Greetings",
		"Salutations",
	}

	// Possible farewells the bot will respond to
	// TODO: make these configurable
	farewells = []string{
		"Goodbye",
		"Farewell",
		"So long",
		"Bye",
		"See you",
		"See ya",
	}

	// Any callbacks that happen onMessageCreate belong in this list.
	// It is the duty of each individual function to decide whether or not to run.
	// These callbacks must be able to safely execute asynchronously.
	onCreateHandlers = [...]onCreateHandler{
		greeting,
		farewell,
	}
)

type onCreateHandler = func(*discordgo.Session, *discordgo.MessageCreate)

func msgIsFromSelf(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	return m.Author.ID == s.State.User.ID
}

func greeting(s *discordgo.Session, m *discordgo.MessageCreate) {
	if msgIsFromSelf(s, m) {
		return
	}

	greetingGroup := buildRegexAltGroup(greetings)

	matched, err := regexp.MatchString(
		fmt.Sprintf("^(?i)%s %s[!.]*$", greetingGroup, buildBotNameRegexp(s.State.User.Username)),
		m.Content,
	)
	if err != nil {
		log.Println("Regex error: ", err)
		return
	}
	if matched {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("%s %s!", greetings[rand.Intn(len(greetings))], m.Author.Username),
		)
	}
}

func farewell(s *discordgo.Session, m *discordgo.MessageCreate) {
	if msgIsFromSelf(s, m) {
		return
	}

	farewellGroup := buildRegexAltGroup(farewells)

	matched, err := regexp.MatchString(
		fmt.Sprintf("^(?i)%s %s[!.]*$", farewellGroup, buildBotNameRegexp(s.State.User.Username)),
		m.Content,
	)
	if err != nil {
		log.Println("Regex error: ", err)
		return
	}
	if matched {
		s.ChannelMessageSend(
			m.ChannelID,
			fmt.Sprintf("%s %s!", farewells[rand.Intn(len(farewells))], m.Author.Username),
		)
	}
}