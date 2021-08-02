package imagecloner

import (
	"context"
	"github.com/google/go-containerregistry/pkg/name"
	"testing"
)

func TestCloner_Clone(t *testing.T) {
	source, err := name.ParseReference("testingnew123/test:latest")
	if err != nil {
		t.Error(err)
		return
	}

	target, err := name.ParseReference("testingnew123/test-123:latest")
	if err != nil {
		t.Error(err)
		return
	}

	err = NewCloner().Clone(context.Background(), source, target)
	if err != nil {
		t.Error(err)
		return
	}
}


func TestNew2(t *testing.T) {
	source, err := name.ParseReference(`index.docker.io/wgarunap/generic-kafka-event-producer:latest`)
	if err != nil {
		t.Fatal(err)
	}

	target, err := name.ParseReference(`index.docker.io/testingnew123/wgarunap_generic-kafka-event-producer:latest`)
	if err != nil {
		t.Fatal(err)
	}

	err = NewCloner().Clone(context.Background(), source, target)
	if err != nil {
		t.Error(err)
		return
	}

}


//img, err := remote.Image(source, remote.WithContext(context.Background()), remote.WithAuthFromKeychain(authn.DefaultKeychain))
//if err != nil {
//	t.Error(err)
//}
//_ = img

//target, err := name.ParseReference("testingnew123/test:latest", name.WithDefaultRegistry(name.DefaultRegistry))
//if err != nil {
//	t.Error(err)
//}
//fmt.Println(source.Name())
//
//err = remote.Write(target, img, remote.WithAuth(authn.FromConfig(authn.AuthConfig{
//	Username:      "testingnew123",
//	Password:      "123123123",
//})))
//if err != nil {
//	t.Error(err)
//}
