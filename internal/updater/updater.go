package updater

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hihumikan/sakkyoku/internal/config"
	"github.com/hihumikan/sakkyoku/internal/notifier"
)

type Updater struct {
	Config   *config.Config
	Notifier *notifier.SlackNotifier
	Projects []string
}

func NewUpdater(cfg *config.Config, notifier *notifier.SlackNotifier) *Updater {
	return &Updater{
		Config:   cfg,
		Notifier: notifier,
	}
}

func (u *Updater) DiscoverProjects() error {
	searchRoot := u.Config.SearchRoot
	if searchRoot == "" {
		return fmt.Errorf("SEARCH_ROOT is not set")
	}

	var projects []string
	err := filepath.Walk(searchRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == ".sakkyoku" {
			proj := filepath.Dir(path)
			projects = append(projects, proj)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking the path: %w", err)
	}

	u.Projects = projects
	return nil
}

func (u *Updater) UpdateProjects() error {
	for _, proj := range u.Projects {
		fmt.Printf("Updating project: %s\n", proj)
		if err := u.UpdateProject(proj); err != nil {
			fmt.Printf("Error updating project %s: %v\n", proj, err)
			if u.Notifier != nil {
				u.Notifier.Notify(fmt.Sprintf("Error updating project %s: %v", proj, err))
			}
		}
	}
	return nil
}

func (u *Updater) UpdateProject(proj string) error {
	// Change to project directory
	cmd := exec.Command("git", "pull", u.Config.RepoGitRemote, "HEAD")
	cmd.Dir = proj
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git pull failed: %v, output: %s", err, out.String())
	}

	// Pull Docker images
	cmd = exec.Command("docker-compose", "pull")
	cmd.Dir = proj
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker-compose pull failed: %v, output: %s", err, out.String())
	}

	// Restart services
	cmd = exec.Command("docker-compose", "up", "-d", "--build")
	cmd.Dir = proj
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("docker-compose up failed: %v, output: %s", err, out.String())
	}

	// Notify success
	if u.Notifier != nil {
		u.Notifier.Notify(fmt.Sprintf("Project %s updated successfully.", proj))
	}

	return nil
}
