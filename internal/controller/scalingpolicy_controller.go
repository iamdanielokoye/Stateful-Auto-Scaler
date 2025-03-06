/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	scalingv1alpha1 "github.com/iamdanielokoye/Stateful-Auto-Scaler/api/v1alpha1"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// ScalingPolicyReconciler reconciles a ScalingPolicy object
type ScalingPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=scaling.example.com,resources=scalingpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scaling.example.com,resources=scalingpolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=scaling.example.com,resources=scalingpolicies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ScalingPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile

func updateStatefulSetReplicas(ctx context.Context, c client.Client, name string, replicas int32) error {
	sts := &appsv1.StatefulSet{}
	err := c.Get(ctx, types.NamespacedName{Name: name, Namespace: "default"}, sts)
	if err != nil {
		return err
	}

	sts.Spec.Replicas = &replicas
	return c.Update(ctx, sts)
}

func (r *ScalingPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ScalingPolicy resource
	var policy scalingv1alpha1.ScalingPolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		log.Error(err, "Failed to get ScalingPolicy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling ScalingPolicy", "name", policy.Name)

	// TODO(user): Add Prometheus integration and scaling logic here
	// Inside the Reconcile function
	cpuUsage, err := PrometheusQuery("avg(rate(container_cpu_usage_seconds_total{namespace='default', pod=~'my-app.*'}[1m]))")
	if err != nil {
		log.Error(err, "Failed to fetch CPU metrics")
		return ctrl.Result{}, err
	}

	// Assume 0.7 as scale-up threshold and 0.3 as scale-down
	var cpuValue float64
	if len(cpuUsage) > 0 {
		cpuValue = float64(cpuUsage[0].Value)
	}

	currentReplicas := int32(len(cpuUsage)) // Mock example, update this logic

	if cpuValue > 0.7 && currentReplicas < 5 { // Max replicas 5
		newReplicas := currentReplicas + 1
		err := updateStatefulSetReplicas(ctx, r.Client, "my-statefulset", newReplicas)
		if err != nil {
			log.Error(err, "Failed to scale up StatefulSet")
		}
	} else if cpuValue < 0.3 && currentReplicas > 1 { // Min replicas 1
		newReplicas := currentReplicas - 1
		err := updateStatefulSetReplicas(ctx, r.Client, "my-statefulset", newReplicas)
		if err != nil {
			log.Error(err, "Failed to scale down StatefulSet")
		}
	}

	return ctrl.Result{}, nil
}

func PrometheusQuery(query string) (model.Vector, error) {
	client, err := api.NewClient(api.Config{
		Address: "http://prometheus-service.monitoring.svc:9090", // Change this if necessary
	})
	if err != nil {
		return nil, err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, _, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		return nil, err
	}

	if vector, ok := result.(model.Vector); ok {
		return vector, nil
	}

	return nil, fmt.Errorf("unexpected result type: %T", result)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScalingPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalingv1alpha1.ScalingPolicy{}).
		Complete(r)
}
