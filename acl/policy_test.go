package acl

import (
	"reflect"
	"strings"
	"testing"
)

func TestACLPolicy_Parse_HCL(t *testing.T) {
	inp := `
event "" {
	policy = "read"
}
event "foo" {
	policy = "write"
}
event "bar" {
	policy = "deny"
}
key "" {
	policy = "read"
}
key "foo/" {
	policy = "write"
}
key "foo/bar/" {
	policy = "read"
}
key "foo/bar/baz" {
	policy = "deny"
}
keyring = "deny"
node "" {
	policy = "read"
}
node "foo" {
	policy = "write"
}
node "bar" {
	policy = "deny"
}
operator = "deny"
service "" {
	policy = "write"
}
service "foo" {
	policy = "read"
}
query "" {
	policy = "read"
}
query "foo" {
	policy = "write"
}
query "bar" {
	policy = "deny"
}
	`
	exp := &Policy{
		Events: []*EventPolicy{
			&EventPolicy{
				Event:  "",
				Policy: PolicyRead,
			},
			&EventPolicy{
				Event:  "foo",
				Policy: PolicyWrite,
			},
			&EventPolicy{
				Event:  "bar",
				Policy: PolicyDeny,
			},
		},
		Keyring: PolicyDeny,
		Keys: []*KeyPolicy{
			&KeyPolicy{
				Prefix: "",
				Policy: PolicyRead,
			},
			&KeyPolicy{
				Prefix: "foo/",
				Policy: PolicyWrite,
			},
			&KeyPolicy{
				Prefix: "foo/bar/",
				Policy: PolicyRead,
			},
			&KeyPolicy{
				Prefix: "foo/bar/baz",
				Policy: PolicyDeny,
			},
		},
		Nodes: []*NodePolicy{
			&NodePolicy{
				Name:   "",
				Policy: PolicyRead,
			},
			&NodePolicy{
				Name:   "foo",
				Policy: PolicyWrite,
			},
			&NodePolicy{
				Name:   "bar",
				Policy: PolicyDeny,
			},
		},
		Operator: PolicyDeny,
		PreparedQueries: []*PreparedQueryPolicy{
			&PreparedQueryPolicy{
				Prefix: "",
				Policy: PolicyRead,
			},
			&PreparedQueryPolicy{
				Prefix: "foo",
				Policy: PolicyWrite,
			},
			&PreparedQueryPolicy{
				Prefix: "bar",
				Policy: PolicyDeny,
			},
		},
		Services: []*ServicePolicy{
			&ServicePolicy{
				Name:   "",
				Policy: PolicyWrite,
			},
			&ServicePolicy{
				Name:   "foo",
				Policy: PolicyRead,
			},
		},
	}

	out, err := Parse(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(out, exp) {
		t.Fatalf("bad: %#v %#v", out, exp)
	}
}

func TestACLPolicy_Parse_JSON(t *testing.T) {
	inp := `{
	"event": {
		"": {
			"policy": "read"
		},
		"foo": {
			"policy": "write"
		},
		"bar": {
			"policy": "deny"
		}
	},
	"key": {
		"": {
			"policy": "read"
		},
		"foo/": {
			"policy": "write"
		},
		"foo/bar/": {
			"policy": "read"
		},
		"foo/bar/baz": {
			"policy": "deny"
		}
	},
	"keyring": "deny",
	"node": {
		"": {
			"policy": "read"
		},
		"foo": {
			"policy": "write"
		},
		"bar": {
			"policy": "deny"
		}
	},
	"operator": "deny",
	"query": {
		"": {
			"policy": "read"
		},
		"foo": {
			"policy": "write"
		},
		"bar": {
			"policy": "deny"
		}
	},
	"service": {
		"": {
			"policy": "write"
		},
		"foo": {
			"policy": "read"
		}
	}
}`
	exp := &Policy{
		Events: []*EventPolicy{
			&EventPolicy{
				Event:  "",
				Policy: PolicyRead,
			},
			&EventPolicy{
				Event:  "foo",
				Policy: PolicyWrite,
			},
			&EventPolicy{
				Event:  "bar",
				Policy: PolicyDeny,
			},
		},
		Keyring: PolicyDeny,
		Keys: []*KeyPolicy{
			&KeyPolicy{
				Prefix: "",
				Policy: PolicyRead,
			},
			&KeyPolicy{
				Prefix: "foo/",
				Policy: PolicyWrite,
			},
			&KeyPolicy{
				Prefix: "foo/bar/",
				Policy: PolicyRead,
			},
			&KeyPolicy{
				Prefix: "foo/bar/baz",
				Policy: PolicyDeny,
			},
		},
		Nodes: []*NodePolicy{
			&NodePolicy{
				Name:   "",
				Policy: PolicyRead,
			},
			&NodePolicy{
				Name:   "foo",
				Policy: PolicyWrite,
			},
			&NodePolicy{
				Name:   "bar",
				Policy: PolicyDeny,
			},
		},
		Operator: PolicyDeny,
		PreparedQueries: []*PreparedQueryPolicy{
			&PreparedQueryPolicy{
				Prefix: "",
				Policy: PolicyRead,
			},
			&PreparedQueryPolicy{
				Prefix: "foo",
				Policy: PolicyWrite,
			},
			&PreparedQueryPolicy{
				Prefix: "bar",
				Policy: PolicyDeny,
			},
		},
		Services: []*ServicePolicy{
			&ServicePolicy{
				Name:   "",
				Policy: PolicyWrite,
			},
			&ServicePolicy{
				Name:   "foo",
				Policy: PolicyRead,
			},
		},
	}

	out, err := Parse(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(out, exp) {
		t.Fatalf("bad: %#v %#v", out, exp)
	}
}

func TestACLPolicy_Keyring_Empty(t *testing.T) {
	inp := `
keyring = ""
	`
	exp := &Policy{
		Keyring: "",
	}

	out, err := Parse(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(out, exp) {
		t.Fatalf("bad: %#v %#v", out, exp)
	}
}

func TestACLPolicy_Operator_Empty(t *testing.T) {
	inp := `
operator = ""
	`
	exp := &Policy{
		Operator: "",
	}

	out, err := Parse(inp)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if !reflect.DeepEqual(out, exp) {
		t.Fatalf("bad: %#v %#v", out, exp)
	}
}

func TestACLPolicy_Bad_Policy(t *testing.T) {
	cases := []string{
		`event "" { policy = "nope" }`,
		`key "" { policy = "nope" }`,
		`keyring = "nope"`,
		`node "" { policy = "nope" }`,
		`operator = "nope"`,
		`query "" { policy = "nope" }`,
		`service "" { policy = "nope" }`,
	}
	for _, c := range cases {
		_, err := Parse(c)
		if err == nil || !strings.Contains(err.Error(), "Invalid") {
			t.Fatalf("expected policy error, got: %#v", err)
		}
	}
}
