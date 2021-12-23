/*
Copyright 2021 The Rook Authors. All rights reserved.

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

package csi

import (
	"github.com/alauda/nativestor/pkg/operator"
	"regexp"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// predicateController is the predicate function to trigger reconcile on operator configuration cm change
func predicateController() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			// if the operator configuration file is created we want to reconcile
			if cm, ok := e.Object.(*v1.ConfigMap); ok {
				// We don't want to use cm.Generation here, it case the operator was stopped and the
				// ConfigMap was created
				return cm.Name == operator.OperatorSettingConfigMapName
			}
			return false
		},

		UpdateFunc: func(e event.UpdateEvent) bool {
			resourceQtyComparer := cmp.Comparer(func(x, y resource.Quantity) bool { return x.Cmp(y) == 0 })
			if old, ok := e.ObjectOld.(*v1.ConfigMap); ok {
				if new, ok := e.ObjectNew.(*v1.ConfigMap); ok {
					if old.Name == operator.OperatorSettingConfigMapName && new.Name == operator.OperatorSettingConfigMapName {
						diff := cmp.Diff(old.Data, new.Data, resourceQtyComparer)
						logger.Debugf("operator configmap diff:\n %s", diff)
						return findCSIChange(diff)
					}
				}
			}

			return false
		},

		DeleteFunc: func(e event.DeleteEvent) bool {
			return false
		},

		GenericFunc: func(e event.GenericEvent) bool {
			return false
		},
	}
}

func findCSIChange(str string) bool {
	var re = regexp.MustCompile(`\"RAW_DEVICE_|\"CSI_|\"KUBELET_.`)
	found := re.FindAllString(str, -1)
	if len(found) > 0 {
		for _, match := range found {
			logger.Infof("raw-device csi config changed with: %q", match)
		}
		return true
	}
	return false
}
