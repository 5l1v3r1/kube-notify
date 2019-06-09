package slack

import (
	"fmt"

	"github.com/masahiro331/kube-notify/pkg/config"
	"github.com/nlopes/slack"
)

const (
	colorCyan   = "#00a1e9f"
	colorYellow = "#e8e800"
	colorBlue   = "#2700e8"
	colorHiRed  = "#e80000"
	colorRed    = "#e84d00"
)

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type SlackWriter struct {
}

var Conf config.SlackConf

func Init(conf config.SlackConf) {
	Conf = conf
}

func (w SlackWriter) NotificationResource(kind, name, namespace string) (err error) {
	api := slack.New(Conf.Token)

	str := fmt.Sprintf(`*%s: %s (%s)*`, kind, name, namespace)

	_, _, err = api.PostMessage(
		Conf.Channel,
		slack.MsgOptionText(str, true),
	)
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	return nil
}

func toSlackAttachments() (attaches []slack.Attachment) {
	for _, v := range vs {
		a := slack.Attachment{
			Title:      "Batch Name",
			Text:       "Failed Text or Success Text",
			MarkdownIn: []string{"text", "pretext"},
			Fields: []slack.AttachmentField{
				{
					Title: "Description",
					Value: "",
					Short: true,
				},
			},
			Color: colorCyan,
		}

		attaches = append(attaches, a)
	}
	return attaches
}
