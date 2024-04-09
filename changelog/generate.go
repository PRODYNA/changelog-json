package changelog

type ChangelogGenerator struct {
	githubToken string
}

func New(githubToken string) (*ChangelogGenerator, error) {
	return &ChangelogGenerator{
		githubToken: githubToken,
	}, nil
}

func (clg *ChangelogGenerator) Generate() (changelog *Changelog, err error) {
	return &Changelog{
		Releases: []Release{
			{
				Tag:   "1.1.0",
				Date:  "2020-01-02",
				Title: "Added new feature",
			},
			{
				Tag:   "1.0.0",
				Date:  "2020-01-01",
				Title: "Initial release",
			},
		},
	}, nil
}
