package db

import "errors"

var (
	ErrorTeamNotFound = errors.New("resource not found")
	ErrorTeamAlreadyExists = errors.New("team_name already exists")
	ErrorUserNotFound = errors.New("user not found")
	ErrorPRSNotFound = errors.New("pull request not found")
)

var (
	CodeTeamNotFound        = "NOT_FOUND"
	CodeTeamAlreadyExists  = "TEAM_EXISTS"
)
