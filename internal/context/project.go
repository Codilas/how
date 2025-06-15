package context

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Codilas/how/pkg/providers"
)

// detectProjectType analyzes the current directory to determine project type
func detectProjectType() (*providers.ProjectContext, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Check for different project types
	if ctx := detectNodeJSProject(wd); ctx != nil {
		return ctx, nil
	}

	if ctx := detectPythonProject(wd); ctx != nil {
		return ctx, nil
	}

	if ctx := detectGoProject(wd); ctx != nil {
		return ctx, nil
	}

	if ctx := detectRustProject(wd); ctx != nil {
		return ctx, nil
	}

	if ctx := detectDockerProject(wd); ctx != nil {
		return ctx, nil
	}

	return nil, nil
}

// detectNodeJSProject checks for Node.js projects
func detectNodeJSProject(dir string) *providers.ProjectContext {
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); err != nil {
		return nil
	}

	ctx := &providers.ProjectContext{
		Type: "nodejs",
	}

	// Read package.json
	if data, err := os.ReadFile(packageJSONPath); err == nil {
		var packageJSON struct {
			Name         string            `json:"name"`
			Version      string            `json:"version"`
			Scripts      map[string]string `json:"scripts"`
			Dependencies map[string]string `json:"dependencies"`
		}

		if err := json.Unmarshal(data, &packageJSON); err == nil {
			ctx.Name = packageJSON.Name
			ctx.Version = packageJSON.Version
			ctx.Scripts = packageJSON.Scripts

			// Detect framework
			if _, hasReact := packageJSON.Dependencies["react"]; hasReact {
				ctx.Framework = "react"
			} else if _, hasVue := packageJSON.Dependencies["vue"]; hasVue {
				ctx.Framework = "vue"
			} else if _, hasAngular := packageJSON.Dependencies["@angular/core"]; hasAngular {
				ctx.Framework = "angular"
			} else if _, hasNext := packageJSON.Dependencies["next"]; hasNext {
				ctx.Framework = "next.js"
			} else if _, hasExpress := packageJSON.Dependencies["express"]; hasExpress {
				ctx.Framework = "express"
			}

			// Get dependency list
			for dep := range packageJSON.Dependencies {
				ctx.Dependencies = append(ctx.Dependencies, dep)
			}
		}
	}

	return ctx
}

// detectPythonProject checks for Python projects
func detectPythonProject(dir string) *providers.ProjectContext {
	pythonFiles := []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"}

	var foundFile string
	for _, file := range pythonFiles {
		if _, err := os.Stat(filepath.Join(dir, file)); err == nil {
			foundFile = file
			break
		}
	}

	if foundFile == "" {
		// Check for .py files
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if strings.HasSuffix(entry.Name(), ".py") {
					foundFile = "python_files"
					break
				}
			}
		}
	}

	if foundFile == "" {
		return nil
	}

	ctx := &providers.ProjectContext{
		Type: "python",
		Name: filepath.Base(dir),
	}

	// Detect framework
	if foundFile == "requirements.txt" {
		if deps, err := readRequirementsTxt(filepath.Join(dir, "requirements.txt")); err == nil {
			ctx.Dependencies = deps

			// Detect framework from dependencies
			for _, dep := range deps {
				dep = strings.ToLower(dep)
				if strings.Contains(dep, "django") {
					ctx.Framework = "django"
				} else if strings.Contains(dep, "flask") {
					ctx.Framework = "flask"
				} else if strings.Contains(dep, "fastapi") {
					ctx.Framework = "fastapi"
				}
			}
		}
	}

	return ctx
}

// detectGoProject checks for Go projects
func detectGoProject(dir string) *providers.ProjectContext {
	goModPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err != nil {
		return nil
	}

	ctx := &providers.ProjectContext{
		Type: "go",
	}

	// Read go.mod
	if data, err := os.ReadFile(goModPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "module ") {
				moduleName := strings.TrimPrefix(line, "module ")
				ctx.Name = filepath.Base(moduleName)
				break
			}
		}
	}

	return ctx
}

// detectRustProject checks for Rust projects
func detectRustProject(dir string) *providers.ProjectContext {
	cargoTomlPath := filepath.Join(dir, "Cargo.toml")
	if _, err := os.Stat(cargoTomlPath); err != nil {
		return nil
	}

	ctx := &providers.ProjectContext{
		Type: "rust",
		Name: filepath.Base(dir),
	}

	// Could parse Cargo.toml for more details
	return ctx
}

// detectDockerProject checks for Docker projects
func detectDockerProject(dir string) *providers.ProjectContext {
	dockerfilePath := filepath.Join(dir, "Dockerfile")
	dockerComposePath := filepath.Join(dir, "docker-compose.yml")
	dockerComposeYamlPath := filepath.Join(dir, "docker-compose.yaml")

	hasDockerfile := false
	hasCompose := false

	if _, err := os.Stat(dockerfilePath); err == nil {
		hasDockerfile = true
	}

	if _, err := os.Stat(dockerComposePath); err == nil {
		hasCompose = true
	} else if _, err := os.Stat(dockerComposeYamlPath); err == nil {
		hasCompose = true
	}

	if !hasDockerfile && !hasCompose {
		return nil
	}

	ctx := &providers.ProjectContext{
		Type: "docker",
		Name: filepath.Base(dir),
	}

	if hasCompose {
		ctx.Framework = "docker-compose"
	}

	return ctx
}

// readRequirementsTxt reads Python requirements.txt file
func readRequirementsTxt(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var deps []string
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			// Extract package name (before version specifiers)
			if idx := strings.IndexAny(line, ">=<!="); idx != -1 {
				line = line[:idx]
			}
			deps = append(deps, line)
		}
	}

	return deps, nil
}
