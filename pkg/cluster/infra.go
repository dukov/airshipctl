package cluster

import (
	"path"

	"sigs.k8s.io/kustomize/v3/pkg/fs"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/k8s"
)

const (
	getInitInfraLong = (`
		Deploy base infrastructure to kubernetes cluster`)

	getInitInfraExample = (`
		# deploy infra to cluster
		airshipctl cluster initinfra`)

	// Annotations
	airshipBase         = "airshipit.org"
	initinfraAnnotation = airshipBase + "/" + "initinfra"
	deployedByLabel     = airshipBase + "/" + "deployed"
	initinfraIdentifier = "initinfra"
)

// Infra is an abstraction used to initialize base infrastructure
type Infra struct {
	fs.FileSystem
	*environment.AirshipCTLSettings
	*k8s.Kubectl

	DryRun   bool
	Prune    bool
	Selector string
}

// NewInfra return instance of Infra
func NewInfra(rs *environment.AirshipCTLSettings) *Infra {
	// At this point AirshipCTLSettings may not be fully initialized
	infra := &Infra{AirshipCTLSettings: rs}
	return infra
}

// Complete Builds FileSystem and Kubectl in runtime
func (infra *Infra) Complete() {
	infra.FileSystem = k8s.Buffer{FileSystem: fs.MakeRealFS()}
	infra.Kubectl = k8s.NewKubectlFromKubeconfigPath(infra.KubeConfigPath())
}

// Deploy method deploys documents
func (infra *Infra) Deploy() error {

	ao, err := infra.ApplyOptions()
	if err != nil {
		return err
	}

	of := "json"
	ao.PrintFlags.OutputFormat = &of
	ao.DryRun = infra.DryRun
	// If prune is true, set selector for purning
	if infra.Prune {
		ao.Prune = infra.Prune
		ao.Selector = deployedByLabel + "=" + initinfraIdentifier
	}

	// TODO (kkalynovskyi) Fix this when config module will provide path to bundle directory.
	dir, _ := path.Split(infra.AirshipCTLSettings.AirshipConfigPath())
	b, err := document.NewBundle(infra.FileSystem, dir, "")
	if err != nil {
		return err
	}

	// Get documents that are annotated to belong to initinfra
	docs, err := b.GetByAnnotation(initinfraAnnotation)
	if err != nil {
		return err
	}

	// Label every document indicating that it was deployed by initinfra module for further reference
	// This may be used later to get all resources that are part of initinfra module, for monitoring, alerting
	// upgrading etc...
	// also if prune is set to true, this fulfills requirement for all labeled document to be labeled.
	// Pruning by annotation is not available, therefore we need to use label.
	for _, doc := range docs {

		res := doc.GetKustomizeResource()
		labels := res.GetLabels()
		labels[deployedByLabel] = initinfraIdentifier
		res.SetLabels(labels)
		doc.SetKustomizeResource(&res)

	}

	return infra.Apply(docs, ao)
}
