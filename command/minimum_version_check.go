package command

import (
	"code.cloudfoundry.org/cli/version"
	"github.com/blang/semver"
)

func MinimumAPIVersionCheck(current string, minimum string) error {
	if current == version.DefaultVersion || minimum == "" {
		return nil
	}

	currentSemvar, err := semver.Make(current)
	if err != nil {
		return err
	}

	minimumSemvar, err := semver.Make(minimum)
	if err != nil {
		return err
	}

	if currentSemvar.Compare(minimumSemvar) == -1 {
		return MinimumAPIVersionNotMetError{
			CurrentVersion: current,
			MinimumVersion: minimum,
		}
	}

	return nil
}

func WarnAPIVersionCheck(config Config, ui UI) error {
	// TODO: make private and refactor commands that use
	err := MinimumAPIVersionCheck(config.BinaryVersion(), config.MinCLIVersion())

	if _, ok := err.(MinimumAPIVersionNotMetError); ok {
		ui.DisplayWarning("Cloud Foundry API version {{.APIVersion}} requires CLI version {{.MinCLIVersion}}. You are currently on version {{.BinaryVersion}}. To upgrade your CLI, please visit: https://github.com/cloudfoundry/cli#downloads",
			map[string]interface{}{
				"APIVersion":    config.APIVersion(),
				"MinCLIVersion": config.MinCLIVersion(),
				"BinaryVersion": config.BinaryVersion(),
			})
		ui.DisplayNewline()
		return nil
	}

	// Only error if there was an issue in parsing versions
	return err
}
