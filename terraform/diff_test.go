// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	"fmt"
	"strconv"
	"testing"
)

func TestInstanceDiff_ChangeType(t *testing.T) {
	t.Parallel()

	cases := []struct {
		diff   *InstanceDiff
		Result diffChangeType
	}{
		{
			&InstanceDiff{},
			diffNone,
		},
		{
			&InstanceDiff{Destroy: true},
			diffDestroy,
		},
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old: "",
						New: "bar",
					},
				},
			},
			diffUpdate,
		},
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old:         "",
						New:         "bar",
						RequiresNew: true,
					},
				},
			},
			diffCreate,
		},
		{
			&InstanceDiff{
				Destroy: true,
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old:         "",
						New:         "bar",
						RequiresNew: true,
					},
				},
			},
			diffDestroyCreate,
		},
		{
			&InstanceDiff{
				DestroyTainted: true,
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old:         "",
						New:         "bar",
						RequiresNew: true,
					},
				},
			},
			diffDestroyCreate,
		},
	}

	for i, tc := range cases {
		actual := tc.diff.ChangeType()
		if actual != tc.Result {
			t.Fatalf("%d: %#v", i, actual)
		}
	}
}

func TestInstanceDiff_Empty(t *testing.T) {
	t.Parallel()

	var rd *InstanceDiff

	if !rd.Empty() {
		t.Fatal("should be empty")
	}

	rd = new(InstanceDiff)

	if !rd.Empty() {
		t.Fatal("should be empty")
	}

	rd = &InstanceDiff{Destroy: true}

	if rd.Empty() {
		t.Fatal("should not be empty")
	}

	rd = &InstanceDiff{
		Attributes: map[string]*ResourceAttrDiff{
			"foo": {
				New: "bar",
			},
		},
	}

	if rd.Empty() {
		t.Fatal("should not be empty")
	}
}

func TestInstanceDiff_RequiresNew(t *testing.T) {
	t.Parallel()

	rd := &InstanceDiff{
		Attributes: map[string]*ResourceAttrDiff{
			"foo": {},
		},
	}

	if rd.RequiresNew() {
		t.Fatal("should not require new")
	}

	rd.Attributes["foo"].RequiresNew = true

	if !rd.RequiresNew() {
		t.Fatal("should require new")
	}
}

func TestInstanceDiff_RequiresNew_nil(t *testing.T) {
	t.Parallel()

	var rd *InstanceDiff

	if rd.RequiresNew() {
		t.Fatal("should not require new")
	}
}

