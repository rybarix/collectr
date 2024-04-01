package collectr_test

import (
	"bytes"
	"testing"

	"github.com/rybarix/email-collectr/internal/collectr"
)

func setup(conf collectr.CollectrConf) (*collectr.Collectr, *bytes.Buffer) {
	var buff bytes.Buffer

	coll := &collectr.Collectr{
		Conf:   conf,
		Writer: &buff,
	}

	return coll, &buff
}

func TestValidateRegexp(t *testing.T) {
	coll, wbuf := setup(collectr.CollectrConf{
		File: "",
		Fields: map[string]string{
			"email": "regexp:.*@.*",
		},
	})

	err := coll.Append(map[string]any{
		"email": "hello@world.com",
	})

	if err != nil {
		t.Fatalf("incorrect validation %s", err)
	}

	expected := "{\"email\":\"hello@world.com\"}\n"
	found := wbuf.String()

	if expected != found {
		t.Fatalf("expected %s but found %s", expected, found)
	}
}

func TestValidateText(t *testing.T) {
	coll, _ := setup(collectr.CollectrConf{
		File: "",
		Fields: map[string]string{
			// "email": "text|required|nonempty|regexp:.*@.*",
			"name": "text",
		},
	})

	err := coll.Append(map[string]any{
		"name": 2131,
	})

	if err == nil {
		t.Fatal("expected validation error")
	}

	expected := "invalid type of " + "name" + ", should have text type"
	found := err.Error()

	if expected != found {
		t.Fatalf("expected %s but found %s", expected, found)
	}
}

func TestValidateNonEmptyText(t *testing.T) {
	coll, _ := setup(collectr.CollectrConf{
		File: "",
		Fields: map[string]string{
			"name": "text|nonempty",
		},
	})

	err := coll.Append(map[string]any{
		"name": "",
	})

	if err == nil {
		t.Fatal("expected validation error")
	}

	expected := "name" + " must be nonempty string"
	found := err.Error()

	if expected != found {
		t.Fatalf("expected %s but found %s", expected, found)
	}
}
