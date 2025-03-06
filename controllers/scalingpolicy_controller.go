package controllers

import (
	"context"
	"fmt"

	scalingv1alpha1 "github.com/iamdanielokoye/Stateful-Auto-Scaler/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ScalingPolicyReconciler reconciles a ScalingPolicy object
type ScalingPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile handles the scaling logic
func (r *ScalingPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var policy scalingv1alpha1.ScalingPolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	fmt.Printf("Reconciling ScalingPolicy: %s\n", policy.Name)

	// Future: Add Prometheus integration and scaling logic

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *ScalingPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&scalingv1alpha1.ScalingPolicy{}).
		Complete(r)
}
