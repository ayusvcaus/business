package graph

import (
	"context"
	"strconv"

	"github.com/alecthomas/log4go"

	"github.com/ayusvcaus/business/auth"
	"github.com/ayusvcaus/business/persist"
)

func (r *mutationResolver) AssociateProject2User(ctx context.Context, projectID int, userID int) (string, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return "", &persist.UserError{Err: errStr}
	}
	project_users := persist.Project_users{ProjectID: projectID, UserID: userID}
	succeeded, err := project_users.AssociateProject2User()
	if err != nil {
		log4go.Error(err.Error())
		return "", err
	}
	return strconv.Itoa(projectID) + " and " + strconv.Itoa(userID) + " " + succeeded, nil

}

func (r *mutationResolver) DecoupleProject2User(ctx context.Context, projectID int, userID int) (string, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return "", &persist.UserError{Err: errStr}
	}
	project_users := persist.Project_users{ProjectID: projectID, UserID: userID}
	succeeded, err := project_users.Delete()
	if err != nil {
		log4go.Error(err.Error())
		return "", err
	}
	return strconv.Itoa(projectID) + " and " + strconv.Itoa(userID) + " mapping was " + succeeded, nil

}
