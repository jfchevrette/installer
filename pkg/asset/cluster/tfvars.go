package cluster

import (
	"os"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/tfvars"
	"github.com/pkg/errors"
)

const (
	// TfVarsFileName is the filename for Terraform variables.
	TfVarsFileName  = "terraform.tfvars"
	tfvarsAssetName = "Terraform Variables"
)

// TerraformVariables depends on InstallConfig and
// Ignition to generate the terrafor.tfvars.
type TerraformVariables struct {
	File *asset.File
}

var _ asset.WritableAsset = (*TerraformVariables)(nil)

// Name returns the human-friendly name of the asset.
func (t *TerraformVariables) Name() string {
	return tfvarsAssetName
}

// Dependencies returns the dependency of the TerraformVariable
func (t *TerraformVariables) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&bootstrap.Bootstrap{},
		&machine.Master{},
		&machine.Worker{},
	}
}

// Generate generates the terraform.tfvars file.
func (t *TerraformVariables) Generate(parents asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	bootstrap := &bootstrap.Bootstrap{}
	master := &machine.Master{}
	worker := &machine.Worker{}
	parents.Get(installConfig, bootstrap, master, worker)

	bootstrapIgn := string(bootstrap.Files()[0].Data)

	masterIgn := string(master.Files()[0].Data)
	workerIgn := string(worker.Files()[0].Data)

	data, err := tfvars.TFVars(installConfig.Config, bootstrapIgn, masterIgn, workerIgn)
	if err != nil {
		return errors.Wrap(err, "failed to get Tfvars")
	}
	t.File = &asset.File{
		Filename: TfVarsFileName,
		Data:     data,
	}

	return nil
}

// Files returns the files generated by the asset.
func (t *TerraformVariables) Files() []*asset.File {
	if t.File != nil {
		return []*asset.File{t.File}
	}
	return []*asset.File{}
}

// Load reads the terraform.tfvars from disk.
func (t *TerraformVariables) Load(f asset.FileFetcher) (found bool, err error) {
	file, err := f.FetchByName(TfVarsFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	t.File = file
	return true, nil
}
