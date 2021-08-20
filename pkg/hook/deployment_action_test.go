package hook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeployment_hook_action(t *testing.T) {
	tests := []struct {
		name        string
		hook        *DeploymentHook
		obj         interface{}
		expectedObj interface{}
		actionType  HookActionType
	}{
		{
			name:        "when deployment hook is nil, object should not be modified",
			hook:        nil,
			obj:         generateFakeDeploymentObj("ns1", "name1", map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			expectedObj: generateFakeDeploymentObj("ns1", "name1", map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			actionType:  HookActionAdd,
		},
		{
			name:        "when deployment hook is configured to add metadata, object should be modified",
			hook:        buildDeploymentHook(map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			obj:         generateFakeDeploymentObj("ns2", "name2", nil, nil),
			expectedObj: generateFakeDeploymentObj("ns2", "name2", map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			actionType:  HookActionAdd,
		},
		{
			name:        "when deployment hook is configured to remove metadata, object should be modified",
			hook:        buildDeploymentHook(map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			obj:         generateFakeDeploymentObj("ns3", "name3", map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			expectedObj: generateFakeDeploymentObj("ns3", "name3", map[string]string{}, []string{}),
			actionType:  HookActionRemove,
		},
		{
			name:        "when deployment hook is configured to remove non-existing metadata, object should not be modified",
			hook:        buildDeploymentHook(map[string]string{"test.com/key": "val"}, []string{"test.com/finalizer"}),
			obj:         generateFakeDeploymentObj("ns4", "name4", map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			expectedObj: generateFakeDeploymentObj("ns4", "name4", map[string]string{"test.io/key": "val"}, []string{"test.io/finalizer"}),
			actionType:  HookActionRemove,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotNil(t, test.obj, "object should not be nil")
			err := deployment_hook_action(test.hook, test.actionType, test.obj)
			assert.Nil(t, err, "deployment_hook_action returned error")
			assert.Equal(t, test.expectedObj, test.obj, "object should match")
		})
	}
}
