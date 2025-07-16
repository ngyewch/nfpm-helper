package generate

type IndexConfig struct {
	Packages []IndexEntry `yaml:"packages"`
}

type IndexEntry struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
}

type Config struct {
	Repositories []RepositoryConfig `yaml:"repositories"`
}

type RepositoryConfig struct {
	Source   string          `yaml:"source"`
	Version  string          `yaml:"version"`
	Type     string          `yaml:"type"`
	Packages []PackageConfig `yaml:"packages"`
}

type PackageConfig struct {
	Name      string   `yaml:"name"`
	Version   string   `yaml:"version"`
	Archs     []string `yaml:"archs"`
	Packagers []string `yaml:"packagers"`
}
