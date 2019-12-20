package models

type Version struct {
	Version     string `json:"version"`
	Revision    string `json:"revision"`
	RevisionURL string `json:"revision_url"`
	Generated   string `json:"generated"`
	Developers  string `json:"developers"`
	BuildNum    string `json:"build_num"`
	DockerImage string `json:"docker_image"`
}