func TestInstanceDiffSame(t *testing.T) {
	t.Parallel()

	cases := []struct {
		One, Two *InstanceDiff
		Same     bool
		Reason   string
	}{
		{
			&InstanceDiff{},
			&InstanceDiff{},
			true,
			"",
		},

		{
			nil,
			nil,
			true,
			"",
		},

		{
			&InstanceDiff{Destroy: false},
			&InstanceDiff{Destroy: true},
			false,
			"diff: Destroy; old: false, new: true",
		},

		{
			&InstanceDiff{Destroy: true},
			&InstanceDiff{Destroy: true},
			true,
			"",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {},
				},
			},
			true,
			"",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"bar": {},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {},
				},
			},
			false,
			"attribute mismatch: bar",
		},

		// Extra attributes
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {},
					"bar": {},
				},
			},
			false,
			"extra attributes: bar",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {RequiresNew: true},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {RequiresNew: false},
				},
			},
			false,
			"diff RequiresNew; old: true, new: false",
		},

		// NewComputed on primitive
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old:         "",
						New:         "${var.foo}",
						NewComputed: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old: "0",
						New: "1",
					},
				},
			},
			true,
			"",
		},

		// NewComputed on primitive, removed
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old:         "",
						New:         "${var.foo}",
						NewComputed: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{},
			},
			true,
			"",
		},

		// NewComputed on set, removed
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old:         "",
						New:         "",
						NewComputed: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.1": {
						Old:        "foo",
						New:        "",
						NewRemoved: true,
					},
					"foo.2": {
						Old: "",
						New: "bar",
					},
				},
			},
			true,
			"",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {NewComputed: true},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.0": {
						Old: "",
						New: "12",
					},
				},
			},
			true,
			"",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.~35964334.bar": {
						Old: "",
						New: "${var.foo}",
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.87654323.bar": {
						Old: "",
						New: "12",
					},
				},
			},
			true,
			"",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old:         "0",
						NewComputed: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{},
			},
			true,
			"",
		},

		// Computed can change RequiresNew by removal, and that's okay
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old:         "0",
						NewComputed: true,
						RequiresNew: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{},
			},
			true,
			"",
		},

		// Computed can change Destroy by removal, and that's okay
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old:         "0",
						NewComputed: true,
						RequiresNew: true,
					},
				},

				Destroy: true,
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{},
			},
			true,
			"",
		},

		// Computed can change Destroy by elements
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old:         "0",
						NewComputed: true,
						RequiresNew: true,
					},
				},

				Destroy: true,
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "1",
						New: "1",
					},
					"foo.12": {
						Old:         "4",
						New:         "12",
						RequiresNew: true,
					},
				},

				Destroy: true,
			},
			true,
			"",
		},

		// Computed sets may not contain all fields in the original diff, and
		// because multiple entries for the same set can compute to the same
		// hash before the values are computed or interpolated, the overall
		// count can change as well.
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.~35964334.bar": {
						Old: "",
						New: "${var.foo}",
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "2",
					},
					"foo.87654323.bar": {
						Old: "",
						New: "12",
					},
					"foo.87654325.bar": {
						Old: "",
						New: "12",
					},
					"foo.87654325.baz": {
						Old: "",
						New: "12",
					},
				},
			},
			true,
			"",
		},

		// Computed values in maps will fail the "Same" check as well
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.%": {
						Old:         "",
						New:         "",
						NewComputed: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.%": {
						Old:         "0",
						New:         "1",
						NewComputed: false,
					},
					"foo.val": {
						Old: "",
						New: "something",
					},
				},
			},
			true,
			"",
		},

		// In a DESTROY/CREATE scenario, the plan diff will be run against the
		// state of the old instance, while the apply diff will be run against an
		// empty state (because the state is cleared when the destroy runs.)
		// For complex attributes, this can result in keys that seem to disappear
		// between the two diffs, when in reality everything is working just fine.
		//
		// Same() needs to take into account this scenario by analyzing NewRemoved
		// and treating as "Same" a diff that does indeed have that key removed.
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"somemap.oldkey": {
						Old:        "long ago",
						New:        "",
						NewRemoved: true,
					},
					"somemap.newkey": {
						Old: "",
						New: "brave new world",
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"somemap.newkey": {
						Old: "",
						New: "brave new world",
					},
				},
			},
			true,
			"",
		},

		// Another thing that can occur in DESTROY/CREATE scenarios is that list
		// values that are going to zero have diffs that show up at plan time but
		// are gone at apply time. The NewRemoved handling catches the fields and
		// treats them as OK, but it also needs to treat the .# field itself as
		// okay to be present in the old diff but not in the new one.
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"reqnew": {
						Old:         "old",
						New:         "new",
						RequiresNew: true,
					},
					"somemap.#": {
						Old: "1",
						New: "0",
					},
					"somemap.oldkey": {
						Old:        "long ago",
						New:        "",
						NewRemoved: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"reqnew": {
						Old:         "",
						New:         "new",
						RequiresNew: true,
					},
				},
			},
			true,
			"",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"reqnew": {
						Old:         "old",
						New:         "new",
						RequiresNew: true,
					},
					"somemap.%": {
						Old: "1",
						New: "0",
					},
					"somemap.oldkey": {
						Old:        "long ago",
						New:        "",
						NewRemoved: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"reqnew": {
						Old:         "",
						New:         "new",
						RequiresNew: true,
					},
				},
			},
			true,
			"",
		},

		// Innner computed set should allow outer change in key
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.~1.outer_val": {
						Old: "",
						New: "foo",
					},
					"foo.~1.inner.#": {
						Old: "0",
						New: "1",
					},
					"foo.~1.inner.~2.value": {
						Old:         "",
						New:         "${var.bar}",
						NewComputed: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.12.outer_val": {
						Old: "",
						New: "foo",
					},
					"foo.12.inner.#": {
						Old: "0",
						New: "1",
					},
					"foo.12.inner.42.value": {
						Old: "",
						New: "baz",
					},
				},
			},
			true,
			"",
		},

		// Innner computed list should allow outer change in key
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.~1.outer_val": {
						Old: "",
						New: "foo",
					},
					"foo.~1.inner.#": {
						Old: "0",
						New: "1",
					},
					"foo.~1.inner.0.value": {
						Old:         "",
						New:         "${var.bar}",
						NewComputed: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old: "0",
						New: "1",
					},
					"foo.12.outer_val": {
						Old: "",
						New: "foo",
					},
					"foo.12.inner.#": {
						Old: "0",
						New: "1",
					},
					"foo.12.inner.0.value": {
						Old: "",
						New: "baz",
					},
				},
			},
			true,
			"",
		},

		// When removing all collection items, the diff is allowed to contain
		// nothing when re-creating the resource. This should be the "Same"
		// since we said we were going from 1 to 0.
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.%": {
						Old:         "1",
						New:         "0",
						RequiresNew: true,
					},
					"foo.bar": {
						Old:         "baz",
						New:         "",
						NewRemoved:  true,
						RequiresNew: true,
					},
				},
			},
			&InstanceDiff{},
			true,
			"",
		},

		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo.#": {
						Old:         "1",
						New:         "0",
						RequiresNew: true,
					},
					"foo.0": {
						Old:         "baz",
						New:         "",
						NewRemoved:  true,
						RequiresNew: true,
					},
				},
			},
			&InstanceDiff{},
			true,
			"",
		},

		// Make sure that DestroyTainted diffs pass as well, especially when diff
		// two works off of no state.
		{
			&InstanceDiff{
				DestroyTainted: true,
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old: "foo",
						New: "foo",
					},
				},
			},
			&InstanceDiff{
				DestroyTainted: true,
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old: "",
						New: "foo",
					},
				},
			},
			true,
			"",
		},
		// RequiresNew in different attribute
		{
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old: "foo",
						New: "foo",
					},
					"bar": {
						Old:         "bar",
						New:         "baz",
						RequiresNew: true,
					},
				},
			},
			&InstanceDiff{
				Attributes: map[string]*ResourceAttrDiff{
					"foo": {
						Old: "",
						New: "foo",
					},
					"bar": {
						Old:         "",
						New:         "baz",
						RequiresNew: true,
					},
				},
			},
			true,
			"",
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			same, reason := tc.One.Same(tc.Two)
			if same != tc.Same {
				t.Fatalf("%d: expected same: %t, got %t (%s)\n\n one: %#v\n\ntwo: %#v",
					i, tc.Same, same, reason, tc.One, tc.Two)
			}
			if reason != tc.Reason {
				t.Fatalf(
					"%d: bad reason\n\nexpected: %#v\n\ngot: %#v", i, tc.Reason, reason)
			}
		})
	}
}

func TestCountFlatmapContainerValues(t *testing.T) {
	t.Parallel()

	for i, tc := range []struct {
		attrs map[string]string
		key   string
		count string
	}{
		{
			attrs: map[string]string{"set.2.list.#": "9999", "set.2.list.0": "x", "set.2.list.0.z": "y", "set.2.attr": "bar", "set.#": "9999"},
			key:   "set.2.list.#",
			count: "1",
		},
		{
			attrs: map[string]string{"set.2.list.#": "9999", "set.2.list.0": "x", "set.2.list.0.z": "y", "set.2.attr": "bar", "set.#": "9999"},
			key:   "set.#",
			count: "1",
		},
		{
			attrs: map[string]string{"set.2.list.0": "x", "set.2.list.0.z": "y", "set.2.attr": "bar", "set.#": "9999"},
			key:   "set.#",
			count: "1",
		},
		{
			attrs: map[string]string{"map.#": "3", "map.a": "b", "map.a.#": "0", "map.b": "4"},
			key:   "map.#",
			count: "2",
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			count := countFlatmapContainerValues(tc.key, tc.attrs)
			if count != tc.count {
				t.Fatalf("expected %q, got %q", tc.count, count)
			}
		})
	}
}
