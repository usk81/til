package google

import (
	"context"
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/indexing/v3"
	"google.golang.org/api/option"
)

// AuthorizeServerToServer is an example of the authorization flow between servers (no consent screen)
//   e.g. Indexing v3 API
func AuthorizeServerToServer() (err error) {
	bs, err := ioutil.ReadFile("your google credential json file path")
	if err != nil {
		return
	}
	cf, err := google.JWTConfigFromJSON(bs, indexing.IndexingScope)
	if err != nil {
		return
	}
	_, err = indexing.NewService(context.Background(), option.WithHTTPClient(cf.Client(context.Background())))
	return
}
