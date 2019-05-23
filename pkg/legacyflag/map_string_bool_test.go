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

func TestMapStringBoolVar(t *testing.T) {
	cases := []struct {
		name string
		args []string

		// start
		target map[string]bool

		// expect
		set   map[string]bool
		merge map[string]bool
		apply bool
	}{
		{
			name: "flag is set",
			args: []string{"--foo=one=false,bar=true"},
			target: map[string]bool{
				"foo": true,
				"bar": false,
			},
			set: map[string]bool{
				"one": false,
				"bar": true,
			},
			merge: map[string]bool{
				"one": false,
				"foo": true,
				"bar": true,
			},
			apply: true,
		},
		{
			name: "flag is not set",
			args: []string{""},
			target: map[string]bool{
				"foo": true,
				"bar": false,
			},
			set: map[string]bool{
				"foo": true,
				"bar": false,
			},
			merge: map[string]bool{
				"foo": true,
				"bar": false,
			},
			apply: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fs := NewFlagSet("")
			val := fs.MapStringBoolVar("foo", c.target, "", &MapOptions{})
			if err := fs.Parse(c.args); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			setTarget := copyMapStringBool(c.target)
			val.Set(&setTarget)
			if !reflect.DeepEqual(setTarget, c.set) {
				t.Errorf("Set: got %#v but expected %#v", setTarget, c.set)
			}

			mergeTarget := copyMapStringBool(c.target)
			val.Merge(&mergeTarget)
			if !reflect.DeepEqual(mergeTarget, c.merge) {
				t.Errorf("Merge: got %#v but expected %#v", mergeTarget, c.merge)
			}

			applied := false
			val.Apply(func(value map[string]bool) {
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

func TestStringMapStringBool(t *testing.T) {
	var nilMap map[string]bool
	cases := []struct {
		name   string
		m      *mapStringBool
		expect string
	}{
		{"nil", newMapStringBool(&nilMap, &MapOptions{}), ""},
		{"empty", newMapStringBool(&map[string]bool{}, &MapOptions{}), ""},
		{"one key", newMapStringBool(&map[string]bool{"one": true}, &MapOptions{}), "one=true"},
		{"two keys", newMapStringBool(&map[string]bool{"one": true, "two": false}, &MapOptions{}), "one=true,two=false"},
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

func TestSetMapStringBool(t *testing.T) {
	var nilMap map[string]bool
	cases := []struct {
		name   string
		vals   []string
		start  *mapStringBool
		expect *mapStringBool
		err    string
	}{
		// we initialize the map with a default key that should be cleared by Set
		{"clears defaults", []string{""},
			newMapStringBool(&map[string]bool{"default": true}, &MapOptions{}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		// make sure we still allocate for "initialized" maps where Map was initially set to a nil map
		{"allocates map if currently nil", []string{""},
			&mapStringBool{initialized: true, m: &nilMap, options: func() *MapOptions {
				o := &MapOptions{}
				o.Default()
				return o
			}()},
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		// for most cases, we just reuse nilMap, which should be allocated by Set, and is reset before each test case
		{"empty", []string{""},
			newMapStringBool(&nilMap, &MapOptions{}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"one key", []string{"one=true"},
			newMapStringBool(&nilMap, &MapOptions{}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{"one": true},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"two keys", []string{"one=true,two=false"},
			newMapStringBool(&nilMap, &MapOptions{}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{"one": true, "two": false},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"two keys, malformed because DisableCommaSeparatedPairs=true", []string{"one=true,two=false"},
			newMapStringBool(&nilMap, &MapOptions{DisableCommaSeparatedPairs: true}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{},
				options: &MapOptions{
					DisableCommaSeparatedPairs: true,
					KeyValueSep:                "=",
					PairSep:                    ",",
				},
			}, `invalid value of one: true,two=false, err: strconv.ParseBool: parsing "true,two=false": invalid syntax`},
		{"two keys, DisableCommaSeparatedPairs=true", []string{"one=true", "two=false"},
			newMapStringBool(&nilMap, &MapOptions{DisableCommaSeparatedPairs: true}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{"one": true, "two": false},
				options: &MapOptions{
					DisableCommaSeparatedPairs: true,
					KeyValueSep:                "=",
					PairSep:                    ",",
				},
			}, ""},
		{"two keys, multiple Set invocations", []string{"one=true", "two=false"},
			newMapStringBool(&nilMap, &MapOptions{}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{"one": true, "two": false},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"two keys with space", []string{"one=true, two=false"},
			newMapStringBool(&nilMap, &MapOptions{}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{"one": true, "two": false},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"empty key", []string{"=true"},
			newMapStringBool(&nilMap, &MapOptions{}),
			&mapStringBool{
				initialized: true,
				m:           &map[string]bool{"": true},
				options: &MapOptions{
					KeyValueSep: "=",
					PairSep:     ",",
				},
			}, ""},
		{"missing value", []string{"one"},
			newMapStringBool(&nilMap, &MapOptions{}),
			nil,
			"malformed pair, expect string=bool"},
		{"non-boolean value", []string{"one=foo"},
			newMapStringBool(&nilMap, &MapOptions{}),
			nil,
			`invalid value of one: foo, err: strconv.ParseBool: parsing "foo": invalid syntax`},
		{"no target", []string{"one=true"},
			newMapStringBool(nil, &MapOptions{}),
			nil,
			"no target (nil pointer to map[string]bool)"},
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

func TestEmptyMapStringBool(t *testing.T) {
	var nilMap map[string]bool
	cases := []struct {
		name   string
		val    *mapStringBool
		expect bool
	}{
		{"nil", newMapStringBool(&nilMap, &MapOptions{}), true},
		{"empty", newMapStringBool(&map[string]bool{}, &MapOptions{}), true},
		{"populated", newMapStringBool(&map[string]bool{"foo": true}, &MapOptions{}), false},
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

func copyMapStringBool(m map[string]bool) map[string]bool {
	if m == nil {
		return nil
	}
	n := map[string]bool{}
	for k, v := range m {
		n[k] = v
	}
	return n
}
