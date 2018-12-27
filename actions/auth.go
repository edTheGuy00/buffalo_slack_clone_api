package actions

import (
	"database/sql"
	"log"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/edTheGuy00/slack_clone_backend/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/mw-tokenauth"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(u.Email))).First(u)

	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password")
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("auth/new.html"))
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}
	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "Welcome Back to Buffalo!")

	redirectURL := "/"
	if redir, ok := c.Session().Get("redirectURL").(string); ok {
		redirectURL = redir
	}

	return c.Redirect(302, redirectURL)
}

// AuthCreateToken creates a new token for the user
func AuthCreateToken(u *models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["userid"] = u.ID
	claims["exp"] = time.Now().Add(time.Minute * 5000).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key, err := tokenauth.GetHMACKey(jwt.SigningMethodHS256)
	if err != nil {
		log.Fatal(errors.Wrap(err, "couldn't get key"))
	}
	return token.SignedString(key)
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out!")
	return c.Redirect(302, "/")
}
