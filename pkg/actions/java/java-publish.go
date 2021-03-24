package java

import (
	ncicommon "github.com/EnvCLI/normalize-ci/pkg/common"
	"github.com/PhilippHeuer/cid/pkg/common/command"
	"github.com/rs/zerolog/log"
)

// Publish
type PublishActionStruct struct {
	stage   string
	name    string
	version string
}

// GetStage returns the stage
func (n PublishActionStruct) GetStage() string {
	return n.stage
}

// GetName returns the name
func (n PublishActionStruct) GetName() string {
	return n.name
}

// GetVersion returns the name
func (n PublishActionStruct) GetVersion() string {
	return n.version
}

// SetConfig is used to pass a custom configuration to each action
func (n PublishActionStruct) SetConfig(config string) {

}

// Check if this package can handle the current environment
func (n PublishActionStruct) Check(projectDir string, env []string) bool {
	loadConfig(projectDir)
	return DetectJavaProject(projectDir)
}

// Check if this package can handle the current environment
func (n PublishActionStruct) Execute(projectDir string, env []string, args []string) {
	log.Debug().Str("action", n.name).Msg("running action")
	loadConfig(projectDir)

	// get release version
	releaseVersion := ncicommon.GetEnvironment(env, `NCI_COMMIT_REF_RELEASE`)

	// run build
	buildSystem := DetectJavaBuildSystem(projectDir)
	if buildSystem == "gradle" {
		command.RunCommand(`gradlew -Pversion="`+releaseVersion+`" publishAllPublicationsToSonatypeRepository --no-daemon --warning-mode=all --console=plain`, env)
		//command.RunCommand(`gradlew closeRepository --no-daemon --warning-mode=all --console=plain`, env)
		//command.RunCommand(`gradlew releaseRepository --no-daemon --warning-mode=all --console=plain`, env)
	} else if buildSystem == "maven" {
		//
	}
}

// PublishAction
func PublishAction() PublishActionStruct {
	entity := PublishActionStruct{
		stage: "publish",
		name: "java-publish",
		version: "0.1.0",
	}

	return entity
}
