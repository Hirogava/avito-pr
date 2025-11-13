package db

import "errors"

var (
	ErrorTeamNotFound = errors.New("resource not found")
	ErrorTeamAlreadyExists = errors.New("team_name already exists")
	ErrorUserNotFound = errors.New("user not found")
	ErrorPRSNotFound = errors.New("pull request not found")
	ErrorPRAlreadyExists = errors.New("PR id already exists")
	ErrorPRMerged = errors.New("cannot reassign on merged PR")
	ErrorReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrorNoCandidateForReviewer = errors.New("no active replacement candidate in team")
)

var (
	CodeTeamNotFound        = "NOT_FOUND"
	CodeTeamAlreadyExists  = "TEAM_EXISTS"
	CodePRExists          = "PR_EXISTS"
	CodePRMerged 		 = "PR_MERGED"
	CodeNotAssigned 	= "NOT_ASSIGNED"
	CodeNoCandidate    = "NO_CANDIDATE"
)
