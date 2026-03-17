package graph

import "github.com/thkx/graphql-sqlc-go-dome/internal/usecase"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	u usecase.Usecase
}

func NewResolver(u usecase.Usecase) *Resolver {
	return &Resolver{u: u}
}
