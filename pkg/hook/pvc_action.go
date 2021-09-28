/*
Copyright 2021 The OpenEBS Authors.

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

package hook

import (
	"github.com/openebs/dynamic-nfs-provisioner/pkg/helper"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

// pvcHookAction will execute the given hook config on the object for the given action
func pvcHookAction(hookCfg *PVCHook, action ActionOp, obj interface{}) error {
	if hookCfg == nil {
		return nil
	}

	pvcObj, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		return errors.Errorf("%T is not a PersistentVolumeClaim type", obj)
	}

	switch action {
	case ActionOpAddOrUpdate:
		pvcHookActionAdd(pvcObj, *hookCfg)
	case ActionOpRemove:
		pvcHookActionRemove(pvcObj, *hookCfg)
	}

	return nil
}

// pvcHookActionAdd will add the given hook config to the given object
func pvcHookActionAdd(obj *corev1.PersistentVolumeClaim, hookCfg PVCHook) {
	if len(hookCfg.Annotations) != 0 {
		AddAnnotations(&obj.ObjectMeta, hookCfg.Annotations)
	}

	if len(hookCfg.Finalizers) != 0 {
		helper.AddFinalizers(&obj.ObjectMeta, hookCfg.Finalizers)
	}
	return
}

// pvcHookActionRemove will remove the given hook config to the given object
func pvcHookActionRemove(obj *corev1.PersistentVolumeClaim, hookCfg PVCHook) {
	if len(hookCfg.Annotations) != 0 {
		helper.RemoveAnnotations(&obj.ObjectMeta, hookCfg.Annotations)
	}

	if len(hookCfg.Finalizers) != 0 {
		helper.RemoveFinalizers(&obj.ObjectMeta, hookCfg.Finalizers)
	}
	return
}
