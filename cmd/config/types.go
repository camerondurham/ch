package config

type BuildOpts struct {
	DockerfilePath string
	Context        string
	Tag            string
}

type PullOpts struct {
	ImageName string
}

type ContainerOpts struct {
	BuildOpts *BuildOpts
	PullOpts  *PullOpts
	Volume    string
	Shell     string
}
