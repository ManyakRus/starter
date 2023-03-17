// Copyright © 2018 Camunda Services GmbH (info@camunda.com)
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

package utils

import (
	"fmt"
	"github.com/camunda/zeebe/clients/go/v8/pkg/pb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/proto"
	"time"
)

const DefaultTestTimeout = 5 * time.Second
const DefaultTestTimeoutInMs = int64(DefaultTestTimeout / time.Millisecond)
const DefaultContainerWaitTimeout = 2 * time.Minute

// RPCTestMsg implements the gomock.Matcher interface
type RPCTestMsg struct {
	Msg proto.Message
}

func (r *RPCTestMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}

	// the long-polling timeout is not controllable so we can only assert it's not higher than expected
	{
		gotActivReq, okGot := msg.(*pb.ActivateJobsRequest)
		wantActivReq, okWant := r.Msg.(*pb.ActivateJobsRequest)
		if okGot && okWant {
			return cmp.Equal(gotActivReq, r.Msg, cmpopts.IgnoreUnexported(pb.ActivateJobsRequest{}),
				cmpopts.IgnoreFields(pb.ActivateJobsRequest{}, "RequestTimeout")) &&
				gotActivReq.RequestTimeout <= wantActivReq.RequestTimeout
		}
	}
	{
		gotCreateReq, okGot := msg.(*pb.CreateProcessInstanceWithResultRequest)
		wantCreateReq, okWant := r.Msg.(*pb.CreateProcessInstanceWithResultRequest)
		if okGot && okWant {
			return cmp.Equal(gotCreateReq, r.Msg, cmpopts.IgnoreUnexported(pb.CreateProcessInstanceWithResultRequest{}),
				cmpopts.IgnoreUnexported(pb.CreateProcessInstanceRequest{}),
				cmpopts.IgnoreFields(pb.CreateProcessInstanceWithResultRequest{}, "RequestTimeout")) &&
				gotCreateReq.RequestTimeout <= wantCreateReq.RequestTimeout
		}

	}

	return proto.Equal(m, r.Msg)
}

func (r *RPCTestMsg) String() string {
	return fmt.Sprintf("is %s", r.Msg)
}
