/*
Copyright The Kubernetes Authors.

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

package disruption

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/utils/clock"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	v1 "sigs.k8s.io/karpenter/pkg/apis/v1"
	"sigs.k8s.io/karpenter/pkg/controllers/state"
	"sigs.k8s.io/karpenter/pkg/metrics"
	"sigs.k8s.io/karpenter/pkg/utils/node"
	nodeclaimutil "sigs.k8s.io/karpenter/pkg/utils/nodeclaim"
)

// Emptiness is a nodeclaim sub-controller that adds or removes status conditions on empty nodeclaims based on TTLSecondsAfterEmpty
type Emptiness struct {
	kubeClient client.Client
	cluster    *state.Cluster
	clock      clock.Clock
}

//nolint:gocyclo
func (e *Emptiness) Reconcile(ctx context.Context, nodePool *v1.NodePool, nodeClaim *v1.NodeClaim) (reconcile.Result, error) {
	hasEmptyCondition := nodeClaim.StatusConditions().Get(v1.ConditionTypeEmpty) != nil

	// From here there are a few scenarios to handle:
	// 1. If ConsolidationPolicyWhenEmpty is not configured or ConsolidateAfter isn't configured, remove the emptiness status condition
	if nodePool.Spec.Disruption.ConsolidationPolicy != v1.ConsolidationPolicyWhenEmpty ||
		nodePool.Spec.Disruption.ConsolidateAfter == nil ||
		nodePool.Spec.Disruption.ConsolidateAfter.Duration == nil {
		if hasEmptyCondition {
			_ = nodeClaim.StatusConditions().Clear(v1.ConditionTypeEmpty)
			log.FromContext(ctx).V(1).Info("removing emptiness status condition, emptiness is disabled")
		}
		return reconcile.Result{}, nil
	}
	// 2. If NodeClaim is not initialized, remove the emptiness status condition
	if !nodeClaim.StatusConditions().Get(v1.ConditionTypeInitialized).IsTrue() {
		if hasEmptyCondition {
			_ = nodeClaim.StatusConditions().Clear(v1.ConditionTypeEmpty)
			log.FromContext(ctx).V(1).Info("removing emptiness status condition, isn't initialized")
		}
		return reconcile.Result{}, nil
	}
	// Get the node to check for pods scheduled to it
	n, err := nodeclaimutil.NodeForNodeClaim(ctx, e.kubeClient, nodeClaim)
	if err != nil {
		// 3. If Node mapping doesn't exist, remove the emptiness status condition
		if nodeclaimutil.IsDuplicateNodeError(err) || nodeclaimutil.IsNodeNotFoundError(err) {
			_ = nodeClaim.StatusConditions().Clear(v1.ConditionTypeEmpty)
			if hasEmptyCondition {
				log.FromContext(ctx).V(1).Info("removing emptiness status condition, doesn't have a single node mapping")
			}
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	// Node is empty, but it is in-use per the last scheduling round, so we don't consider it empty
	// We perform a short requeue if the node is nominated, so we can check the node for emptiness when the node
	// nomination time ends since we don't watch node nomination events
	// 4. If the Node is nominated for pods to schedule to it, remove the emptiness status condition
	if e.cluster.IsNodeNominated(n.Spec.ProviderID) {
		_ = nodeClaim.StatusConditions().Clear(v1.ConditionTypeEmpty)
		if hasEmptyCondition {
			log.FromContext(ctx).V(1).Info("removing emptiness status condition, is nominated for pods")
		}
		return reconcile.Result{RequeueAfter: time.Second * 30}, nil
	}
	pods, err := node.GetReschedulablePods(ctx, e.kubeClient, n)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("retrieving node pods, %w", err)
	}
	// 5. If there are pods that are actively scheduled to the Node, remove the emptiness status condition
	if len(pods) > 0 {
		_ = nodeClaim.StatusConditions().Clear(v1.ConditionTypeEmpty)
		if hasEmptyCondition {
			log.FromContext(ctx).V(1).Info("removing emptiness status condition, not empty")
		}
		return reconcile.Result{}, nil
	}
	// 6. Otherwise, add the emptiness status condition
	nodeClaim.StatusConditions().SetTrue(v1.ConditionTypeEmpty)
	if !hasEmptyCondition {
		log.FromContext(ctx).V(1).Info("marking empty")

		metrics.NodeClaimsDisruptedCounter.With(prometheus.Labels{
			metrics.TypeLabel:     metrics.EmptinessReason,
			metrics.NodePoolLabel: nodeClaim.Labels[v1.NodePoolLabelKey],
		}).Inc()
	}
	return reconcile.Result{}, nil
}
