package cmd

import (
	"fmt"
	"github.com/cidverse/cid/pkg/app"
	"github.com/cidverse/cid/pkg/common/api"
	"github.com/cidverse/cid/pkg/common/protectoutput"
	"github.com/cidverse/cid/pkg/common/workflowrun"
	"github.com/cidverse/cid/pkg/core/cliutil"
	"github.com/cidverse/cid/pkg/core/config"
	"github.com/cidverse/cid/pkg/core/rules"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"text/tabwriter"
)

func init() {
	rootCmd.AddCommand(workflowRootCmd)
	workflowRootCmd.AddCommand(workflowListCmd)
	workflowRootCmd.AddCommand(workflowRunCmd)
	workflowRunCmd.Flags().StringArrayP("stage", "s", []string{}, "limit execution to the specified stage(s)")
	workflowRunCmd.Flags().StringArrayP("module", "m", []string{}, "limit execution to the specified module(s)")
}

var workflowRootCmd = &cobra.Command{
	Use:     "workflow",
	Aliases: []string{"wf"},
	Short:   ``,
	Long:    ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
		os.Exit(0)
	},
}

var workflowListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "lists all workflows",
	Run: func(cmd *cobra.Command, args []string) {
		// find project directory and load config
		projectDir := api.FindProjectDir()
		cfg := app.Load(projectDir)
		env := api.GetCIDEnvironment(projectDir)

		// print list
		w := tabwriter.NewWriter(protectoutput.NewProtectedWriter(nil, os.Stdout), 1, 1, 1, ' ', 0)
		_, _ = fmt.Fprintln(w, "WORKFLOW\tRULES\tSTAGES\tACTIONS\tENABLED")
		for _, workflow := range cfg.Workflows {
			_, _ = fmt.Fprintln(w, workflow.Name+"\t"+
				rules.EvaluateRulesAsText(workflow.Rules, rules.GetRuleContext(env))+"\t"+
				strconv.Itoa(len(workflow.Stages))+"\t"+
				strconv.Itoa(workflow.ActionCount())+"\t"+
				cliutil.BoolToChar(workflow.Enabled))
		}
		_ = w.Flush()
	},
}

var workflowRunCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Short:   "runs the specified workflow, requires exactly one argument",
	Run: func(cmd *cobra.Command, args []string) {
		modules, _ := cmd.Flags().GetStringArray("module")
		stages, _ := cmd.Flags().GetStringArray("stage")

		// find project directory and load config
		projectDir := api.FindProjectDir()
		cfg := app.Load(projectDir)
		env := api.GetCIDEnvironment(projectDir)

		var wf *config.Workflow
		if len(args) == 0 {
			// evaluate rules to pick workflow
			wf = workflowrun.FirstWorkflowMatchingRules(cfg.Workflows, env)
		} else if len(args) == 1 {
			// find workflow
			wf = cfg.FindWorkflow(args[0])
			if wf == nil {
				log.Fatal().Str("workflow", args[0]).Msg("workflow does not exist")
			}
		} else {
			// error
			_ = cmd.Help()
			os.Exit(0)
		}

		// run
		log.Info().Str("workflow", wf.Name).Msg("running workflow")
		workflowrun.RunWorkflow(cfg, wf, env, projectDir, stages, modules)
	},
}
