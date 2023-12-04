/*
Copyright 2018 The Kubernetes Authors.

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

// Package cluster implements the `create cluster` command
package cluster

import (
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/log"

	"sigs.k8s.io/kind/pkg/internal/cli"
	"sigs.k8s.io/kind/pkg/internal/runtime"
)

type flagpole struct {
	Name       string
	Kubeconfig string
}

// NewCommand returns a new cobra.Command for cluster creation
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args:  cobra.NoArgs,
		Use:   "cluster",
		Short: "Starts a local Kubernetes cluster",
		Long:  "Starts a local Kubernetes cluster using Docker container 'nodes'",
		RunE: func(cmd *cobra.Command, args []string) error {
			cli.OverrideDefaultName(cmd.Flags())
			return startCluster(logger, flags)
		},
	}
	cmd.Flags().StringVarP(
		&flags.Name,
		"name",
		"n",
		"",
		"cluster name, overrides KIND_CLUSTER_NAME, config (default kind)",
	)
	cmd.Flags().StringVar(
		&flags.Kubeconfig,
		"kubeconfig",
		"",
		"sets kubeconfig path instead of $KUBECONFIG or $HOME/.kube/config",
	)
	return cmd
}

func startCluster(logger log.Logger, flags *flagpole) error {
	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		runtime.GetDefault(logger),
	)
	// Start individual cluster
	logger.V(0).Infof("Starting cluster %q ...", flags.Name)
	if err := provider.Start(flags.Name, flags.Kubeconfig); err != nil {
		return errors.Wrapf(err, "failed to start cluster %q", flags.Name)
	}
	return nil
}
