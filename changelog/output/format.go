package output

import (
	"encoding/json"
	"github.com/Masterminds/semver/v3"
	"log/slog"
	"slices"
)

type Entry struct {
	Tag         string
	Name        string
	Component   string
	Description string
	Date        string
}

type Component struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

type Release struct {
	Tag        string       `json:"tag"`
	Components *[]Component `json:"components"`
}

type Changelog struct {
	Releases *[]Release `json:"releases"`
}

func (c *Changelog) AddEntry(entry Entry) error {
	if c.Releases == nil {
		c.Releases = &[]Release{}
	}

	// create a new release and insert it into the changelog in the right order
	if len(*c.Releases) == 0 {
		// no releaes, create it anyways
		slog.Info("Inserting first release", "entry", entry)
		*c.Releases = append(*c.Releases, Release{
			Tag:        entry.Tag,
			Components: &[]Component{}})
	} else {

		// check if it already exists
		exists := false
		for _, r := range *c.Releases {
			if r.Tag == entry.Tag {
				slog.Info("Release already exists", "tag", entry.Tag)
				exists = true
				break
			}
		}

		if !exists {
			// find the right place to insert the release
			added := false
			for i, r := range *c.Releases {
				v0, err := semver.NewVersion(r.Tag)
				if err != nil {
					slog.Error("unable to parse version", "tag", r.Tag, "error", err)
					return err
				}
				v1, err := semver.NewVersion(entry.Tag)
				if err != nil {
					slog.Error("unable to parse version", "tag", entry.Tag, "error", err)
					return err
				}
				if v0.Equal(v1) {
					slog.Info("Release already exists", "tag", entry.Tag)
					break
				}
				if v0.LessThan(v1) {
					slog.Info("Inserting release", "tag", entry.Tag, "before", r.Tag)
					// add release before the current release
					*c.Releases = slices.Insert(*c.Releases, i, Release{
						Tag:        entry.Tag,
						Components: &[]Component{},
					})
					added = true
					break
				}
			}

			if !added {
				// add release at the end
				slog.Info("Adding release at the end", "tag", entry.Tag)
				*c.Releases = append(*c.Releases, Release{
					Tag:        entry.Tag,
					Components: &[]Component{},
				})
			}
		}
	}

	// find the right release and add the component
	for _, release := range *c.Releases {
		if release.Tag == entry.Tag {
			*release.Components = append(*release.Components, Component{
				Name:        entry.Component,
				Title:       entry.Name,
				Description: entry.Description,
			})
			break
		}
	}
	return nil
}

func (c *Changelog) RenderJSON() (output []byte, err error) {
	output, err = json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}
