package fetcher

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
	"io/ioutil"
)

type GroupFetcher struct {
	adminService    *admin.Service
	ctx             context.Context
	subject         string
	maxRecurseDepth int
	groups          []string
}

func NewGroupFetcher(keyFile, impersonate, subject string, depth int) (*GroupFetcher, error) {
	keyJson, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(keyJson, admin.AdminDirectoryGroupReadonlyScope, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Unable to read JSON service account %s", err)
	}
	config.Subject = impersonate
	ctx := context.Background()
	svc, err := admin.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))

	return &GroupFetcher{
		adminService:    svc,
		ctx:             ctx,
		maxRecurseDepth: depth,
		subject:         subject,
	}, nil
}

func (f *GroupFetcher) Search(visited map[string]bool, subject string, depth int) error {
	newGroups := []string{}
	call := f.adminService.Groups.List().UserKey(subject).Fields("nextPageToken", "groups(email)")
	if err := call.Pages(f.ctx, func(groups *admin.Groups) error {
		for _, group := range groups.Groups {
			if _, ok := visited[group.Email]; ok {
				continue
			}
			visited[group.Email] = true
			newGroups = append(newGroups, group.Email)
		}
		// Only recursively search for new groups that we haven't seen
		for _, email := range newGroups {
			if err := f.Search(visited, email, depth-1); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (f *GroupFetcher) ListGoogleGroups() ([]string, error) {
	userGroupsMap := make(map[string]bool)

	if err := f.Search(userGroupsMap, f.subject, f.maxRecurseDepth); err != nil {
		return nil, err
	}
	// Convert set to list
	var userGroups = make([]string, 0, len(userGroupsMap))
	for email := range userGroupsMap {
		userGroups = append(userGroups, email)
	}
	return userGroups, nil
}
