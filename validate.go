//
// Copyright 2015 The CoreOS Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/coreos/coreos-cloudinit/config/validate"
	ignConfig "github.com/coreos/ignition/config"
)

func main() {
	os.Exit(validateConfig())
}

func validateConfig() int {
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from input: %s\n", err)
		return 2
	}

	config := bytes.Replace(body, []byte("\r"), []byte{}, -1)

	_, report, err := ignConfig.Parse(config)
	code := 0

	switch err {
	case ignConfig.ErrCloudConfig, ignConfig.ErrEmpty:
		report, err := validate.Validate(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}

		for _, e := range report.Entries() {
			fmt.Fprintf(os.Stderr, "%v\n", e)
			code = 1
		}
	case nil:
	default:
		report.Sort()
		fmt.Fprintf(os.Stderr, "%v\n", err)
		if len(report.Entries) > 0 {
			fmt.Fprintf(os.Stderr, "%v\n", report)
		}
		code = 1
	}
	return code
}
