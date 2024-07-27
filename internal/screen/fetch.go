package screen

import (
	"fmt"
	"log"
	"os"

	"github.com/snusEbjoer/projectManager/internal/config"
)

type Fetcher struct {
	config config.Config
}

func NewFetcher() *Fetcher {
	return &Fetcher{
		config: config.NewConfig(),
	}
}

type Project struct {
	workDir string
	name    string
}

func (f *Fetcher) FetchProjects() []Project {
	res := []Project{}
	for _, wd := range f.config.WorkDirs {
		entry, err := os.ReadDir(wd)
		fmt.Println(res)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, dir := range entry {
			res = append(res, Project{
				workDir: wd,
				name:    dir.Name(),
			})
		}
	}
	fmt.Println(res)
	return res
}
