package domain

import (
	"context"
	"github.com/google/go-containerregistry/pkg/name"
)

//go:generate mockgen -destination=../mocks/image_clioner.go -package=mocks -source=./imagecloner.go

type Cloner interface {
	Clone(ctx context.Context, sourceImage name.Reference, targetImage name.Reference) error
	IsExistInClones(ctx context.Context, targetImage name.Reference) (error, bool)
}
