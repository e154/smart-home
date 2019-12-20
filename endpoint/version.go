package endpoint

import (
	m "github.com/e154/smart-home/models"
	"github.com/e154/smart-home/version"
)

type VersionEndpoint struct {
	*CommonEndpoint
}

func NewVersionEndpoint(common *CommonEndpoint) *VersionEndpoint {
	return &VersionEndpoint{
		CommonEndpoint: common,
	}
}

func (v *VersionEndpoint) ServerVersion() (ver *m.Version) {

	ver = &m.Version{
		Version:     version.VersionString,
		Revision:    version.RevisionString,
		RevisionURL: version.RevisionURLString,
		Generated:   version.GeneratedString,
		Developers:  version.DevelopersString,
		BuildNum:    version.BuildNumString,
		DockerImage: version.DockerImageString,
	}

	return
}
