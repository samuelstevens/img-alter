package util

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestMakeAbsRelativeTo(t *testing.T) {
	cases := []struct {
		root  string
		other string
		want  string
	}{
		{
			root:  "../test-dir/test.txt",
			other: "nested-dir/nested.txt",
			want:  "../test-dir/nested-dir/nested.txt",
		},
		{
			root:  "../test-dir/",
			other: "nested-dir/nested.txt",
			want:  "../test-dir/nested-dir/nested.txt",
		},
		{
			root:  "../test-dir/nested-dir",
			other: "nested-dir/nested.txt",
			want:  "../test-dir/nested-dir/nested.txt",
		},
		{
			root:  "../test-dir/nested-dir/nested.txt",
			other: "nested-dir/sub-dir/sub.txt",
			want:  "../test-dir/nested-dir/sub-dir/sub.txt",
		},
		{
			root:  "../test-dir/nested-dir/sub-dir/sub.txt",
			other: "nested-dir/other-dir/other.txt",
			want:  "../test-dir/nested-dir/other-dir/other.txt",
		},
		{
			root:  "../test-dir/nested-dir/sub-dir",
			other: "nested-dir/sub-dir/sub.txt",
			want:  "../test-dir/nested-dir/sub-dir/sub.txt",
		},
	}

	for _, c := range cases {
		root, _ := filepath.Abs(c.root)
		got, err := MakeAbsRelativeTo(root, c.other)
		want, _ := filepath.Abs(c.want)

		if err != nil {
			t.Errorf("wanted no error; got %s", err.Error())
		}

		if got != want {
			t.Errorf("wanted %s; got %s", want, got)
		}
	}

}

func TestRoot(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{
			path: "/hello/world",
			want: "/hello",
		},
		{
			path: "hello/world",
			want: "hello",
		},
		{
			path: "./hello/world",
			want: ".",
		},
		{
			path: "./hello",
			want: ".",
		},
	}

	for _, c := range cases {
		got := root(c.path)

		if got != c.want {
			t.Errorf("want %s; got %s", c.want, got)
		}
	}
}

func TestJoin(t *testing.T) {
	cases := []struct {
		base  []string
		other []string
		want  string
	}{
		{
			base:  []string{},
			other: []string{},
			want:  "/",
		},
		{
			base:  []string{"hello"},
			other: []string{},
			want:  "/hello",
		},
		{
			base:  []string{},
			other: []string{"hello"},
			want:  "/hello",
		},
		{
			base:  []string{"hello"},
			other: []string{"world"},
			want:  "/hello/world",
		},
	}

	for _, c := range cases {
		got := Join(c.base, c.other)

		if got != c.want {
			t.Errorf("Join(%s, %s): got %s, wanted %s", c.base, c.other, got, c.want)
		}
	}
}

func TestSplitPath(t *testing.T) {
	cases := []struct {
		path string
		want []string
	}{
		{
			path: "/hello/world/",
			want: []string{"hello", "world"},
		},
		{
			path: "/hello/",
			want: []string{"hello"},
		},
	}

	for _, c := range cases {
		got := splitPath(c.path)

		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("splitPath(%s) = %v (len: %d); want %v (len: %d)", c.path, got, len(got), c.want, len(c.want))
		}
	}
}

func TestGreatestCommonAbs(t *testing.T) {
	cases := []struct {
		general  string
		specific string
		want     string
	}{
		{
			general:  "/hello/world/",
			specific: "/hello/world/stuff.txt",
			want:     "/hello/world",
		},
		{
			general:  "/hello/",
			specific: "/hello/world/stuff.txt",
			want:     "/hello",
		},
	}

	for _, c := range cases {
		got, err := GreatestCommonAbs(c.general, c.specific)

		if err != nil {
			t.Errorf(err.Error())
		}

		if got != c.want {
			t.Errorf("common(%v, %v) = %s; want %s", c.general, c.specific, got, c.want)
		}
	}
}
