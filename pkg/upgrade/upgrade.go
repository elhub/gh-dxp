package upgrade

import "github.com/elhub/gh-dxp/pkg/utils"

// RunUpgrade upgrades gh dxp to the latest version
func RunUpgrade(exe utils.Executor) error {

	_, err := exe.GH("extension", "upgrade", "elhub/gh-dxp", "--force")
	return err

}
