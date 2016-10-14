package module

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/johntdyer/slackrus"
)

// Slack ...
type Slack struct {
	Errmsg string
	Info   string
}

func (slack *Slack) SendSlack() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stderr)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.AddHook(&slackrus.SlackrusHook{
		HookURL:        "https://hooks.slack.com/services/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",
		AcceptedLevels: slackrus.LevelThreshold(logrus.DebugLevel),
		Channel:        "#go_push",
		IconEmoji:      ":ghost:",
		Username:       "senseoki",
	})

	logrus.WithFields(logrus.Fields{
		"Go":  "PUSH",
		"msg": slack.Errmsg,
	}).Info(slack.Info)
}
