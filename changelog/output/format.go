package output

import (
	"encoding/json"
	"github.com/coreos/go-semver/semver"
)

type Entry struct {
	Tag         string
	Name        string
	Component   string
	Description string
}

type Component struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Release struct {
	Tag        string       `json:"tag"`
	Components *[]Component `json:"components"`
}

type Changelog struct {
	Releases *[]Release `json:"releases"`
}

func (c *Changelog) AddEntry(entry Entry) {
	if c.Releases == nil {
		c.Releases = &[]Release{}
	}

	// create a new release and insert it into the changelog in the right order
	if len(*c.Releases) == 0 {
		// no releaes, create it anyways
		*c.Releases = append(*c.Releases, Release{
			Tag:        entry.Tag,
			Components: &[]Component{}})
	} else {
		// find the right place to insert the release
		for _, r := range *c.Releases {
			if r.Tag == entry.Tag {
				// release already exists, break here
				break
			}
			v0 := semver.New(r.Tag)
			v1 := semver.New(entry.Tag)
			if v0.LessThan(*v1) {
				// add release before the current release
				*c.Releases = append(*c.Releases, Release{
					Tag:        entry.Tag,
					Components: &[]Component{}})
			}
			break
		}
	}

	// find the right release and add the component
	for i, release := range *c.Releases {
		if release.Tag == entry.Tag {
			*(*c.Releases)[i].Components = append(*(*c.Releases)[i].Components, Component{
				Name:        entry.Component,
				Description: entry.Description,
			})
			break
		}
	}

}

func (c *Changelog) RenderJSON() (output []byte, err error) {
	output, err = json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}
