// Package db defines errors used in the database layer.
package db

import "errors"

var (
	// ErrorTeamNotFound - ошибка, команда не найдена
	ErrorTeamNotFound = errors.New("resource not found")
	// ErrorTeamAlreadyExists - ошибка, команда уже существует
	ErrorTeamAlreadyExists = errors.New("team_name already exists")
	// ErrorUserNotFound - ошибка, пользователь не найден
	ErrorUserNotFound = errors.New("user not found")
	// ErrorPRSNotFound - ошибка, PR не найден
	ErrorPRSNotFound = errors.New("pull request not found")
	// ErrorPRAlreadyExists - ошибка, PR уже существует
	ErrorPRAlreadyExists = errors.New("PR id already exists")
	// ErrorPRMerged - ошибка, PR уже был объединен
	ErrorPRMerged = errors.New("cannot reassign on merged PR")
	// ErrorReviewerNotAssigned - ошибка, ревьювер не назначен
	ErrorReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
	// ErrorNoCandidateForReviewer - ошибка, нет кандидата для ревьювера
	ErrorNoCandidateForReviewer = errors.New("no active replacement candidate in team")
)

var (
	// CodeTeamNotFound - код ошибки, команда не найдена
	CodeTeamNotFound = "NOT_FOUND"
	// CodeTeamAlreadyExists - код ошибки, команда уже существует
	CodeTeamAlreadyExists = "TEAM_EXISTS"
	// CodePRExists - код ошибки, PR уже существует
	CodePRExists = "PR_EXISTS"
	// CodePRMerged - код ошибки, PR уже был объединен
	CodePRMerged = "PR_MERGED"
	// CodeNotAssigned - код ошибки, ревьювер не назначен
	CodeNotAssigned = "NOT_ASSIGNED"
	// CodeNoCandidate - код ошибки, нет кандидата для ревьювера
	CodeNoCandidate = "NO_CANDIDATE"
)
