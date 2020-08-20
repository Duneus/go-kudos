package gokudos

type SettingsStorage interface {
	SetScheduleSettings(settings Settings) error
	GetScheduleSettingsForTeam(teamId string) (*Settings, error)
}

type Settings struct {
	TeamId    string
	ChannelId string
}
