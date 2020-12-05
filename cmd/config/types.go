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

type Environment struct {
	Name string
	Opts *ContainerOpts
}

type config struct {
	Envs struct {
		Values map[string]*ContainerOpts
	}
}
