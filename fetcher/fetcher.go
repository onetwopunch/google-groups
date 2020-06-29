package fetcher

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
	"io/ioutil"
)

type GroupFetcher struct {
	Config       *jwt.Config
	AdminService *admin.Service
	Context      context.Context
}

func NewDefaultGroupFetcher(keyFile, impersonate string) (*GroupFetcher, error) {
	ctx := context.Background()
	f, err := NewGroupFetcher(keyFile, impersonate)
	if err != nil {
		return nil, err
	}
	err = f.SetContext(ctx)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *GroupFetcher) SetContext(ctx context.Context) error {
	if svc, err := admin.NewService(ctx, option.WithHTTPClient(f.Config.Client(ctx))); err != nil {
		return err
	} else {
		f.AdminService = svc
		return nil
	}
}

func NewGroupFetcher(keyFile, impersonate string) (*GroupFetcher, error) {
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(keyJson, admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to read JSON service account %s", err)
	}
	config.Subject = impersonate

	return &GroupFetcher{
		Config: config,
	}, nil
}

func (f *GroupFetcher) Search(visited map[string]bool, subject string, depth int) error {
	newGroups := []string{}
	call := f.AdminService.Groups.List().UserKey(subject).Fields("nextPageToken", "groups(email)")
	if err := call.Pages(f.Context, func(groups *admin.Groups) error {
		for _, group := range groups.Groups {
			if _, ok := visited[group.Email]; ok {
				continue
			}
			visited[group.Email] = true
			newGroups = append(newGroups, group.Email)
		}
		// Only recursively search for new groups that we haven't seen
		if depth > 0 {
			for _, email := range newGroups {
				if err := f.Search(visited, email, depth-1); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (f *GroupFetcher) ListGoogleGroups(subject string, depth int) ([]string, error) {
	userGroupsMap := make(map[string]bool)

	if err := f.Search(userGroupsMap, subject, depth); err != nil {
		return nil, err
	}
	// Convert set to list
	var userGroups = make([]string, 0, len(userGroupsMap))
	for email := range userGroupsMap {
		userGroups = append(userGroups, email)
	}
	return userGroups, nil
}
