package webpage

import (
	"path/filepath"
	"testing"
)

func TestNewWebPage(t *testing.T) {
	cases := []struct {
		path       string
		rootDir    string
		wantedPath string
		wantedRoot string
		err        error
	}{
		{
			path:       "./test.html",
			wantedPath: "./test.html",
			err:        nil,
		},
		{
			path:       "./test.go",
			wantedPath: "",
			err:        &FileTypeError{"./test.go"},
		},
	}

	for _, c := range cases {
		got, err := New(c.path)

		wantedPath, _ := filepath.Abs(c.wantedPath)

		if err != nil {
			if c.err == nil {
				t.Errorf("got error %s; wanted no error", err.Error())
			}

			if err.Error() != c.err.Error() {
				t.Errorf("got error %s; wanted error %s", err.Error(), c.err.Error())
			}
		} else {
			if got.absolutePath != wantedPath {
				t.Errorf("webpage.New(%s, %s).absolutePath == %s, want %s", c.path, c.rootDir, got.absolutePath, wantedPath)
			}
		}
	}
}

func TestLabelImages(t *testing.T) {
	var labelFunc func(string, string) string = func(imgPath string, prevDescription string) string {
		return "hello, world!"
	}

	preHTML := "<html><head></head><body>"
	postHTML := "</body></html>"

	cases := []struct {
		html string
		want string
	}{
		{
			html: "<p>Hello!</p>",
			want: "<p>Hello!</p>",
		},
		{
			html: "<img src=\"hello.png\" />",
			want: "<img alt=\"hello, world!\" src=\"hello.png\"/>",
		},
		{
			html: "<img src=\"hello.png\" alt=\"hello, world!\"/>",
			want: "<img alt=\"hello, world!\" src=\"hello.png\"/>",
		},
		{
			html: "<script>var x = \"hello\"</script>",
			want: "<script>var x = \"hello\"</script>",
		},
	}

	for _, c := range cases {
		html := preHTML + c.html + postHTML
		want := preHTML + c.want + postHTML

		got, err := LabelImages(html, labelFunc)

		if err != nil {

			t.Errorf("got error %s; wanted no error", err.Error())

		}

		if got != want {
			t.Errorf("LabelImages(%s) == %s, want %s", html, got, want)
		}

	}
}
