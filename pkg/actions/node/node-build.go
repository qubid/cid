package node

import (
	"github.com/cidverse/cid/pkg/common/api"
	"github.com/cidverse/cid/pkg/common/command"
	"github.com/rs/zerolog/log"
)

type BuildActionStruct struct{}

// GetDetails retrieves information about the action
func (action BuildActionStruct) GetDetails(ctx api.ActionExecutionContext) api.ActionDetails {
	return api.ActionDetails{
		Stage:     "build",
		Name:      "node-build",
		Version:   "0.1.0",
		UsedTools: []string{"yarn"},
	}
}

// Check evaluates if the action should be executed or not
func (action BuildActionStruct) Check(ctx api.ActionExecutionContext) bool {
	return DetectNodeProject(ctx.ProjectDir)
}

// Execute runs the action
func (action BuildActionStruct) Execute(ctx api.ActionExecutionContext, state *api.ActionStateContext) error {
	// parse package.json
	packageConfig := ParsePackageJSON(ctx.ProjectDir + `/package.json`)

	// dependencies
	command.RunCommand(`yarn install --frozen-lockfile --cache-folder `+api.GetCacheDir(Config.Paths, "yarn"), ctx.Env, ctx.ProjectDir)

	// dependency specific
	reactDependencyVersion, reactDependencyPresent := packageConfig.Dependencies[`react`]
	if reactDependencyPresent {
		log.Debug().Str("react", reactDependencyVersion).Msg("found library")
		ctx.Env["BUILD_PATH"] = ctx.ProjectDir + `/` + Config.Paths.Artifact + `/html` // overwrite build dir - react - react-scripts at v4.0.2+
		ctx.Env["CI"] = "false"                                                        // if ci=true, then react warnings will result in errors - allow warnings // TODO: remove
	}

	// build script
	buildScriptLine, buildScriptPresent := packageConfig.Scripts[`build`]
	if buildScriptPresent {
		log.Debug().Str("build", buildScriptLine).Msg("found build script")
		command.RunCommand(`yarn build --cache-folder `+api.GetCacheDir(Config.Paths, "yarn")+` `+ctx.ProjectDir, ctx.Env, ctx.ProjectDir)
	}

	return nil
}

func init() {
	api.RegisterBuiltinAction(BuildActionStruct{})
}
