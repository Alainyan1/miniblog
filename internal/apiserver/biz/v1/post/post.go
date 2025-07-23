// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package post

import (
	"context"
	"miniblog/internal/apiserver/model"
	"miniblog/internal/apiserver/pkg/conversion"
	"miniblog/internal/apiserver/store"
	"miniblog/internal/pkg/contextx"

	apiv1 "miniblog/pkg/api/apiserver/v1"

	"github.com/jinzhu/copier"
	"github.com/onexstack/onexstack/pkg/store/where"
)

type PostBiz interface {
	Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error)
	Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error)
	List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error)

	PostExpansion
}

type PostExpansion interface{}

type postBiz struct {
	store store.IStore
}

var _ PostBiz = (*postBiz)(nil)

func New(store store.IStore) *postBiz {
	return &postBiz{store: store}
}

func (b *postBiz) Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	var postM model.PostM
	_ = copier.Copy(&postM, rq)

	postM.UserID = contextx.UserID(ctx)

	if err := b.store.Post().Create(ctx, &postM); err != nil {
		return nil, err
	}

	return &apiv1.CreatePostResponse{PostID: postM.PostID}, nil
}

func (b *postBiz) Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	// 1. 构建查询条件
	whr := where.T(ctx).F("postID", rq.GetPostID())

	// 2. 调用store层的postModel的Get方法, 传入查询条件获取对应的postM结构体
	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	if rq.Title != nil {
		postM.Title = rq.GetTitle()
	}

	if rq.Content != nil {
		postM.Content = rq.GetContent()
	}

	// 3. 调用store层的postModel的Update方法更新postM结构体
	if err := b.store.Post().Update(ctx, postM); err != nil {
		return nil, err
	}

	return &apiv1.UpdatePostResponse{}, nil
}

func (b *postBiz) Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostIDs())

	if err := b.store.Post().Delete(ctx, whr); err != nil {
		return nil, err
	}

	return &apiv1.DeletePostResponse{}, nil
}

func (b *postBiz) Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	whr := where.T(ctx).F("postID", rq.GetPostID())

	postM, err := b.store.Post().Get(ctx, whr)
	if err != nil {
		return nil, err
	}

	return &apiv1.GetPostResponse{Post: conversion.PostModelToPostV1(postM)}, nil
}

func (b *postBiz) List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	whr := where.T(ctx).P(int(rq.GetOffset()), int(rq.GetLimit()))

	count, postList, err := b.store.Post().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	posts := make([]*apiv1.Post, 0, len(postList))
	for _, post := range postList {
		converted := conversion.PostModelToPostV1(post)
		posts = append(posts, converted)
	}

	return &apiv1.ListPostResponse{TotalCount: count, Posts: posts}, nil
}
