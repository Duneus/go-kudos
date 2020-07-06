package service

import (
	"github.com/Duneus/go-kudos/pkg/gokudos"
	"github.com/slack-go/slack"
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
	hideKudosText := slack.NewTextBlockObject("plain_text", "Hide kudos", false, false)
	showKudos := slack.NewButtonBlockElement("show_kudos", "show_kudos", showKudosText)
	hideKudos := slack.NewButtonBlockElement("hide_kudos", "hide_kudos", hideKudosText)
	action := slack.NewActionBlock("actions", showKudos, hideKudos)

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
			kudosText := slack.NewTextBlockObject("plain_text", k.Message, true, false)
			removeButtonText := slack.NewTextBlockObject("plain_text", "Remove", false, false, )
			removeButton := slack.NewButtonBlockElement("remove_kudos", k.ID, removeButtonText)
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
