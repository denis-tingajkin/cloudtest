// Copyright (c) 2019 Cisco Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"github.com/denis-tingajkin/cloudtest/pkg/config"
	"io/ioutil"
	"testing"

	"github.com/onsi/gomega"

)

func TestImport(t *testing.T) {
	g := gomega.NewWithT(t)
	testConfig := &config.CloudTestConfig{
		Imports: []string{"samples/*"},
	}
	err := performImport(testConfig)
	g.Expect(err).Should(gomega.BeNil())
	files, _ := ioutil.ReadDir(testConfig.Imports[0][:len(testConfig.Imports[0])-1])
	g.Expect(len(testConfig.Executions) == len(files)).Should(gomega.BeTrue())
}
