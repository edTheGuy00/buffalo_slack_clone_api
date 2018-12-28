package actions

import (
	"github.com/edTheGuy00/slack_clone_backend/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		return c.Render(406, r.JSON(verrs))
	}

	token, err := AuthCreateToken(u)
	if err != nil {
		return errors.WithStack(err)
	}
	u.JwtToken = token

	return c.Render(200, r.JSON(u))
}
