package actions

import (
	"github.com/edTheGuy00/slack_clone_backend/actions/websocket"
	"github.com/edTheGuy00/slack_clone_backend/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Message)
// DB Table: Plural (messages)
// Resource: Plural (Messages)
// Path: Plural (/messages)
// View Template Folder: Plural (/templates/messages/)

// MessagesResource is the resource for the Message model
type MessagesResource struct {
	buffalo.Resource
}

var hub *websocket.Hub

// MessagesHandler handles the websocket connection for messages
func MessagesHandler(c buffalo.Context) error {

	if hub == nil {
		hub = websocket.NewHub()
		go hub.Run()
	}

	return websocket.ServeWS(hub, c.Response(), c.Request())
}

// List gets the list of messages from the database
func (v MessagesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	messages := &models.Messages{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Messages from the DB
	if err := q.Eager("User").Where("channel_id = ?", c.Param("channel_id")).All(messages); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, messages))
}

// Create adds a Message to the DB. This function is mapped to the
// path POST /messages
func (v MessagesResource) Create(c buffalo.Context) error {
	// Allocate an empty Message
	message := &models.Message{}

	// Bind message to the html form elements
	if err := c.Bind(message); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	channelID, err := uuid.FromString(c.Param("channel_id"))
	if err != nil {
		return errors.WithStack(err)
	}

	uid := AuthGetUserID(c)
	userID, err := uuid.FromString(uid)
	if err != nil {
		return errors.WithStack(err)
	}

	message.ChannelID = channelID
	message.UserID = userID

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(message)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, message))
	}

	// and redirect to the messages index page
	return c.Render(201, r.Auto(c, message))
}

// Destroy deletes a Message from the DB. This function is mapped
// to the path DELETE /messages/{message_id}
func (v MessagesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Message
	message := &models.Message{}

	// To find the Message the parameter message_id is used.
	if err := tx.Find(message, c.Param("message_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(message); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Message was destroyed successfully")

	// Redirect to the messages index page
	return c.Render(200, r.Auto(c, message))
}
