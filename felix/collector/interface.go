// Copyright (c) 2018-2025 Tigera, Inc. All rights reserved.
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

package collector

import (
	"github.com/projectcalico/calico/felix/collector/types"
	"github.com/projectcalico/calico/felix/proto"
)

type Collector interface {
	Start() error
	ReportingChannel() chan<- *proto.DataplaneStats
	RegisterMetricsReporter(types.Reporter)
	SetDataplaneInfoReader(types.DataplaneInfoReader)
	SetPacketInfoReader(types.PacketInfoReader)
	SetConntrackInfoReader(types.ConntrackInfoReader)
}
