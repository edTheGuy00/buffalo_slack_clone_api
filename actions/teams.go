package actions

import (
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
// Model: Singular (Team)
// DB Table: Plural (teams)
// Resource: Plural (Teams)
// Path: Plural (/teams)
// View Template Folder: Plural (/templates/teams/)

// TeamsResource is the resource for the Team model
type TeamsResource struct {
	buffalo.Resource
}

// List gets all Teams. This function is mapped to the path
// GET /teams
func (v TeamsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context

	//userID := AuthGetUserID(c)
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	user := &models.User{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Teams from the DB
	if err := q.Eager().First(user); err != nil {
		return errors.WithStack(err)
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, user.Teams))
}

// Show gets the data for one Team. This function is mapped to
// the path GET /teams/{team_id}
func (v TeamsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	// To find the Team the parameter team_id is used.
	if err := tx.Eager("channels").Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, team))
}

// New renders the form for creating a new Team.
// This function is mapped to the path GET /teams/new
func (v TeamsResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &models.Team{}))
}

// Create adds a Team to the DB. This function is mapped to the
// path POST /teams
func (v TeamsResource) Create(c buffalo.Context) error {
	// Allocate an empty Team
	team := &models.Team{}

	// Bind team to the html form elements
	if err := c.Bind(team); err != nil {
		return errors.WithStack(err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(team)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, team))
	}

	userID := AuthGetUserID(c)

	uid, err2 := uuid.FromString(userID)
	if err2 != nil {
		return errors.WithStack(err2)
	}

	teamMember := &models.TeamMember{
		UserID: uid,
		TeamID: team.ID,
		Admin:  true,
	}

	if err := tx.Create(teamMember); err != nil {
		return errors.WithStack(err)
	}

	// and redirect to the teams index page
	return c.Render(201, r.Auto(c, team))
}

// Edit renders a edit form for a Team. This function is
// mapped to the path GET /teams/{team_id}/edit
func (v TeamsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	if err := tx.Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, team))
}

// Update changes a Team in the DB. This function is mapped to
// the path PUT /teams/{team_id}
func (v TeamsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	if err := tx.Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind Team to the html form elements
	if err := c.Bind(team); err != nil {
		return errors.WithStack(err)
	}

	verrs, err := tx.ValidateAndUpdate(team)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, team))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", "Team was updated successfully")

	// and redirect to the teams index page
	return c.Render(200, r.Auto(c, team))
}

// Destroy deletes a Team from the DB. This function is mapped
// to the path DELETE /teams/{team_id}
func (v TeamsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	// Allocate an empty Team
	team := &models.Team{}

	// To find the Team the parameter team_id is used.
	if err := tx.Find(team, c.Param("team_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(team); err != nil {
		return errors.WithStack(err)
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", "Team was destroyed successfully")

	// Redirect to the teams index page
	return c.Render(200, r.Auto(c, team))
}
