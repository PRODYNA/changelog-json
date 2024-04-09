package changelog

import "encoding/json"

type Release struct {
	Tag   string `json:"tag"`
	Date  string `json:"date"`
	Title string `json:"title"`
}

type Changelog struct {
	Releases []Release `json:"releases"`
}

func (c *Changelog) AddRelease(release Release) {
	c.Releases = append(c.Releases, release)
}

func (c *Changelog) RenderJSON() (output []byte, err error) {
	output, err = json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return output, nil
}
