package cli

import "testing"

const (
	something = "something"
	other     = "other"
)

func TestBetterConfigValue(t *testing.T) {

	cases := []struct {
		config string
		flag   string
		want   string
	}{
		{
			config: "",
			flag:   "",
			want:   "",
		},
		{
			config: "",
			flag:   something,
			want:   something,
		},
		{
			config: other,
			flag:   "",
			want:   other,
		},
		{
			config: other,
			flag:   something,
			want:   something,
		},
	}

	for _, c := range cases {
		got := betterConfigString(c.config, c.flag)

		if got != c.want {
			t.Errorf("better(%s, %s) = %s; wanted %s", c.config, c.flag, got, c.want)
		}
	}
}
