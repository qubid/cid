package githubaction

import (
	commonapi "github.com/cidverse/cid/pkg/common/api"
	"github.com/cidverse/cid/pkg/core/registry"
	"github.com/cidverse/cid/pkg/core/state"
)

type Executor struct{}

func (e Executor) GetName() string {
	return "githubaction"
}

func (e Executor) GetVersion() string {
	return "0.1.0"
}

func (e Executor) GetType() string {
	return string(registry.ActionTypeGitHubAction)
}

func (e Executor) Execute(ctx *commonapi.ActionExecutionContext, localState *state.ActionStateContext, catalogAction *registry.Action, action *registry.WorkflowAction) error {
	return nil
}
