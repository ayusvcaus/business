package graph

import (
	"context"

	"github.com/alecthomas/log4go"

	"github.com/ayusvcaus/business/auth"
	"github.com/ayusvcaus/business/graph/model"
	"github.com/ayusvcaus/business/persist"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return "", &persist.UserError{Err: errStr}
	} else if cuser.Username != "admin" {
		errStr := "The user has no previlige!"
		log4go.Error(errStr)
		return "", &persist.UserError{Err: errStr}
	}
	var user persist.User
	user.Username = input.Username
	user.Password = input.Password
	user.Create()
	if user.IsExist() {
		token, err := auth.GenerateToken(user.Username)
		if err != nil {
			log4go.Error(err.Error())
			return "", err
		}
		log4go.Info(user.Username + " has been created")
		return token, nil
	}
	errStr := "Creating a new User failed"
	log4go.Error(errStr)
	return "", &UserResolverError{Err: errStr}
}

func (r *mutationResolver) UpdatePassword(ctx context.Context, input model.NewUser) (*model.User, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	if cuser.Username != input.Username && cuser.Username != "admin" {
		errStr := "User can chang his/her own password only and admin change anyone's password!"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var user persist.User
	user.Username = input.Username
	user.Password = input.Password
	err := user.UpdatePassword()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	muser := model.User{Name: user.Username, ID: user.ID}
	return &muser, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user persist.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		errStr := "No such username or password"
		log4go.Error(errStr)
		return "", &persist.UserError{Err: errStr}
	}
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		log4go.Error(err.Error())
		return "", err
	}
	log4go.Info(user.Username + " has logged in")
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.NewRefreshToken) (string, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return "", &persist.UserError{Err: errStr}
	}
	username, err := auth.ParseToken(input.Token)
	if username != cuser.Username {
		errStr := "User can refresh his/her own token only!"
		log4go.Error(errStr)
		return "", &persist.UserError{Err: errStr}
	}
	if err != nil {
		log4go.Error(err.Error())
		return "", err
	}
	token, err := auth.GenerateToken(username)
	if err != nil {
		log4go.Error(err.Error())
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) DeleteUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	if cuser.Username != "admin" {
		errStr := "admin can delete a user only!"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var user persist.User
	user.Username = input.Username
	_, err := user.DeleteUser()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	muser := model.User{Name: user.Username, ID: user.ID}
	return &muser, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var dbuser persist.User
	dbusers, _ := dbuser.AllUsers()
	var musers []*model.User
	for _, du := range dbusers {
		var mlinks []*model.Link
		for _, dl := range du.Links {
			ml := model.Link{ID: dl.ID, Title: dl.Title, Address: &dl.Address}
			mlinks = append(mlinks, &ml)
		}
		var mprojects []*model.Project
		for _, dp := range du.Projects {
			mp := model.Project{ID: dp.ID, Name: dp.Name, Quantity: &dp.Quantity, Buget: &dp.Buget}
			mprojects = append(mprojects, &mp)
		}
		mu := model.User{ID: du.ID, Name: du.Username, Links: mlinks, Projects: mprojects}
		musers = append(musers, &mu)

	}
	return musers, nil
}

func (r *queryResolver) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var puser persist.User
	user, err := puser.GetUserByUsername(username)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var mlinks []*model.Link
	for _, dl := range user.Links {
		ml := model.Link{ID: dl.ID, Title: dl.Title, Address: &dl.Address}
		mlinks = append(mlinks, &ml)
	}
	var mprojects []*model.Project
	for _, dp := range user.Projects {
		mp := model.Project{ID: dp.ID, Name: dp.Name, Quantity: &dp.Quantity, Buget: &dp.Buget}
		mprojects = append(mprojects, &mp)
	}
	var muser = model.User{ID: user.ID, Name: user.Username, Links: mlinks, Projects: mprojects}
	return &muser, nil
}

type UserResolverError struct {
	Err string
}

func (m *UserResolverError) Error() string {
	return m.Err
}
