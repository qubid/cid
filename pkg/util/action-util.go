package util

import (
	"github.com/PhilippHeuer/cid/pkg/common/api"
	"github.com/PhilippHeuer/cid/pkg/actions/golang"
)

// GetName returns the name
func GetAllActions() []api.ActionStep {
	var actions []api.ActionStep
	actions = append(actions, golang.BuildAction())
	actions = append(actions, golang.TestAction())

	return actions
}

func FindAction(stage string, projectDir string) api.ActionStep {
	for _, action := range GetAllActions() {
		if stage == action.GetStage() && action.Check(projectDir) {
			return action
		}
	}

	return nil
}
