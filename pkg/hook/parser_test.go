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
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

func getTestHookData(version string) []byte {
	var hook Hook
	hook.Config = make(map[ActionType]HookConfig)
	hook.Config[ActionAddOnCreateVolumeEvent] = HookConfig{
		Name: "createHook",
		BackendPVConfig: &PVHook{
			Annotations: map[string]string{
				"example.io/track": "true",
				"test.io/owner":    "teamA",
			},
			Finalizers: []string{"test.io/tracking-protection"},
		},
		NFSPVConfig: &PVHook{
			Annotations: map[string]string{
				"example.io/track": "true",
				"test.io/owner":    "teamA",
			},
			Finalizers: []string{"test.io/tracking-protection"},
		},

		BackendPVCConfig: &PVCHook{
			Annotations: map[string]string{
				"example.io/track": "true",
				"test.io/owner":    "teamA",
			},
			Finalizers: []string{"test.io/tracking-protection"},
		},

		NFSServiceConfig: &ServiceHook{
			Annotations: map[string]string{
				"example.io/track": "true",
				"test.io/owner":    "teamA",
			},
			Finalizers: []string{"test.io/tracking-protection"},
		},
		NFSDeploymentConfig: &DeploymentHook{
			Annotations: map[string]string{
				"example.io/track": "true",
				"test.io/owner":    "teamA",
			},
			Finalizers: []string{"test.io/tracking-protection"},
		},
	}

	hook.Version = version
	data, _ := yaml.Marshal(hook)
	return data
}

func getTestHookDataWithInvalidAction(version string) []byte {
	var hook Hook
	hook.Config = make(map[ActionType]HookConfig)
	hook.Config["invalidAction"] = HookConfig{
		Name: "createHook",
		BackendPVConfig: &PVHook{
			Annotations: map[string]string{
				"example.io/track": "true",
				"test.io/owner":    "teamA",
			},
			Finalizers: []string{"test.io/tracking-protection"},
		},
	}

	hook.Version = version
	data, _ := yaml.Marshal(hook)
	return data
}

func TestParseHooks(t *testing.T) {
	invalidHookData := `
hook:
NFSDeployment:
    annotations:
      example.io/track: "true"
      test.io/owner: teamA
  finalizers:
    - test.io/tracking-protection
`

	hookWithInvalidAction := getTestHookDataWithInvalidAction("1.0.0")

	tests := []struct {
		name          string
		hookData      []byte
		shouldErrored bool
	}{
		{
			name:          "when correct hook data is passed",
			hookData:      getTestHookData("1.0.0"),
			shouldErrored: false,
		},
		{
			name:          "when invalid versioned hook data is passed",
			hookData:      getTestHookData("0.0.0"),
			shouldErrored: true,
		},
		{
			name:          "when invalid hook data is passed",
			hookData:      []byte(invalidHookData),
			shouldErrored: true,
		},
		{
			name:          "when hook data is having invalid actionEvent",
			hookData:      hookWithInvalidAction,
			shouldErrored: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h, err := ParseHooks(test.hookData)
			assert.Equal(t, test.shouldErrored, err != nil)
			if !test.shouldErrored {
				assert.NotNil(t, h, "Hook obj should not be nil")
			}
		})
	}
}
