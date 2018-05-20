package regression

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"gopkg.in/src-d/go-log.v0"
)

// RepoDescription holds the information about a single repository
type RepoDescription struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Complexity  int    `json:"complexity"`
}

var defaultRepos = []RepoDescription{
	{
		Name:        "cangallo",
		URL:         "git://github.com/jfontan/cangallo.git",
		Description: "Small repository that should be fast to clone",
		Complexity:  0,
	},
	{
		Name:        "octoprint-tft",
		URL:         "https://github.com/mcuadros/OctoPrint-TFT",
		Description: "Small repository that should be fast to clone",
		Complexity:  0,
	},
	{
		Name:        "upsilon",
		URL:         "git://github.com/upsilonproject/upsilon-common.git",
		Description: "Average repository",
		Complexity:  1,
	},
	{
		Name:        "numpy",
		URL:         "git://github.com/numpy/numpy.git",
		Description: "Average repository",
		Complexity:  2,
	},
	{
		Name:        "tensorflow",
		URL:         "git://github.com/tensorflow/tensorflow.git",
		Description: "Average repository",
		Complexity:  3,
	},
	{
		Name:        "bismuth",
		URL:         "git://github.com/hclivess/Bismuth.git",
		Description: "Big files repo (100Mb)",
		Complexity:  4,
	},
}

// Repositories struct has the information about a set of repositories and
// functionality to download them.
type Repositories struct {
	Repos  []RepoDescription
	config Config
}

// NewDefaultRepositories creates a new Repositories with default set.
// hardcoded.
func NewDefaultRepositories(config Config) *Repositories {
	return NewRepositories(config, defaultRepos)
}

// NewRepositoriesFromYaml creates a new Repositories struct loading the list
// of repos from a yaml file.
func NewRepositoriesFromYaml(
	config Config,
	file string,
) (*Repositories, error) {
	repos, err := loadReposYaml(file)
	if err != nil {
		return nil, err
	}

	return NewRepositories(config, repos), nil
}

func loadReposYaml(file string) ([]RepoDescription, error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var repos []RepoDescription
	err = yaml.Unmarshal(text, &repos)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

// NewRepositories creates a new Repositories set.
func NewRepositories(config Config, repos []RepoDescription) *Repositories {
	return &Repositories{
		Repos:  repos,
		config: config,
	}
}

// Download clones all repositories in the list that have equal or lower
// complexity specified in config.
func (r *Repositories) Download() error {
	for _, repo := range r.Repos {
		if repo.Complexity > r.config.Complexity {
			continue
		}

		logger, _ := log.New()
		logger = logger.New(log.Fields{"name": repo.Name})

		path := filepath.Join(r.config.RepositoriesCache, repo.Name)
		exist, err := fileExist(path)
		if err != nil {
			return err
		}
		if exist {
			logger.Debugf("Repository already downloaded")
			continue
		}

		logger = logger.New(log.Fields{
			"url":  repo.URL,
			"path": path,
		})

		logger.Debugf("Downloading repository")
		err = os.MkdirAll(r.config.RepositoriesCache, 0755)
		if err != nil {
			return err
		}

		err = downloadRepo(logger, repo.URL, path)
		if err != nil {
			logger.Errorf(err, "Could not download repository")
			return err
		}
	}

	return nil
}

// Path returns the repository cache path.
func (r *Repositories) Path() string {
	return r.config.RepositoriesCache
}

// Names returns the names of repositories withing concurrency level.
func (r *Repositories) Names() []string {
	names := make([]string, 0, len(r.Repos))
	for _, repo := range r.Repos {
		if repo.Complexity <= r.config.Complexity {
			names = append(names, repo.Name)
		}
	}

	return names
}

// LinksDir returns a path of a temporary directory with repos within
// the config complexity.
func (r *Repositories) LinksDir() (string, error) {
	dir, err := createTempDir()
	if err != nil {
		return "", err
	}

	for _, name := range r.Names() {
		from := filepath.Join(r.config.RepositoriesCache, name)
		to := filepath.Join(dir, name)

		err = recursiveCopy(from, to)
		if err != nil {
			os.RemoveAll(dir)
			return "", err
		}
	}

	return dir, err
}

func downloadRepo(l log.Logger, url, path string) error {
	downloadPath := fmt.Sprintf("%s.download", path)
	exist, err := fileExist(downloadPath)
	if err != nil {
		return err
	}

	if exist {
		err = os.RemoveAll(downloadPath)
		if err != nil {
			return err
		}
	}

	clone, err := NewExecutor("git", "clone", "--mirror", url, downloadPath)
	if err != nil {
		l.Errorf(err, "Could not create executor")
		return err
	}

	err = clone.Run()
	if err != nil {
		out, _ := clone.Out()
		l.New(log.Fields{"output": out}).Errorf(err, "Could not execute git clone")
		return err
	}

	err = os.Rename(downloadPath, path)
	return err
}

// ShowRepos prints information about all repositories.
func (r *Repositories) ShowRepos() {
	for _, repo := range r.Repos {
		fmt.Printf("* Name: %s\n", repo.Name)
		fmt.Printf("  URL: %s\n", repo.URL)
		fmt.Printf("  Complexity: %d\n", repo.Complexity)
	}
}
