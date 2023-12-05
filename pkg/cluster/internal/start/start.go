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

package start

import (
	"sigs.k8s.io/kind/pkg/cluster/internal/providers"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/log"
)

// Cluster starts the cluster identified by ctx
// explicitKubeconfigPath is --kubeconfig, following the rules from
// https://kubernetes.io/docs/reference/generated/kubectl/kubectl-commands
func Cluster(logger log.Logger, p providers.Provider, name, explicitKubeconfigPath string) error {
	n, err := p.ListNodes(name)
	if err != nil {
		return errors.Wrap(err, "error listing nodes")
	}

	if len(n) > 0 {
		err = p.StartNodes(n)
		if err != nil {
			return err
		}
		logger.V(0).Infof("Started nodes: %q", n)
	}

	// identify external load balancer node
	loadBalancerNode, err := nodeutils.ExternalLoadBalancerNode(n)
	if err != nil {
		return err
	}

	if loadBalancerNode != nil {
		//err = p.DeleteNodes(lb)
		//if err != nil {
		//	return err
		//}

		//var backendServers = map[string]string{}
		//controlPlaneNodes, err := nodeutils.SelectNodesByRole(
		//	n,
		//	constants.ControlPlaneNodeRoleValue,
		//)
		//if err != nil {
		//	return err
		//}
		//for _, n := range controlPlaneNodes {
		//	backendServers[n.String()] = fmt.Sprintf("%s:%d", n.String(), common.APIServerInternalPort)
		//}

		//// create loadbalancer config data
		//loadbalancerConfig, err := loadbalancer.Config(&loadbalancer.ConfigData{
		//	ControlPlanePort: common.APIServerInternalPort,
		//	BackendServers:   backendServers,
		//	IPv6:             false, //ctx.Config.Networking.IPFamily == config.IPv6Family,
		//})
		//if err != nil {
		//	return errors.Wrap(err, "failed to generate loadbalancer config data")
		//}
		//
		//// create loadbalancer config on the node
		//if err := nodeutils.WriteFile(loadBalancerNode, loadbalancer.ConfigPath, loadbalancerConfig); err != nil {
		//	// TODO: logging here
		//	return errors.Wrap(err, "failed to copy loadbalancer config to node")
		//}

		// reload the config. haproxy will reload on SIGHUP
		if err := loadBalancerNode.Command("kill", "-s", "HUP", "1").Run(); err != nil {
			return errors.Wrap(err, "failed to reload loadbalancer")
		}

		logger.V(0).Infof("Reloaded node: %q", loadBalancerNode)
	}

	return nil
}
