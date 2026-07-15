package tmplfn

import (
	"slices"
	"strings"
	"time"

	"github.com/shelepuginivan/makewww/pkg/resource"
)

func inPath(path string, resources []resource.Resource) []resource.Resource {
	filtered := make([]resource.Resource, 0, len(resources))

	for _, res := range resources {
		if strings.HasPrefix(res.Path().Absolute(), path) {
			filtered = append(filtered, res)
		}
	}

	return filtered
}

func latest(resources []resource.Resource) resource.Resource {
	var (
		latest          resource.Resource
		latestCreatedAt time.Time
	)

	for _, res := range resources {
		meta, ok := res.(resource.WithMetadata)
		if !ok {
			continue
		}

		currentCreatedAt := meta.Metadata().CreatedAt

		if currentCreatedAt.After(latestCreatedAt) {
			latest = res
			latestCreatedAt = currentCreatedAt
		}
	}

	return latest
}

func sortLatestFirst(resources []resource.Resource) []resource.Resource {
	slices.SortFunc(resources, func(a, b resource.Resource) int {
		aMeta, ok := a.(resource.WithMetadata)
		if !ok {
			return 0
		}

		bMeta, ok := b.(resource.WithMetadata)
		if !ok {
			return 0
		}

		aTime := aMeta.Metadata().CreatedAt
		bTime := bMeta.Metadata().CreatedAt

		if aTime.Before(bTime) {
			return 1
		} else {
			return -1
		}
	})
	return resources
}

func sortByOrder(resources []resource.Resource) []resource.Resource {
	slices.SortFunc(resources, func(a, b resource.Resource) int {
		aMeta, ok := a.(resource.WithMetadata)
		if !ok {
			return 0
		}

		bMeta, ok := b.(resource.WithMetadata)
		if !ok {
			return 0
		}

		if aMeta.Metadata().Order > bMeta.Metadata().Order {
			return 1
		} else {
			return -1
		}
	})
	return resources
}

func draft(resources []resource.Resource) []resource.Resource {
	out := make([]resource.Resource, 0, len(resources))

	for _, res := range resources {
		meta, ok := res.(resource.WithMetadata)
		if ok && meta.Metadata().Draft {
			out = append(out, res)
		}
	}

	return out
}

func notDraft(resources []resource.Resource) []resource.Resource {
	out := make([]resource.Resource, 0, len(resources))

	for _, res := range resources {
		meta, ok := res.(resource.WithMetadata)
		if !ok || !meta.Metadata().Draft {
			out = append(out, res)
		}
	}

	return out
}
