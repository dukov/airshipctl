/*
Copyright 2014 The Kubernetes Authors.

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

package cluster

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/cluster"
	"opendev.org/airship/airshipctl/pkg/environment"
)

const (
	getInitInfraLong = (`
		Deploy base infrastructure to kubernetes cluster`)

	getInitInfraExample = (`
		# deploy infra to cluster
		airshipctl cluster initinfra`)
)

// NewCmdInitInfra creates a command to manage initial airship infrastructure
func NewCmdInitInfra(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	i := cluster.NewInfra(rootSettings)
	initFlagsCmd := &cobra.Command{
		Use:     "initinfra",
		Short:   "deploy base infra components to cluster",
		Long:    getInitInfraLong,
		Example: getInitInfraExample,
		// Since this command calls kubectl modules underneath this, we are trying to use same pattern
		// for checking errors, instead of using RunE cobra command field when creating command
		RunE: func(cmd *cobra.Command, args []string) error {
			i.Complete()
			return i.Deploy()
		},
	}
	initFlagsCmd.Flags().BoolVar(&i.DryRun, "dry-run", false,
		"Don't deliver infra to the cluster, see the changes instead")
	initFlagsCmd.Flags().BoolVar(&i.Prune, "prune", false,
		`If set to true, command will delete all kubernetes resources that are not
		defined in airship documents and have airshipit.org/deployed=initinfra label`)
	return initFlagsCmd
}
