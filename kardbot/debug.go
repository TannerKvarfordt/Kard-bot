package kardbot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Map logrus log levels to discordgo log levels
func logrusToDiscordGo() map[log.Level]int {
	return map[log.Level]int{
		log.PanicLevel: discordgo.LogError,
		log.FatalLevel: discordgo.LogError,
		log.ErrorLevel: discordgo.LogError,
		log.WarnLevel:  discordgo.LogWarning,
		log.InfoLevel:  discordgo.LogInformational,
		log.DebugLevel: discordgo.LogInformational,
		log.TraceLevel: discordgo.LogDebug,
	}
}

func updateLogLevel(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// TODO: move this check into a helper function
	author, authorID, err := getInteractionCreateAuthorNameAndID(i)
	if err != nil {
		log.Error(err)
		return
	}
	if authorID == s.State.User.ID {
		log.Trace("Ignoring message from self")
		return
	}

	if isOwner, err := authorIsOwner(i); err != nil {
		log.Error(err)
		return
	} else if !isOwner {
		log.Warnf("User %s (%s) does not have privilege to update log level", author, authorID)
		return
	}

	levelStr := strings.ToLower(i.ApplicationCommandData().Options[0].StringValue())

	if lvl, err := log.ParseLevel(levelStr); err == nil {
		info := fmt.Sprintf(`Set logging level to "%s"`, levelStr)
		log.Info(info)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: info,
			},
		})
		log.SetLevel(lvl)
		if bot().EnableDGLogging {
			// TODO: make this thread safe somehow (logrus is already thread safe)
			s.LogLevel = logrusToDiscordGo()[lvl]
		}
	} else {
		log.Error(err)
	}
}
