package buildpack_test

import (
	"github.com/cloudfoundry/cli/cf/api/apifakes"
	"github.com/cloudfoundry/cli/cf/commandregistry"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/flags"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	"github.com/cloudfoundry/cli/cf/commands/buildpack"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListBuildpacks", func() {
	var (
		ui                  *testterm.FakeUI
		buildpackRepo       *apifakes.OldFakeBuildpackRepository
		requirementsFactory *testreq.FakeReqFactory
		deps                commandregistry.Dependency
	)

	updateCommandDependency := func(pluginCall bool) {
		deps.UI = ui
		deps.RepoLocator = deps.RepoLocator.SetBuildpackRepository(buildpackRepo)
		commandregistry.Commands.SetCommand(commandregistry.Commands.FindCommand("buildpacks").SetDependency(deps, pluginCall))
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		buildpackRepo = new(apifakes.OldFakeBuildpackRepository)
		requirementsFactory = &testreq.FakeReqFactory{}
	})

	runCommand := func(args ...string) bool {
		return testcmd.RunCLICommand("buildpacks", args, requirementsFactory, updateCommandDependency, false, ui)
	}

	Context("when arguments are provided", func() {
		var cmd commandregistry.Command
		var flagContext flags.FlagContext

		BeforeEach(func() {
			cmd = &buildpack.ListBuildpacks{}
			cmd.SetDependency(deps, false)
			flagContext = flags.NewFlagContext(cmd.MetaData().Flags)
		})

		It("should fail with usage", func() {
			flagContext.Parse("blahblah")

			reqs := cmd.Requirements(requirementsFactory, flagContext)

			err := testcmd.RunRequirements(reqs)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Incorrect Usage"))
			Expect(err.Error()).To(ContainSubstring("No argument required"))
		})
	})

	It("fails requirements when login fails", func() {
		Expect(runCommand()).To(BeFalse())
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
		})

		It("lists buildpacks", func() {
			p1 := 5
			p2 := 10
			p3 := 15
			t := true
			f := false

			buildpackRepo.Buildpacks = []models.Buildpack{
				{Name: "Buildpack-1", Position: &p1, Enabled: &t, Locked: &f},
				{Name: "Buildpack-2", Position: &p2, Enabled: &f, Locked: &t},
				{Name: "Buildpack-3", Position: &p3, Enabled: &t, Locked: &f},
			}

			runCommand()

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting buildpacks"},
				[]string{"buildpack", "position", "enabled"},
				[]string{"Buildpack-1", "5", "true", "false"},
				[]string{"Buildpack-2", "10", "false", "true"},
				[]string{"Buildpack-3", "15", "true", "false"},
			))
		})

		It("tells the user if no build packs exist", func() {
			runCommand()
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting buildpacks"},
				[]string{"No buildpacks found"},
			))
		})
	})

})
