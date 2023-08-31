package controllers

import (
	"context"
	"github.com/myoperator/common_operator/pkg/apis/generic/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type GenericController struct {
	client.Client
	E record.EventRecorder
}

func NewGenericController(e record.EventRecorder) *GenericController {
	return &GenericController{E: e}
}

// Reconcile 调节方法
func (r *GenericController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	g := &v1alpha1.Generic{}

	err := r.Get(ctx, req.NamespacedName, g)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	// FIXME: 需要把interface{}转掉的float64转回int
	// https://bbs.huaweicloud.com/blogs/309752
	for k, v := range g.Spec.Template {
		vv, ok := v.(float64)
		if ok {
			vvv := int(vv)
			g.Spec.Template[k] = vvv
		}
	}

	// 处理调协的主要逻辑
	err = handler(g)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *GenericController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

func (r *GenericController) OnUpdatePodHandler(event event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.ObjectNew.GetOwnerReferences() {
		if ref.Kind == v1alpha1.GenericKind && ref.APIVersion == v1alpha1.GenericApiVersion {
			// 重新放入Reconcile调协方法
			limitingInterface.Add(reconcile.Request{
				types.NamespacedName{
					Name: ref.Name, Namespace: event.ObjectNew.GetNamespace(),
				},
			})
		}
	}

}

func (r *GenericController) OnDeletePodHandler(event event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.Object.GetOwnerReferences() {
		if ref.Kind == v1alpha1.GenericKind && ref.APIVersion == v1alpha1.GenericApiVersion {
			// 重新入列，这样删除pod后，就会进入调和loop，发现ownerReference还在，会立即创建出新的pod。
			klog.Info("delete pod: ", event.Object.GetName(), event.Object.GetObjectKind())
			limitingInterface.Add(reconcile.Request{
				NamespacedName: types.NamespacedName{Name: ref.Name,
					Namespace: event.Object.GetNamespace()}})
		}
	}
}
