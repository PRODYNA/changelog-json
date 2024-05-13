package output

import (
	"fmt"
	"testing"
)

// test method
func TestChangelog_AddEntry(t *testing.T) {
	changelog := &Changelog{}
	changelog.AddEntry(Entry{
		Tag:         "1.0.0",
		Name:        "Initial Release",
		Component:   "frontend",
		Description: "Initial frontend version",
	})
	changelog.AddEntry(Entry{
		Tag:         "1.0.0",
		Name:        "Initial Release",
		Component:   "backend",
		Description: "Initial backend version",
	})
	changelog.AddEntry(Entry{
		Tag:         "1.1.0",
		Name:        "Feature 1",
		Component:   "frontend",
		Description: "This cool frontend feature",
	})
	changelog.AddEntry(Entry{
		Tag:         "1.1.0",
		Name:        "Feature 1",
		Component:   "backend",
		Description: "This cool backend feature",
	})
	changelog.AddEntry(Entry{
		Tag:         "1.2.0",
		Name:        "Feature 2",
		Component:   "frontend",
		Description: "More cool frontend feature",
	})

	fmt.Printf("%+v\n", changelog)
	output, err := changelog.RenderJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	fmt.Printf("%s", output)

	if len(*changelog.Releases) != 3 {
		t.Errorf("expected 3 releases, got %d", len(*changelog.Releases))
	}

	if len((*changelog.Releases)[0].Components) != 1 {
		t.Errorf("expected 2 components, got %d", len((*changelog.Releases)[0].Components))
	}

	if len((*changelog.Releases)[1].Components) != 2 {
		t.Errorf("expected 2 components, got %d", len((*changelog.Releases)[1].Components))
	}

	if len((*changelog.Releases)[2].Components) != 2 {
		t.Errorf("expected 1 component, got %d", len((*changelog.Releases)[2].Components))
	}

	// expect that the first release is 1.2.0
	if (*changelog.Releases)[0].Tag != "1.2.0" {
		t.Errorf("expected first release to be 1.2.0, got %s", (*changelog.Releases)[0].Tag)
	}

	// expect that the second release is 1.1.0
	if (*changelog.Releases)[1].Tag != "1.1.0" {
		t.Errorf("expected second release to be 1.1.0, got %s", (*changelog.Releases)[1].Tag)
	}

	// expect that the third release is 1.0.0
	if (*changelog.Releases)[2].Tag != "1.0.0" {
		t.Errorf("expected third release to be 1.0.0, got %s", (*changelog.Releases)[2].Tag)
	}
}
