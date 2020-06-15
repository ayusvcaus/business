package graph

import (
	"context"
	"fmt"

	"github.com/alecthomas/log4go"

	"github.com/ayusvcaus/business/auth"
	"github.com/ayusvcaus/business/graph/model"
	"github.com/ayusvcaus/business/persist"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	link := persist.Link{Address: input.Address,
		Title:  input.Title,
		UserID: input.UserID}
	plink, err := link.Create()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	mlink := model.Link{ID: plink.ID,
		Title:   plink.Title,
		Address: &plink.Address,
		UserID:  &plink.UserID}
	return &mlink, nil
}

func (r *mutationResolver) UpdateLink(ctx context.Context, linkID int, input model.NewLink) (*model.Link, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	link := persist.Link{ID: linkID, Title: input.Title, Address: input.Address, UserID: input.UserID}
	plink, err := link.Update()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	mlink := model.Link{ID: plink.ID, Title: plink.Title, Address: &plink.Address, UserID: &plink.UserID}
	return &mlink, nil
}

func (r *mutationResolver) DeleteLink(ctx context.Context, linkID int) (*model.Link, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	link := persist.Link{ID: linkID}
	err := link.Delete()
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	mlink := model.Link{ID: link.ID, Title: link.Title, Address: &link.Address, UserID: &link.UserID}
	return &mlink, nil
}

func (r *mutationResolver) SetLink2User(ctx context.Context, linkID int, username string) (string, error) {
	panic(fmt.Errorf("not implemented"))

}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var dblink persist.Link
	dblinks, _ := dblink.Links()
	var mlinks []*model.Link
	for _, dl := range dblinks {
		ml := model.Link{ID: dl.ID, Title: dl.Title, Address: &dl.Address, UserID: &dl.UserID}
		mlinks = append(mlinks, &ml)
	}
	return mlinks, nil
}

func (r *queryResolver) GetLinksByTitle(ctx context.Context, title string) ([]*model.Link, error) {
	cuser := auth.ForContext(ctx)
	if cuser == nil {
		errStr := "Not authenticated"
		log4go.Error(errStr)
		return nil, &persist.UserError{Err: errStr}
	}
	var link persist.Link
	plinks, err := link.GetLinksByTitle(title)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var mlinks []*model.Link
	for _, dl := range plinks {
		ml := model.Link{ID: dl.ID, Title: dl.Title, Address: &dl.Address}
		mlinks = append(mlinks, &ml)
	}
	return mlinks, nil
}
