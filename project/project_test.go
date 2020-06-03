package project

import (
	"testing"

	"github.com/samuelstevens/goimglabeler/api"
)

func TestNewProject(t *testing.T) {
	client := api.New()

	cases := []struct {
		dir  string
		want Project
	}{
		{dir: ".", want: Project{dir: ".", client: client}},
	}

	for _, c := range cases {
		got := New(c.dir, client)

		if *got != c.want {
			t.Errorf("NewProject(%s) == %v, want %v", c.dir, got, c.want)
		}
	}
}
