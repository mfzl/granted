package registry

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/clio"
	grantedConfig "github.com/common-fate/granted/pkg/config"
	"github.com/common-fate/granted/pkg/testable"

	"github.com/urfave/cli/v2"
)

var RemoveCommand = cli.Command{
	Name:        "remove",
	Description: "Remove subscription to provided registry",
	Usage:       "Remove subscribing to existing registry",
	Action: func(c *cli.Context) error {

		gConf, err := grantedConfig.Load()
		if err != nil {
			return err
		}

		if len(gConf.ProfileRegistryURLS) <= 0 {
			clio.Error("There are no profile registries configured currently.\n Please use 'granted registry add <https://github.com/your-org/your-registry.git>' to add a new registry")

			return nil
		}

		in := survey.Select{Message: "Please select the git repository you would like to unsubscribe:", Options: gConf.ProfileRegistryURLS}
		var out string
		err = testable.AskOne(&in, &out)
		if err != nil {
			return err
		}

		index := -1
		for i, v := range gConf.ProfileRegistryURLS {
			if out == v {
				index = i
				break
			}
		}

		if index != -1 {
			u, err := parseGitURL(out)
			if err != nil {
				return err
			}

			repoDir, err := getRegistryLocation(u)
			if err != nil {
				return err
			}

			err = os.RemoveAll(repoDir)
			if err != nil {
				return err
			}

			gConf.ProfileRegistryURLS = remove(gConf.ProfileRegistryURLS, index)

			if err := gConf.Save(); err != nil {
				return err
			}
		}

		clio.Successf("Successfully removed %s", out)

		return nil
	},
}

func remove(slice []string, i int) []string {
	return append(slice[:i], slice[i+1:]...)
}
