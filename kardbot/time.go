package kardbot

import (
	"fmt"
	"strings"
	"time"

	"github.com/TannerKvarfordt/Kard-bot/kardbot/dg_helpers"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

const (
	timeCmd             = "time"
	timeCmdOptEphemeral = "ephemeral"
)

const (
	timeSubCmdGroupTZ = "timezones"

	tzSubCmdHelp = "help"

	tzSubCmdInfo           = "info"
	tzSubCmdInfoTZOpt      = "timezone"
	tzSubCmdInfoFmtOpt     = "format"
	tzSubCmdInfoFmtOptDflt = "Monday, 2006-01-02 3:04PM MST"
)

func tzFormatOpts() []*discordgo.ApplicationCommandOptionChoice {
	return []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:  "Default",
			Value: tzSubCmdInfoFmtOptDflt,
		},
		{
			Name:  "Layout",
			Value: time.Layout,
		},
		{
			Name:  "ANSIC",
			Value: time.ANSIC,
		},
		{
			Name:  "UnixDate",
			Value: time.UnixDate,
		},
		{
			Name:  "RubyDate",
			Value: time.RubyDate,
		},
		{
			Name:  "RFC822",
			Value: time.RFC822,
		},
		{
			Name:  "RFC822Z",
			Value: time.RFC822Z,
		},
		{
			Name:  "RFC850",
			Value: time.RFC850,
		},
		{
			Name:  "RFC1123",
			Value: time.RFC1123,
		},
		{
			Name:  "RFC1123Z",
			Value: time.RFC1123Z,
		},
		{
			Name:  "RFC3339",
			Value: time.RFC3339,
		},
		{
			Name:  "RFC3339Nano",
			Value: time.RFC3339Nano,
		},
		{
			Name:  "Kitchen",
			Value: time.Kitchen,
		},
		{
			Name:  "Stamp",
			Value: time.Stamp,
		},
		{
			Name:  "StampMilli",
			Value: time.StampMilli,
		},
		{
			Name:  "StampMicro",
			Value: time.StampMicro,
		},
		{
			Name:  "StampNano",
			Value: time.StampNano,
		},
	}
}

func timeCmdOpts() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
			Name:        timeSubCmdGroupTZ,
			Description: "Timezone related commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        tzSubCmdHelp,
					Description: "Get a list of valid time zones.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        timeCmdOptEphemeral,
							Description: "Should the bot's response be ephemeral? Defaults to true.",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        tzSubCmdInfo,
					Description: "Get information about a given timezone",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        tzSubCmdInfoTZOpt,
							Description: "The IANA timezone to get information for.",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        tzSubCmdInfoFmtOpt,
							Description: "The format in which the date should be displayed.",
							Choices:     tzFormatOpts(),
						},
						{
							Type:        discordgo.ApplicationCommandOptionBoolean,
							Name:        timeCmdOptEphemeral,
							Description: "Should the bot's response be ephemeral? Defaults to true.",
						},
					},
				},
			},
		},
	}
}

func handleTimeCmd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if s == nil || i == nil {
		err := fmt.Errorf("nil Session pointer (%v) and/or InteractionCreate pointer (%v)", s, i)
		interactionRespondEphemeralError(s, i, true, err)
		log.Error(err)
		return
	}

	var (
		err           error                          = nil
		reportableErr                                = false
		resp          *discordgo.InteractionResponse = nil
	)
	subCmdOrGroup := i.ApplicationCommandData().Options[0].Name
	switch subCmdOrGroup {
	case timeSubCmdGroupTZ:
		resp, reportableErr, err = handleTZSubCmd(s, i)
	default:
		interactionRespondEphemeralError(s, i, true, fmt.Errorf("unknown subcommand: %s", subCmdOrGroup))
		return
	}

	if err != nil {
		interactionRespondEphemeralError(s, i, reportableErr, err)
		return
	}
	if resp == nil {
		interactionRespondEphemeralError(s, i, true, fmt.Errorf("nil response returned"))
		log.Error(err)
		return
	}

	err = s.InteractionRespond(i.Interaction, resp)
	if err != nil {
		interactionRespondEphemeralError(s, i, true, err)
		log.Error(err)
		return
	}
}

