/*
Copyright 2021.

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

package controllers

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	demov1 "helloworld-operator/api/v1"
)

// HelloworldReconciler reconciles a Helloworld object
type HelloworldReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demo.hw.io,resources=helloworlds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.hw.io,resources=helloworlds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.hw.io,resources=helloworlds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Helloworld object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *HelloworldReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("helloworld", req.NamespacedName)

	// your logic here
	helloworld := &demov1.Helloworld{}
	if err := r.Get(ctx, req.NamespacedName, helloworld); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if helloworld.StatusSetDefault() {
		if err := r.Status().Update(ctx, helloworld); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	p := &v1.Pod{}
	err := r.Get(ctx, req.NamespacedName, p)
	if err == nil {
		if v1.PodSucceeded == p.Status.Phase && demov1.Running != helloworld.Status.Phase {
			helloworld.Status.Phase = demov1.Running
			if err := r.Status().Update(ctx, helloworld); err != nil {
				return ctrl.Result{}, err
			}
		} else if v1.PodSucceeded != p.Status.Phase && demov1.Pending != helloworld.Status.Phase {
			helloworld.Status.Phase = demov1.Pending
			if err := r.Status().Update(ctx, helloworld); err != nil {
				return ctrl.Result{}, err
			}
		} else {
			return ctrl.Result{}, nil
		}
	} else {
		if errors.IsNotFound(err) {
			pod := makePod(helloworld)
			if err := controllerutil.SetControllerReference(helloworld, pod, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			if err := r.Create(ctx, pod); err != nil && !errors.IsAlreadyExists(err) {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, err
		} else {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
func makePod(hello *demov1.Helloworld) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      hello.Name,
			Namespace: hello.Namespace,
		}, Spec: *hello.Spec.Template.Spec.DeepCopy(),
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *HelloworldReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.Helloworld{}).
		Complete(r)
}
