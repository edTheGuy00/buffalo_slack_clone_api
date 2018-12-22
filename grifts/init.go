package grifts

import (
	"github.com/edTheGuy00/slack_clone_backend/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
