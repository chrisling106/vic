// Copyright 2016 VMware, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/vic/lib/spec"
	"github.com/vmware/vic/pkg/vsphere/session"
	"github.com/vmware/vic/pkg/vsphere/test/env"

	"golang.org/x/net/context"
)

// Session returns a session.Session struct
func Session(ctx context.Context, t *testing.T) *session.Session {
	config := &session.Config{
		Service:        env.URL(t),
		Insecure:       true,
		Keepalive:      time.Duration(5) * time.Minute,
		DatacenterPath: "",
		DatastorePath:  "/ha-datacenter/datastore/*",
		HostPath:       "/ha-datacenter/host/*/*",
		PoolPath:       "/ha-datacenter/host/*/Resources",
	}

	s, err := session.NewSession(config).Create(ctx)
	if err != nil {
		// FIXME: See session_test.go TestSession for detail. We never get to PickRandomHost in the case of multiple hosts
		if strings.Contains(err.Error(), "resolves to multiple hosts") {
			t.SkipNow()
		} else {
			t.Errorf("ERROR: %s", err)
			t.SkipNow()
		}
	}
	return s
}

// SpecConfig returns a spec.VirtualMachineConfigSpecConfig struct
func SpecConfig(session *session.Session, name string) *spec.VirtualMachineConfigSpecConfig {
	return &spec.VirtualMachineConfigSpecConfig{
		NumCPUs:       2,
		MemoryMB:      2048,
		VMForkEnabled: true,

		ConnectorURI: "tcp://1.2.3.4:9876",

		ID:            name,
		Name:          "zombie_attack",
		BootMediaPath: session.Datastore.Path("brainz.iso"),
		VMPathName:    fmt.Sprintf("[%s]", session.Datastore.Name()),
	}
}

// PickRandomHost returns a random object.HostSystem from the hosts attached to the datastore and also lives in the same cluster
func PickRandomHost(ctx context.Context, session *session.Session, t *testing.T) *object.HostSystem {
	hosts, err := session.Datastore.AttachedClusterHosts(ctx, session.Cluster)
	if err != nil {
		t.Errorf("ERROR: %s", err)
		t.SkipNow()
	}
	return hosts[rand.Intn(len(hosts))]
}
