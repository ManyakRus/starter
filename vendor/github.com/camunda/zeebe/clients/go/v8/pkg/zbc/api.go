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
//

package zbc

import (
	"github.com/camunda/zeebe/clients/go/v8/pkg/commands"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
)

type Client interface {
	NewTopologyCommand() *commands.TopologyCommand
	// Deprecated: Use NewDeployResourceCommand instead. To be removed in 8.1.0.
	NewDeployProcessCommand() *commands.DeployCommand
	NewDeployResourceCommand() *commands.DeployResourceCommand

	NewCreateInstanceCommand() commands.CreateInstanceCommandStep1
	NewCancelInstanceCommand() commands.CancelInstanceStep1
	NewSetVariablesCommand() commands.SetVariablesCommandStep1
	NewResolveIncidentCommand() commands.ResolveIncidentCommandStep1

	NewEvaluateDecisionCommand() commands.EvaluateDecisionCommandStep1

	NewPublishMessageCommand() commands.PublishMessageCommandStep1

	NewBroadcastSignalCommand() commands.BroadcastSignalCommandStep1

	NewActivateJobsCommand() commands.ActivateJobsCommandStep1
	NewCompleteJobCommand() commands.CompleteJobCommandStep1
	NewFailJobCommand() commands.FailJobCommandStep1
	NewUpdateJobRetriesCommand() commands.UpdateJobRetriesCommandStep1
	NewThrowErrorCommand() commands.ThrowErrorCommandStep1

	NewJobWorker() worker.JobWorkerBuilderStep1

	Close() error
}
