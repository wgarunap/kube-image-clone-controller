package imagecloner

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/wgarunap/kube-image-clone-controller/domain"
)

type cloner struct{}

func (c *cloner) Clone(ctx context.Context, sourceImage name.Reference, targetImage name.Reference) error {
	if err := c.isExistInPublic(ctx, sourceImage); err != nil {
		return err
	}

	img, err := remote.Image(sourceImage, remote.WithContext(ctx))
	if err != nil {
		return err
	}

	err = remote.Write(targetImage, img, remote.WithContext(ctx), remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return err
	}
	return nil
}

func (c *cloner) isExistInPublic(ctx context.Context, sourceImage name.Reference) error {
	_, err := remote.Head(sourceImage, remote.WithContext(ctx))
	if err != nil {
		if e, ok := err.(*transport.Error); ok {
			fmt.Println(fmt.Sprintf(`error reading the source image:%s, errorcode:%d`, sourceImage.Name(), e.StatusCode))
		}
		return err
	}
	return nil
}

func (c *cloner) IsExistInClones(ctx context.Context, targetImage name.Reference) (error, bool) {
	_, err := remote.Head(targetImage, remote.WithContext(ctx), remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		if e, ok := err.(*transport.Error); ok {
			fmt.Println(fmt.Sprintf(`error reading the source image:%s, errorcode:%d`, targetImage.Name(), e.StatusCode))
		}
		return err, false
	}
	return nil, true
}

func NewCloner() domain.Cloner {
	return &cloner{}
}
