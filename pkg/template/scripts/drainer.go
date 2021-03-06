// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package scripts

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/pingcap-incubator/tiup/pkg/localdata"
)

// DrainerScript represent the data to generate drainer config
type DrainerScript struct {
	NodeID    string
	IP        string
	Port      uint64
	DeployDir string
	DataDir   string
	NumaNode  string
	CommitTs  int64
	Endpoints []*PDScript
}

// NewDrainerScript returns a DrainerScript with given arguments
func NewDrainerScript(nodeID, ip, deployDir, dataDir string) *DrainerScript {
	return &DrainerScript{
		NodeID:    nodeID,
		IP:        ip,
		Port:      8249,
		DeployDir: deployDir,
		DataDir:   dataDir,
		CommitTs:  -1,
	}
}

// WithPort set Port field of DrainerScript
func (c *DrainerScript) WithPort(port uint64) *DrainerScript {
	c.Port = port
	return c
}

// WithNumaNode set NumaNode field of DrainerScript
func (c *DrainerScript) WithNumaNode(numa string) *DrainerScript {
	c.NumaNode = numa
	return c
}

// WithCommitTs set CommitTs field of DrainerScript
func (c *DrainerScript) WithCommitTs(ts int64) *DrainerScript {
	c.CommitTs = ts
	return c
}

// AppendEndpoints add new DrainerScript to Endpoints field
func (c *DrainerScript) AppendEndpoints(ends ...*PDScript) *DrainerScript {
	c.Endpoints = append(c.Endpoints, ends...)
	return c
}

// Config read ${localdata.EnvNameComponentInstallDir}/templates/scripts/run_drainer.sh.tpl as template
// and generate the config by ConfigWithTemplate
func (c *DrainerScript) Config() (string, error) {
	fp := path.Join(os.Getenv(localdata.EnvNameComponentInstallDir), "templates", "scripts", "run_drainer.sh.tpl")
	tpl, err := ioutil.ReadFile(fp)
	if err != nil {
		return "", err
	}
	return c.ConfigWithTemplate(string(tpl))
}

// ConfigWithTemplate generate the Drainer config content by tpl
func (c *DrainerScript) ConfigWithTemplate(tpl string) (string, error) {
	tmpl, err := template.New("Drainer").Parse(tpl)
	if err != nil {
		return "", err
	}

	content := bytes.NewBufferString("")
	if err := tmpl.Execute(content, c); err != nil {
		return "", err
	}

	return content.String(), nil
}