func handleTZSubCmd(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, bool, error) {
	subCmdName := i.ApplicationCommandData().Options[0].Options[0].Name
	switch subCmdName {
	case tzSubCmdHelp:
		return handleTZSubCmdHelp(s, i)
	case tzSubCmdInfo:
		return handleTZSubCmdInfo(s, i)
	default:
		return nil, true, fmt.Errorf("unknown %s sub command: %s", timeSubCmdGroupTZ, subCmdName)
	}
}

func handleTZSubCmdHelp(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, bool, error) {
	flags := InteractionResponseFlagEphemeral
	for _, opt := range i.ApplicationCommandData().Options[0].Options[0].Options {
		switch opt.Name {
		case timeCmdOptEphemeral:
			if !opt.BoolValue() {
				flags = 0
			}
		default:
			log.Warn("Unknown option: ", opt.Name)
		}
	}

	c, _ := fastHappyColorInt64()
	e := dg_helpers.NewEmbed()
	e.SetTitle("Timezones").
		SetURL("https://en.wikipedia.org/wiki/List_of_tz_database_time_zones").
		SetColor(int(c)).
		SetDescription("This bot supports [Internet Assigned Numbers Authority (IANA)](https://www.iana.org/time-zones) governed timezones. "+
			"A convenient list of valid timezones can be found on [Wikipedia](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List). "+
			"See the _TZ database name_ or _Time zone abbreviation_ column. Note that inputs are case-sensitive.\n"+
			"\n*Valid Timezone Input Examples*\n"+
			"- America/Boise\n"+
			"- Asia/Hong_Kong\n"+
			"- Europe/Berlin\n"+
			"- EET\n"+
			"- MST\n"+
			"- MDT\n"+
			"\nSubcommands and their usage are documented below.\n").
		AddField(tzSubCmdHelp, "Prints this help message. Response is optionally ephemeral.").
		AddField(tzSubCmdInfo, "Provides general information about a given timezone. "+
			"Requires an [IANA timezone database name or abbreviation](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List) as input. "+
			"Optionally takes a date format in which the provided timezone should be displayed. "+
			"Response is optionally ephemeral.")

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  flags,
			Embeds: []*discordgo.MessageEmbed{e.Truncate().SetType(discordgo.EmbedTypeRich).MessageEmbed},
		},
	}, false, nil
}

func handleTZSubCmdInfo(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, bool, error) {
	flags := InteractionResponseFlagEphemeral
	tz := ""
	format := tzSubCmdInfoFmtOptDflt
	for _, opt := range i.ApplicationCommandData().Options[0].Options[0].Options {
		switch opt.Name {
		case timeCmdOptEphemeral:
			if !opt.BoolValue() {
				flags = 0
			}
		case tzSubCmdInfoTZOpt:
			tz = opt.StringValue()
		case tzSubCmdInfoFmtOpt:
			format = opt.StringValue()
		default:
			log.Warn("Unknown option: ", opt.Name)
		}
	}

	tz = strings.TrimSpace(tz)
	if strings.ToLower(tz) == "local" {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   InteractionResponseFlagEphemeral,
				Content: fmt.Sprintf(`For privacy reasons, this bot does not track user timezones. Please specify a specific IANA timezone rather than "%s".`, tz),
			},
		}, false, nil
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   InteractionResponseFlagEphemeral,
				Content: fmt.Sprintf(`"%s" is not a valid [IANA Timezone](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones).`, tz),
			},
		}, false, nil
	}

	c, _ := fastHappyColorInt64()
	e := dg_helpers.NewEmbed()
	t := time.Now().In(loc)
	abbrev, offset := t.Zone()
	e.SetTitle(loc.String()).
		SetDescription(t.Format(format)).
		SetColor(int(c)).
		AddField("Abbreviation", abbrev).
		AddField("Daylight Savings Time in Effect?", fmt.Sprintf("%t", t.IsDST())).
		AddField("UTC/GMT Offset (hh:mm)", fmt.Sprintf(`%+03d:%02d`, offset/3600, func() int {
			seconds := (offset % 3600) / 60
			if seconds < 0 {
				return seconds * -1
			}
			return seconds
		}()))

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  flags,
			Embeds: []*discordgo.MessageEmbed{e.Truncate().SetType(discordgo.EmbedTypeRich).MessageEmbed},
		},
	}, false, nil
}
