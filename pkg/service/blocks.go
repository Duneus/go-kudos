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

func homeTabHeader() *slack.ActionBlock {
	showKudosText := slack.NewTextBlockObject("plain_text", "Show kudos", false, false)
	showAllKudosText := slack.NewTextBlockObject("plain_text", "Show all kudos", false, false)
	scheduleText := slack.NewTextBlockObject("plain_text", "Schedule next kudos", false, false)
	channelSelectText := slack.NewTextBlockObject("plain_text", "Select channel", false, false)

	showKudos := slack.NewButtonBlockElement(ShowKudosView, "show_kudos", showKudosText)
	showAllKudos := slack.NewButtonBlockElement(ShowAllKudosView, "show_all_kudos", showAllKudosText)
	schedule := slack.NewButtonBlockElement(ShowSchedulingView, "schedule", scheduleText)
	channelSelect := slack.NewButtonBlockElement(ShowChannelSelectView, "channel_select", channelSelectText)

	return slack.NewActionBlock("actions", showKudos, showAllKudos, schedule, channelSelect)
}

func handleAppHomeTab() (slack.HomeTabViewRequest, error) {
	hello := NewSection("Welcome to KudosBot")
	divider := NewDivider()
	prompt := NewSection("hello!")

	header := homeTabHeader()

	return slack.HomeTabViewRequest{
		Type: "home",
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{hello, divider, header, prompt},
		},
	}, nil
}

func handleAppHomeTabWithKudosList(kudos []gokudos.Kudos) (slack.HomeTabViewRequest, error) {
	var blockSet []slack.Block
	hello := NewSection("Welcome to KudosBot")
	divider := NewDivider()
	prompt := NewSection("hello!")

	header := homeTabHeader()

	blockSet = append(blockSet, hello, divider, header, prompt)

	if len(kudos) > 0 {
		for _, k := range kudos {
			kudosId := strconv.Itoa(k.ID)
			kudosText := slack.NewTextBlockObject("mrkdwn", k.Message, false, false)
			removeButtonText := slack.NewTextBlockObject("plain_text", "Remove", false, false, )
			removeButton := slack.NewButtonBlockElement(RemoveKudos, kudosId, removeButtonText)
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

func handleAppHomeSchedulingTab() (slack.HomeTabViewRequest, error) {
	var blockSet []slack.Block
	hello := NewSection("Welcome to KudosBot")
	divider := NewDivider()
	prompt := NewSection("Select when do you want to schedule your next posting of kudos!")

	header := homeTabHeader()

	datepicker := slack.NewDatePickerBlockElement(SetSchedule)
	datepickerText := slack.NewTextBlockObject("plain_text", "Select date", false, false)

	acc := slack.Accessory{
		DatePickerElement: datepicker,
	}

	datePickerSection := slack.NewSectionBlock(datepickerText, nil, &acc)

	blockSet = append(blockSet, hello, divider, header, prompt, datePickerSection)

	return slack.HomeTabViewRequest{
		Type: "home",
		Blocks: slack.Blocks{
			BlockSet: blockSet,
		},
	}, nil
}

func handleAppHomeTabChannelPicker(selectedChannel string) (slack.HomeTabViewRequest, error) {
	var blockSet []slack.Block
	hello := NewSection("Welcome to KudosBot")
	divider := NewDivider()
	prompt := NewSection("Please select which channel you want to use to publish kudos!")

	header := homeTabHeader()

	text := slack.NewTextBlockObject("plain_text", "Select channel", false, false)

	dropdown := slack.NewOptionsSelectBlockElement("channels_select", text, SetChannel)

	if selectedChannel != "" {
		dropdown.InitialChannel = selectedChannel
	}

	acc := slack.NewAccessory(dropdown)

	dropdownSection := slack.NewSectionBlock(text, nil, acc)

	blockSet = append(blockSet, hello, divider, header, prompt, dropdownSection)

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
