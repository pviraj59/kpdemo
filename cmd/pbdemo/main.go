package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/matthewmcnew/pbdemo/buildpacks"
	"github.com/matthewmcnew/pbdemo/logs"
	"github.com/matthewmcnew/pbdemo/populate"
	"github.com/matthewmcnew/pbdemo/rebase"
	"github.com/matthewmcnew/pbdemo/server"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "A tool to demo build service & kpack",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the Build Service Demo")
	},
}

func main() {
	_ = rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(populateCmd(),
		serveCmd(),
		updateRunImageCmd(),
		cleanupCmd(),
		logsCmd(),
		updateBPCmd(),
	)
}

func populateCmd() *cobra.Command {
	var registry string
	var cacheSize string
	var count int32
	var cmd = &cobra.Command{
		Use:     "populate",
		Aliases: []string{"setup"},
		Short:   "Populate Build Service with Images",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Relocating Buildpacks and Run Image. This will take a moment.")
			imageTag := fmt.Sprintf("%s/pbdemo", registry)

			fmt.Printf("Writing all images to: %s\n", imageTag)

			relocated, err := populate.Relocate(imageTag)
			if err != nil {
				return err
			}

			return populate.Populate(count, relocated.Order, imageTag, cacheSize)
		},
	}
	cmd.Flags().StringVarP(&cacheSize, "cache-size", "s", "500Mi", "the cache size to use for build service images")

	cmd.Flags().StringVarP(&registry, "registry", "r", "", "registry to deploy images into")
	_ = cmd.MarkFlagRequired("registry")

	cmd.Flags().Int32VarP(&count, "count", "c", 0, "the number of images to populate in build service")
	_ = cmd.MarkFlagRequired("count")

	return cmd
}

func updateBPCmd() *cobra.Command {
	var buildpack string
	var cmd = &cobra.Command{
		Use:   "update-buildpack",
		Short: "Create new buildpack to simulate update",
		RunE: func(cmd *cobra.Command, args []string) error {
			return buildpacks.UpdateBuildpack(buildpack)
		},
	}
	cmd.Flags().StringVarP(&buildpack, "buildpack", "b", "", "the id of the buildpack to update")
	_ = cmd.MarkFlagRequired("buildpack")

	return cmd
}

func serveCmd() *cobra.Command {
	var port string
	var cmd = &cobra.Command{
		Use:     "serve",
		Aliases: []string{"visualization", "ui"},
		Short:   "Setup a local web server build service visualization ",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Starting Up")
			go func() {
				time.Sleep(500 * time.Millisecond)

				url := fmt.Sprintf("http://localhost:%s", port)
				fmt.Printf("Open up a browser to %s\n", url)

				server.OpenBrowser(url)
			}()

			server.Serve(port)

			return nil
		},
	}

	cmd.Flags().StringVarP(&port, "port", "p", "8080", "registry to deploy images into")

	return cmd
}

func updateRunImageCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "update-stack",
		Aliases: []string{"rebase", "update-run-image"},
		Short:   "Demo an update by pushing an updated stack run image",
		RunE: func(cmd *cobra.Command, args []string) error {
			return rebase.UpdateRunImage()
		},
	}

	return cmd
}

func cleanupCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "cleanup",
		Short: "Remove build service demo images",
		RunE: func(cmd *cobra.Command, args []string) error {
			return populate.Cleanup()
		},
	}

	return cmd
}

func logsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "logs",
		Short:   "Stream build logs from an image",
		Example: "pbdemo logs <image-name>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("no image name provided")
			}

			image := args[0]

			return logs.Logs(image)
		},
	}

	return cmd
}
