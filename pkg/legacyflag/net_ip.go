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
	"github.com/spf13/pflag"
	"net"
)

// IPValue is a reference to a registered net.IP flag value.
type IPValue struct {
	name string
	value net.IP
	fs *pflag.FlagSet
}

// IPVar registers a flag for net.IP against the FlagSet, and returns
// a IPValue reference to the registered flag value.
func (fs *FlagSet) IPVar(name string, def net.IP, usage string) *IPValue {
	v := &IPValue{
		name: name,
		fs: fs.fs,
	}
	fs.fs.IPVar(&v.value, name, def, usage)
	return v
}

// Set copies the flag value to the target if the flag was set.
func (v *IPValue) Set(target *net.IP) {
	if v.fs.Changed(v.name) {
		*target = v.value
	}
}

// Apply calls the apply func with the flag value if the flag was set.
func (v *IPValue) Apply(apply func(value net.IP)) {
	if v.fs.Changed(v.name) {
		apply(v.value)
	}
}
