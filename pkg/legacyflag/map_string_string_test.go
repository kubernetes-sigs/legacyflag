/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package legacyflag

import (
	"reflect"
	"testing"
)

func TestMapStringStringVar(t *testing.T) {
	cases := []struct {
		name string
		args []string

		// start
		target map[string]string

		// expect
		set   map[string]string
		merge map[string]string
		apply bool
	}{
		{
			name: "flag is set",
			args: []string{"--foo=one=quux,bar=baz"},
			target: map[string]string{
				"foo": "baz",
				"bar": "quux",
			},
			set: map[string]string{
				"one": "quux",
				"bar": "baz",
			},
			merge: map[string]string{
				"one": "quux",
				"foo": "baz",
				"bar": "baz",
			},
			apply: true,
		},
		{
			name: "flag is not set",
			args: []string{""},
			target: map[string]string{
				"foo": "baz",
				"bar": "quux",
			},
			set: map[string]string{
				"foo": "baz",
				"bar": "quux",
			},
			merge: map[string]string{
				"foo": "baz",
				"bar": "quux",
			},
			apply: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fs := NewFlagSet("")
			val := fs.MapStringStringVar("foo", c.target, "", &MapOptions{})
			if err := fs.Parse(c.args); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			setTarget := copyMapStringString(c.target)
			val.Set(&setTarget)
			if !reflect.DeepEqual(setTarget, c.set) {
				t.Errorf("Set: got %#v but expected %#v", setTarget, c.set)
			}

			mergeTarget := copyMapStringString(c.target)
			val.Merge(&mergeTarget)
			if !reflect.DeepEqual(mergeTarget, c.merge) {
				t.Errorf("Merge: got %#v but expected %#v", mergeTarget, c.merge)
			}

			applied := false
			val.Apply(func(value map[string]string) {
				applied = true
				// value passed to apply func should match the expected result of Set
				if !reflect.DeepEqual(value, c.set) {
					t.Errorf("Apply: got %#v but expected %#v", value, c.set)
				}
			})
			if c.apply && !applied {
				t.Errorf("Apply: apply func not called")
			} else if !c.apply && applied {
				t.Errorf("Apply: apply func called, should not have been")
			}
		})
	}
}

func TestStringMapStringString(t *testing.T) {
	var nilMap map[string]string
	cases := []struct {
		name   string
		m      *mapStringString
		expect string
	}{
		{"nil", newMapStringString(&nilMap, &MapOptions{}), ""},
		{"empty", newMapStringString(&map[string]string{}, &MapOptions{}), ""},
		{"one key", newMapStringString(&map[string]string{"one": "foo"}, &MapOptions{}), "one=foo"},
		{"two keys", newMapStringString(&map[string]string{"one": "foo", "two": "bar"}, &MapOptions{}), "one=foo,two=bar"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			str := c.m.String()
			if c.expect != str {
				t.Fatalf("expect %q but got %q", c.expect, str)
			}
		})
	}
}

func TestSetMapStringString(t *testing.T) {
	var nilMap map[string]string
	cases := []struct {
		name   string
		vals   []string
		start  *mapStringString
		expect *mapStringString
		err    string
	}{
		// we initialize the map with a default key that should be cleared by Set
		{"clears defaults", []string{""},
			newMapStringString(&map[string]string{"default": ""}, &MapOptions{}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		// make sure we still allocate for "initialized" maps where m was initially set to a nil map
		{"allocates map if currently nil", []string{""},
			&mapStringString{initialized: true, m: &nilMap, options: func() *MapOptions {
				o := &MapOptions{}
				o.Default()
				return o
			}()},
			&mapStringString{
				initialized: true,
				m:           &map[string]string{},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		// for most cases, we just reuse nilMap, which should be allocated by Set, and is reset before each test case
		{"empty", []string{""},
			newMapStringString(&nilMap, &MapOptions{}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"one key", []string{"one=foo"},
			newMapStringString(&nilMap, &MapOptions{}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{"one": "foo"},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"two keys", []string{"one=foo,two=bar"},
			newMapStringString(&nilMap, &MapOptions{}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{"one": "foo", "two": "bar"},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"one key, DisableCommaSeparatedPairs=true", []string{"one=foo,bar"},
			newMapStringString(&nilMap, &MapOptions{DisableCommaSeparatedPairs: true}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{"one": "foo,bar"},
				options: &MapOptions{
					DisableCommaSeparatedPairs: true,
					KeyValueSep:                "=",
					PairSep:                    ",",
				},
			}, ""},
		{"two keys, DisableCommaSeparatedPairs=true", []string{"one=foo,bar", "two=foo,bar"},
			newMapStringString(&nilMap, &MapOptions{DisableCommaSeparatedPairs: true}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{"one": "foo,bar", "two": "foo,bar"},
				options: &MapOptions{
					DisableCommaSeparatedPairs: true,
					KeyValueSep:                "=",
					PairSep:                    ",",
				},
			}, ""},
		{"two keys, multiple Set invocations", []string{"one=foo", "two=bar"},
			newMapStringString(&nilMap, &MapOptions{}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{"one": "foo", "two": "bar"},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"two keys with space", []string{"one=foo, two=bar"},
			newMapStringString(&nilMap, &MapOptions{}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{"one": "foo", "two": "bar"},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"empty key", []string{"=foo"},
			newMapStringString(&nilMap, &MapOptions{}),
			&mapStringString{
				initialized: true,
				m:           &map[string]string{"": "foo"},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"missing value", []string{"one"},
			newMapStringString(&nilMap, &MapOptions{}),
			nil,
			"malformed pair, expect string=string"},
		{"no target", []string{"a:foo"},
			newMapStringString(nil, &MapOptions{}),
			nil,
			"no target (nil pointer to map[string]string)"},
	}
	for _, c := range cases {
		nilMap = nil
		t.Run(c.name, func(t *testing.T) {
			var err error
			for _, val := range c.vals {
				err = c.start.Set(val)
				if err != nil {
					break
				}
			}
			if c.err != "" {
				if err == nil || err.Error() != c.err {
					t.Fatalf("expect error %s but got %v", c.err, err)
				}
				return
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(c.expect, c.start) {
				t.Fatalf("expect options: %#v, map: %#v but got options: %#v, map: %#v",
					c.expect.options, c.expect.m, c.start.options, c.start.m)
			}
		})
	}
}

func TestEmptyMapStringString(t *testing.T) {
	var nilMap map[string]string
	cases := []struct {
		name   string
		val    *mapStringString
		expect bool
	}{
		{"nil", newMapStringString(&nilMap, &MapOptions{}), true},
		{"empty", newMapStringString(&map[string]string{}, &MapOptions{}), true},
		{"populated", newMapStringString(&map[string]string{"foo": ""}, &MapOptions{}), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.val.Empty()
			if result != c.expect {
				t.Fatalf("expect %t but got %t", c.expect, result)
			}
		})
	}
}

func copyMapStringString(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	n := map[string]string{}
	for k, v := range m {
		n[k] = v
	}
	return n
}
