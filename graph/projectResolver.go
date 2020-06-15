package graph

import (
	"context"

	"github.com/alecthomas/log4go"

	"github.com/ayusvcaus/business/auth"
	"github.com/ayusvcaus/business/graph/model"
	"github.com/ayusvcaus/business/persist"
)

func (r *mutationResolver) CreateProject(ctx context.Context, input model.NewProject) (*model.Project, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	project := persist.Project{Name: input.Name,
		Quantity: input.Quantity,
		Buget:    input.Buget}
	dbproject, err := project.Create()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	mproject := model.Project{ID: dbproject.ID,
		Name:     dbproject.Name,
		Quantity: &dbproject.Quantity,
		Buget:    &dbproject.Buget}
	return &mproject, nil
}

func (r *mutationResolver) UpdateProject(ctx context.Context, projectID int, input model.NewProject) (*model.Project, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	project := persist.Project{ID: projectID,
		Name:     input.Name,
		Quantity: input.Quantity,
		Buget:    input.Buget}
	dbproject, err := project.Update()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	mproject := model.Project{ID: dbproject.ID,
		Name:     dbproject.Name,
		Quantity: &dbproject.Quantity,
		Buget:    &dbproject.Buget}
	return &mproject, nil
}

func (r *mutationResolver) DeleteProject(ctx context.Context, projectID int) (*model.Project, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	project := persist.Project{ID: projectID}
	_, err := project.Delete()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	mproject := model.Project{ID: project.ID,
		Name:     project.Name,
		Quantity: &project.Quantity,
		Buget:    &project.Buget}
	return &mproject, nil
}

func (r *queryResolver) Projects(ctx context.Context) ([]*model.Project, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var project persist.Project
	dbprojects, _ := project.Projects()
	var mprojects []*model.Project
	for _, dp := range dbprojects {
		var musers []*model.User
		for _, du := range dp.Users {
			mu := model.User{ID: du.ID, Name: du.Username}
			musers = append(musers, &mu)
		}
		mp := model.Project{ID: dp.ID,
			Name:     dp.Name,
			Quantity: &dp.Quantity,
			Buget:    &dp.Buget,
			Users:    musers}
		mprojects = append(mprojects, &mp)
	}
	return mprojects, nil
}

func (r *queryResolver) GetProjectsByName(ctx context.Context, name string) ([]*model.Project, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var project persist.Project
	dpprojects, err := project.GetProjectsByName(name)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var mprojects []*model.Project
	for _, dp := range dpprojects {
		var musers []*model.User
		for _, du := range dp.Users {
			mu := model.User{ID: du.ID, Name: du.Username}
			musers = append(musers, &mu)
		}
		mp := model.Project{ID: dp.ID,
			Name:     dp.Name,
			Quantity: &dp.Quantity,
			Buget:    &dp.Buget,
			Users:    musers}
		mprojects = append(mprojects, &mp)
	}
	return mprojects, nil
}
