package service

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"github.com/slack-go/slack"
	"math/rand"
	"strconv"
)

func NewDivider() slack.Block {
	return slack.DividerBlock{
		Type: "divider",
	}
}

func NewSection(text string) slack.Block {
	return slack.SectionBlock{
		Type: "section",
		Text: &slack.TextBlockObject{
			Type: "mrkdwn",
			Text: text,
		},
	}
}

func handleAppHomeTab() (slack.HomeTabViewRequest, error) {
	hello := NewSection("hello!")
	divider := NewDivider()
	prompt := NewSection("hello!")

	showKudosText := slack.NewTextBlockObject("plain_text", "Show kudos", false, false)
	showAllKudosText := slack.NewTextBlockObject("plain_text", "Show all kudos", false, false)
	hideKudosText := slack.NewTextBlockObject("plain_text", "Hide kudos", false, false)
	showKudos := slack.NewButtonBlockElement("show_kudos", "show_kudos", showKudosText)
	showAllKudos := slack.NewButtonBlockElement("show_all_kudos", "show_all_kudos", showAllKudosText)
	hideKudos := slack.NewButtonBlockElement("hide_kudos", "hide_kudos", hideKudosText)
	action := slack.NewActionBlock("actions", showKudos, showAllKudos, hideKudos)

	return slack.HomeTabViewRequest{
		Type: "home",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{hello, divider, prompt, action},
		},
	}, nil
}

func handleAppHomeTabWithKudosList(kudos []gokudos.Kudos) (slack.HomeTabViewRequest, error) {
	var blockSet []slack.Block
	hello := NewSection("hello!")
	divider := NewDivider()
	prompt := NewSection("hello!")

	hideKudosText := slack.NewTextBlockObject("plain_text", "Hide kudos", false, false)
	hideKudos := slack.NewButtonBlockElement("hide_kudos", "hide_kudos", hideKudosText)
	action := slack.NewActionBlock("actions", hideKudos)

	blockSet = append(blockSet, hello, divider, prompt, action, divider)

	if len(kudos) > 0 {
		for _, k := range kudos {
			kudosId := strconv.Itoa(k.ID)
			kudosText := slack.NewTextBlockObject("mrkdwn", k.Message, false, false)
			removeButtonText := slack.NewTextBlockObject("plain_text", "Remove", false, false, )
			removeButton := slack.NewButtonBlockElement("remove_kudos", kudosId, removeButtonText)
			acc := slack.Accessory{
				ButtonElement: removeButton,
			}
			kudosSection := slack.NewSectionBlock(kudosText, nil, &acc)
			blockSet = append(blockSet, kudosSection)
		}
	}

	return slack.HomeTabViewRequest{
		Type: "home",
		Blocks: slack.Blocks{
			BlockSet: blockSet,
		},
	}, nil
}

var hearts = []string{":purple_heart:", ":heart:", ":yellow_heart:", ":green_heart:", ":blue_heart:"}

func getHeart() string {
	el := rand.Intn(len(hearts))

	return hearts[el]
}

func prependHeart(message string) string {
	return getHeart() + " " + message
}
