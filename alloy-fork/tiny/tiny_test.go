package tiny

import (
	"testing"

	"github.com/grafana/alloy/syntax"
)

func TestUnmarshalAlloy(t *testing.T) {
	var args Arguments
	if err := syntax.Unmarshal([]byte(``), &args); err != nil {
		t.Fatal(err)
	}
	if args != DefaultArguments {
		t.Fatalf("got %+v, want defaults", args)
	}

	if err := syntax.Unmarshal([]byte(`loadavg_path = "/host/proc/loadavg"`), &args); err != nil {
		t.Fatal(err)
	}
	if args.LoadavgPath != "/host/proc/loadavg" {
		t.Fatalf("got %q", args.LoadavgPath)
	}
}
