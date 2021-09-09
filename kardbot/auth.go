package kardbot

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	BotTokenEnv = "KARDBOT_TOKEN"
	BotOwnerEnv = "KARDBOT_OWNER_ID"
	Intents     = discordgo.IntentsAllWithoutPrivileged
)

var (
	getBotToken func() string
	getOwnerID  func() string
)

// Retrieves the bot's auth token from the environment
func init() {
	// This will only add new environment variables,
	// and will NOT overwrite existing ones.
	_ = godotenv.Load( /*.env by default*/ )

	token, tokenFound := os.LookupEnv(BotTokenEnv)
	if !tokenFound {
		log.Fatalf("%s not found in environment", BotTokenEnv)
	} else if token == "" {
		log.Fatalf("%s is the empty string", BotTokenEnv)
	}
	getBotToken = func() string { return token }

	owner, ownerFound := os.LookupEnv(BotOwnerEnv)
	if !ownerFound {
		log.Warnf("%s not found in environment. No commands requiring this privilege can be executed.", BotOwnerEnv)
	} else if owner == "" {
		log.Warnf("%s is the empty string. No commands requiring this privilege can be executed.", BotOwnerEnv)
	}
	getOwnerID = func() string { return owner }
}
