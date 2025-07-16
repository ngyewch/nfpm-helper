package build

type Config struct {
	Name            string              `yaml:"name"`
	Download        DownloadBaseConfig  `yaml:"download"`
	StripComponents int                 `yaml:"strip_components"`
	Packaging       PackagingBaseConfig `yaml:"packaging"`
	Outputs         []Output            `yaml:"outputs"`
}

type Output struct {
	Arch      string          `yaml:"arch"`
	Download  DownloadConfig  `yaml:"download"`
	Packaging PackagingConfig `yaml:"packaging"`
}

type DownloadBaseConfig struct {
	UrlTemplate string `yaml:"url_template"`
}

type DownloadConfig struct {
	UrlTemplate string            `yaml:"url_template"`
	Env         map[string]string `yaml:"env"`
}

type PackagingBaseConfig struct {
	FilenameTemplate string `yaml:"filename_template"`
}

type PackagingConfig struct {
	FilenameTemplate string            `yaml:"filename_template"`
	Env              map[string]string `yaml:"env"`
}
