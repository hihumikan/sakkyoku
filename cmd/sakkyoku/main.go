package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/hihumikan/sakkyoku/internal/config"
	"github.com/hihumikan/sakkyoku/internal/notifier"
	"github.com/hihumikan/sakkyoku/internal/updater"
	"github.com/hihumikan/sakkyoku/internal/utils"
	"github.com/spf13/cobra"
)

var cfgPath string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "sakkyoku",
		Short: "Continuous Deployment tool for docker-compose",
		Long:  `sakkyoku is a Continuous Deployment tool for managing docker-compose services.`,
	}

	var installCmd = &cobra.Command{
		Use:   "install",
		Short: "Install sakkyoku",
		Run:   install,
	}

	var uninstallCmd = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall sakkyoku",
		Run:   uninstall,
	}

	var statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Check status of sakkyoku managed projects",
		Run:   status,
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update sakkyoku managed projects",
		Run:   update,
	}

	var cleanupCmd = &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup unused Docker images",
		Run:   cleanup,
	}

	rootCmd.AddCommand(installCmd, uninstallCmd, statusCmd, updateCmd, cleanupCmd)
	rootCmd.Execute()
}

func install(cmd *cobra.Command, args []string) {
	fmt.Println("Installing sakkyoku...")
}

func uninstall(cmd *cobra.Command, args []string) {
	fmt.Println("Uninstalling sakkyoku...")
}

func status(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig("/etc/sakkyoku")
	utils.CheckError(err, "Loading config")

	sn := notifier.NewSlackNotifier(cfg.SlackWebhook, "sakkyoku")
	up := updater.NewUpdater(cfg, sn)
	err = up.DiscoverProjects()
	utils.CheckError(err, "Discovering projects")

	for _, proj := range up.Projects {
		fmt.Printf("Project: %s\n", proj)
	}
}

func update(cmd *cobra.Command, args []string) {
	cfg, err := config.LoadConfig("/etc/sakkyoku")
	utils.CheckError(err, "Loading config")

	sn := notifier.NewSlackNotifier(cfg.SlackWebhook, "sakkyoku")
	up := updater.NewUpdater(cfg, sn)
	err = up.DiscoverProjects()
	utils.CheckError(err, "Discovering projects")

	err = up.UpdateProjects()
	utils.CheckError(err, "Updating projects")
}

func cleanup(cmd *cobra.Command, args []string) {
	fmt.Println("Cleaning up Docker images...")

	pruneCmd := exec.Command("docker", "image", "prune", "-a", "-f")
	pruneCmd.Stdout = os.Stdout
	pruneCmd.Stderr = os.Stderr

	err := pruneCmd.Run()
	if err != nil {
		fmt.Printf("Error during cleanup: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Cleanup completed.")
}
