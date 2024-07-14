package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tupagen",
	Short: "Generate a web app usin Tupã framework",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'tupagen create' to generate a new project")
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Run:   runCreate,
}

var runCmd = &cobra.Command{
	Use:   "run [project_name]",
	Short: "Run the specified project",
	Args:  cobra.ExactArgs(1),
	Run:   runProject,
}

func Execute() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(runCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCreate(cmd *cobra.Command, args []string) {
	apiType := promptAPIType()
	if apiType == "REST" {
		restType := promptRESTType()
		if restType == "CRUD" {
			fmt.Println("Generating CRUD REST API")
			port := promptPort()
			projectName := promptProjectName()
			createProjectStructure(projectName)
			generateAPIServer(projectName, port)
			buildProjectBinary(projectName)
			fmt.Printf("Project '%s' created successfully!\n", projectName)
			fmt.Printf("Run the project with './tupagen run %s'\n", projectName)
		}
	}
}

func promptAPIType() string {
	var apiType string
	fmt.Println("Select API Type:")
	fmt.Println("1. REST")
	fmt.Println("Enter choice: ")
	fmt.Scanln(&apiType)
	if strings.TrimSpace(apiType) == "1" {
		return "REST"
	}
	return ""
}

func promptRESTType() string {
	var restType string
	fmt.Println("Select REST API Type:")
	fmt.Println("1. CRUD")
	fmt.Println("Enter choice: ")
	fmt.Scanln(&restType)
	if strings.TrimSpace(restType) == "1" {
		return "CRUD"
	}
	return ""
}

func promptPort() string {
	var port string
	fmt.Println("Enter port number (default 8080): ")
	fmt.Scanln(&port)
	if strings.TrimSpace(port) == "" {
		return ":8080"
	}
	return ":" + port
}

func promptProjectName() string {
	var projectName string
	fmt.Print("Enter project name (default 'default'): ")
	fmt.Scanln(&projectName)
	if strings.TrimSpace(projectName) == "" {
		return "default"
	}
	return projectName
}

func createProjectStructure(projectName string) {
	basePath := filepath.Join("projects", projectName)
	os.MkdirAll(basePath, os.ModePerm)
}

func generateAPIServer(projectName, port string) {
	code := fmt.Sprintf(`package main

	import (
		"github.com/tupatech/tupa"
	)

	func main() {
		server := tupa.NewAPIServer("%s", nil)
		server.New()
	}
	`, port)

	projectPath := filepath.Join("projects", projectName)
	mainFilePath := filepath.Join(projectPath, "main.go")
	os.WriteFile(mainFilePath, []byte(code), 0644)
}

func buildProjectBinary(projectName string) {
	projectPath := filepath.Join("projects", projectName)
	cmd := exec.Command("go", "build", "-o", projectName, ".")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build the project: %v\n", err)
	}
}

func runProject(cmd *cobra.Command, args []string) {
	projectName := args[0]
	projectPath := filepath.Join("projects", projectName, projectName)
	log.Println(projectPath)
	// projectPath := filepath.Join(projectFolder, "main.go")
	if _, err := os.Stat(projectPath); err == nil {
		cmd := exec.Command(projectPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		fmt.Printf("Project '%s' does not exist.\n", projectName)
	}
}
