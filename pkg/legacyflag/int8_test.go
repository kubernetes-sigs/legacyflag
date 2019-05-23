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

// This file is generated. DO NOT EDIT.

package legacyflag

import(
	"testing"
	"reflect"
	
)

func TestInt8Var(t *testing.T) {
	cases := []struct {
		name string
		args []string
		set   int8
		apply bool
	}{
		{
			name: "flag is set",
			args: []string{"--foo=-1"},
			set: -1,
			apply: true,
		},
		{
			name: "flag is not set",
			args: []string{""},
			apply: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var target int8

			fs := NewFlagSet("")
			val := fs.Int8Var("foo", target, "")
			if err := fs.Parse(c.args); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			val.Set(&target)
			if !reflect.DeepEqual(target, c.set) {
				t.Errorf("Set: got %#v but expected %#v", target, c.set)
			}

			applied := false
			val.Apply(func(value int8) {
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
