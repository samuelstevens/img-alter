package webpage

import (
	"io/ioutil"
	"os"
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
			rootDir:    ".",
			wantedPath: "./test.html",
			wantedRoot: ".",
			err:        nil,
		},
		{
			path:       "./test.go",
			rootDir:    ".",
			wantedPath: "",
			wantedRoot: "",
			err:        &FileTypeError{"./test.go"},
		},
	}

	for _, c := range cases {
		got, err := New(c.path, c.rootDir)

		if err != nil {
			if c.err == nil {
				t.Errorf("got error %s; wanted no error", err.Error())
			}

			if err.Error() != c.err.Error() {
				t.Errorf("got error %s; wanted error %s", err.Error(), c.err.Error())
			}
		} else {
			if got.absolutePath != c.wantedPath {
				t.Errorf("webpage.New(%s, %s).absolutePath == %s, want %s", c.path, c.rootDir, got.absolutePath, c.wantedPath)
			}

			if got.rootDir != c.wantedRoot {
				t.Errorf("webpage.New(%s, %s).rootDir == %s, want %s", c.path, c.rootDir, got.rootDir, c.wantedRoot)
			}
		}
	}
}

func TestCloseNoFile(t *testing.T) {
	webpage, err := New("./test.html", ".")

	if err != nil {
		t.Errorf("Constructor error: %s", err.Error())
	}

	testErr := webpage.close()

	if testErr != nil {
		t.Errorf("safe close error: %s; wanted none", testErr.Error())
	}

}

func TestCloseTmpFile(t *testing.T) {
	webpage, err := New("./test.html", ".")

	if err != nil {
		t.Errorf("Constructor error: %s", err.Error())
	}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Errorf("tmpfile error: %s", err.Error())
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name()) // clean up

	webpage.file = tmpfile

	testErr := webpage.close()

	if testErr != nil {
		t.Errorf("safe close error: %s; wanted none", testErr.Error())
	}
}

func TestCloseAlreadyClosed(t *testing.T) {
	webpage, err := New("./test.html", ".")

	if err != nil {
		t.Errorf("Constructor error: %s", err.Error())
	}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Errorf("tmpfile error: %s", err.Error())
	}
	err = tmpfile.Close()
	if err != nil {
		t.Errorf("tmpfile error: %s", err.Error())
	}
	defer os.Remove(tmpfile.Name()) // clean up

	webpage.file = tmpfile

	testErr := webpage.close()

	if testErr == nil {
		t.Errorf("wanted 'already closed error', got none")
	}
}
