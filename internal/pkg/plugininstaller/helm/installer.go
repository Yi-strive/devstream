package helm

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/helm"
	"github.com/devstream-io/devstream/pkg/util/k8s"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	DefaultCreateOperations = plugininstaller.ExecuteOperations{
		DealWithNsWhenInstall,
		InstallOrUpdate,
	}
	DefaultUpdateOperations = plugininstaller.ExecuteOperations{
		InstallOrUpdate,
	}
	DefaultDeleteOperations = plugininstaller.ExecuteOperations{
		Delete,
	}
	DefaultTerminateOperations = plugininstaller.TerminateOperations{
		Delete,
	}
)

// InstallOrUpdate will install or update service by input options
func InstallOrUpdate(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	h, err := helm.NewHelm(opts.GetHelmParam())
	if err != nil {
		return err
	}

	log.Info("Creating or updating helm chart ...")
	if err := h.InstallOrUpgradeChart(); err != nil {
		log.Errorf("Failed to install or upgrade the chart: %s.", err)
		return err
	}
	return err
}

// DealWithNsWhenInstall will create namespace by input options
func DealWithNsWhenInstall(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	log.Debugf("Prepare to create the namespace: %s.", opts.GetNamespace())

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return err
	}
	return kubeClient.UpsertNameSpace(opts.Chart.Namespace)
}

// Delete will delete service base on input options
func Delete(options plugininstaller.RawOptions) error {
	opts, err := NewOptions(options)
	if err != nil {
		return err
	}
	h, err := helm.NewHelm(opts.GetHelmParam())
	if err != nil {
		return err
	}

	log.Infof("Uninstalling %s helm chart.", opts.GetReleaseName())
	if err = h.UninstallHelmChartRelease(); err != nil {
		return err
	}
	return nil
}
