package db

import "errors"

var (
	ErrorTeamNotFound = errors.New("resource not found")
	ErrorTeamAlreadyExists = errors.New("team_name already exists")
)

var (
	CodeTeamNotFound        = "NOT_FOUND"
	CodeTeamAlreadyExists  = "TEAM_EXISTS"
)