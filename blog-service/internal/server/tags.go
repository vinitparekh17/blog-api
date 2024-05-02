package server

import (
	"context"

	"github.com/jay-bhogayata/blogapi/internal/database"
	"github.com/jay-bhogayata/blogapi/internal/helper"
	"github.com/jay-bhogayata/blogapi/internal/pb/blogservice"
	"github.com/jay-bhogayata/blogapi/internal/store"
)

func (srv *Server) CreateBlogTag(ctx context.Context, blogTag *blogservice.BlogTag) (*blogservice.BlogTag, error) {

	tag, err := store.BlogStore.Query.CreateTag(ctx, blogTag.Name)

	if err != nil {
		store.BlogStore.Logger.Error("error while creating Tag", "error", err.Error())
		return nil, err
	}

	return &blogservice.BlogTag{
		Id:   tag.ID,
		Name: tag.Name,
	}, nil
}

func (srv *Server) ListBlogTag(ctx context.Context, _ *blogservice.Void) (*blogservice.ListBlogTagResponse, error) {
	Tags, err := store.BlogStore.Query.GetTags(ctx)
	if err != nil {
		store.BlogStore.Logger.Error("error while fetching tags", "error", err.Error())
		return nil, err
	}

	var blogList []*blogservice.BlogTag
	for _, tag := range Tags {
		blogList = append(blogList, &blogservice.BlogTag{
			Id:   tag.ID,
			Name: tag.Name,
		})
	}
	return &blogservice.ListBlogTagResponse{
		BlogTag: blogList,
	}, nil
}

func (srv *Server) UpdateBlogTag(ctx context.Context, blogTag *blogservice.BlogTag) (*blogservice.BlogTag, error) {
	tag, err := helper.UnMarshalToDbType[*blogservice.BlogTag, database.UpdateTagParams](blogTag)

	if err != nil {
		store.BlogStore.Logger.Error("error while unmarshalling tag", "error", err.Error())
		return nil, err
	}

	updatedTag, err := store.BlogStore.Query.UpdateTag(ctx, tag)
	if err != nil {
		store.BlogStore.Logger.Error("error while updating Tag", "error", err.Error())
		return nil, err
	}

	return &blogservice.BlogTag{
		Id:   updatedTag.ID,
		Name: updatedTag.Name,
	}, nil
}

func (srv *Server) DeleteBlogTag(ctx context.Context, tagId *blogservice.BlogTagId) (*blogservice.BlogTagId, error) {
	if _, err := store.BlogStore.Query.GetTagById(ctx, tagId.Id); err != nil {
		if err.Error() == "no rows in result set" {
			store.BlogStore.Logger.Error("no Tag found with given id", "error", err.Error())
			return nil, err
		}
		return nil, err
	}

	if err := store.BlogStore.Query.DeleteTag(ctx, tagId.Id); err != nil {
		store.BlogStore.Logger.Error("error while deleting Tag", "error", err.Error())
		return nil, err
	}

	return &blogservice.BlogTagId{
		Id: tagId.Id,
	}, nil
}

func (srv *Server) GetBlogTag(ctx context.Context, tagId *blogservice.BlogTagId) (*blogservice.BlogTag, error) {
	Tag, err := store.BlogStore.Query.GetTagById(ctx, tagId.Id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			store.BlogStore.Logger.Error("no tag found with given id", "error", err.Error())
			return &blogservice.BlogTag{}, nil
		}
		store.BlogStore.Logger.Error("error while fetching tag", "error", err.Error())
		return nil, err
	}

	return &blogservice.BlogTag{
		Id:   Tag.ID,
		Name: Tag.Name,
	}, nil
}
