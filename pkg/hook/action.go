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

// Action run hooks for the given object type as per the event
// Action will skip further hook execution if any error occurred
func (h *Hook) Action(obj interface{}, resourceType int, eventType EventType) error {
	var err error
	for actionType, cfg := range h.Config {
		actionEvent, ok := ActionForEventMap[actionType]
		if !ok {
			continue
		}

		if actionEvent.evType != eventType {
			continue
		}

		switch resourceType {
		case ResourceBackendPVC:
			err = pvcHookAction(cfg.BackendPVCConfig, actionEvent.actOp, obj)
			if err != nil {
				return err
			}
		case ResourceBackendPV:
			err = pvHookAction(cfg.BackendPVConfig, actionEvent.actOp, obj)
			if err != nil {
				return err
			}
		case ResourceNFSService:
			err = serviceHookAction(cfg.NFSServiceConfig, actionEvent.actOp, obj)
			if err != nil {
				return err
			}
		case ResourceNFSPV:
			err = pvHookAction(cfg.NFSPVConfig, actionEvent.actOp, obj)
			if err != nil {
				return err
			}
		case ResourceNFSServerDeployment:
			err = deploymentHookAction(cfg.NFSDeploymentConfig, actionEvent.actOp, obj)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ActionExists will check if action exists for the give resource type and event type
func (h *Hook) ActionExists(resourceType int, eventType EventType) bool {
	_, actionExist := h.availableActions[eventType][resourceType]
	return actionExist
}
