/*
Copyright 2024 Wayne Wang<net_use@bzhy.com>.

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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strings"
	"time"

	monitoringv1beta1 "github.com/kubesysadm/kubesysadm/api/v1beta1"
)

// PodCleanRuleReconciler reconciles a PodCleanRule object
type PodCleanRuleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.sysadm.cn,resources=podcleanrule,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.sysadm.cn,resources=podcleanrule/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.sysadm.cn,resources=podcleanrule/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the PodCleanRule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *PodCleanRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("req info: ", "NS:", req.Namespace, "Name: ", req.Name, "toString:", req.String())

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodCleanRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	b := ctrl.NewControllerManagedBy(mgr)
	b.For(&monitoringv1beta1.PodCleanRule{})
	b.Watches(&corev1.Pod{}, handler.EnqueueRequestsFromMapFunc(r.startPodCleanProcess),
		builder.WithPredicates(predicate.GenerationChangedPredicate{}))
	return b.Complete(r)
}

func (r *PodCleanRuleReconciler) startPodCleanProcess(ctx context.Context, pod client.Object) []reconcile.Request {
	logger := log.FromContext(ctx)

	var podCleanRuleList monitoringv1beta1.PodCleanRuleList
	e := r.List(ctx, &podCleanRuleList)
	if e != nil {
		logger.Error(e, "get pod clean rule list error")
		return nil
	}

	for _, item := range podCleanRuleList.Items {
		e := r.cleaningPodForRule(ctx, item)
		if e != nil {
			logger.Error(e, "cleaning pods with rule error ", "namespace", item.Spec.NameSpace, "age", item.Spec.Age)
		}
	}

	return nil
}

func (r *PodCleanRuleReconciler) cleaningPodForRule(ctx context.Context, rule monitoringv1beta1.PodCleanRule) error {
	logger := log.FromContext(ctx)

	ns := rule.Spec.NameSpace
	age := rule.Spec.Age
	prefixName := strings.TrimSpace(rule.Spec.PrefixName)

	podList := corev1.PodList{}
	listOpts := &client.ListOptions{Namespace: ns}
	e := r.Client.List(ctx, &podList, listOpts)
	if e != nil {
		return e
	}

	for _, item := range podList.Items {
		if item.Status.Phase == corev1.PodPending || item.Status.Phase == corev1.PodRunning || item.Status.Phase == corev1.PodUnknown {
			continue
		}
		podName := item.Name
		createTime := item.ObjectMeta.CreationTimestamp
		nowTime := time.Now()
		difference := nowTime.Sub(createTime.Time)
		if difference.Seconds() > float64(age*1000*1000) {
			if prefixName != "" && !strings.HasPrefix(podName, prefixName) {
				continue
			}
			deletePod := item
			e := r.Client.Delete(ctx, &deletePod, &client.DeleteOptions{})
			if e != nil {
				logger.Error(e, "clean pod error", "namespace", ns, "podName", podName)
			} else {
				logger.Info("pod has be deleted ", "namespace", ns, "podName", podName)
			}

		}
	}

	return nil
}
