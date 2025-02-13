package jenkinspipeline

import (
	_ "embed"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/util"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/types"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

//go:embed tpl/Jenkinsfile_offline.tpl
var offlineJenkinsScript string

// setJenkinsDefault config default fields for usage
func setJenkinsDefault(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	const ciType = server.CIJenkinsType
	opts := new(jobOptions)
	if err := util.DecodePlugin(options, opts); err != nil {
		return nil, err
	}

	// if jenkins is offline, just use offline Jenkinsfile
	if opts.needOfflineConfig() {
		opts.CIFileConfig = &cifile.CIFileConfig{
			Type: ciType,
			ConfigContentMap: map[string]string{
				"Jenkinsfile": offlineJenkinsScript,
			},
			Vars: opts.Pipeline.GenerateCIFileVars(opts.ProjectRepo),
		}
	} else {
		if opts.Pipeline.ConfigLocation == "" {
			opts.Pipeline.ConfigLocation = "https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/jenkins-pipeline/general/Jenkinsfile"
		}
		opts.CIFileConfig = opts.Pipeline.BuildCIFileConfig(ciType, opts.ProjectRepo)
	}
	// set field value if empty
	if opts.Jenkins.Namespace == "" {
		opts.Jenkins.Namespace = "jenkins"
	}
	if opts.JobName == "" {
		opts.JobName = jenkinsJobName(opts.ProjectRepo.GetRepoName())
	}
	return types.EncodeStruct(opts)
}

// validateJenkins will validate jenkins jobName field and jenkins field
func validateJenkins(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts := new(jobOptions)
	if err := util.DecodePlugin(options, opts); err != nil {
		return nil, err
	}

	if err := validator.CheckStructError(opts).Combine(); err != nil {
		return nil, err
	}

	// check jenkins job name
	if err := opts.JobName.checkValid(); err != nil {
		log.Debugf("jenkins validate pipeline invalid: %+v", err)
		return nil, err
	}

	return options, nil
}
