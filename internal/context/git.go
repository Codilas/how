package context

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Codilas/how/pkg/providers"
)

// getGitContext gathers git repository information
func getGitContext() (*providers.GitContext, error) {
	// Check if we're in a git repository
	if !isGitRepository() {
		return nil, nil
	}

	ctx := &providers.GitContext{}

	// Get repository name
	if repo, err := getGitRepository(); err == nil {
		ctx.Repository = repo
	}

	// Get current branch
	if branch, err := getGitBranch(); err == nil {
		ctx.Branch = branch
	}

	// Get current commit hash
	if commit, err := getGitCommit(); err == nil {
		ctx.CommitHash = commit
	}

	// Get git status
	if status, err := getGitStatus(); err == nil {
		ctx.Status = status
	}

	// Get recent commits
	if commits, err := getRecentCommits(5); err == nil {
		ctx.RecentCommits = commits
	}

	// Get remote URL
	if remote, err := getGitRemote(); err == nil {
		ctx.RemoteURL = remote
	}

	return ctx, nil
}

// isGitRepository checks if current directory is in a git repository
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// getGitRepository gets the repository name
func getGitRepository() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	path := strings.TrimSpace(string(output))
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}

	return "", nil
}

// getGitBranch gets the current branch name
func getGitBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// getGitCommit gets the current commit hash
func getGitCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	commit := strings.TrimSpace(string(output))
	if len(commit) > 7 {
		return commit[:7], nil // Return short hash
	}

	return commit, nil
}

// getGitStatus gets the git status
func getGitStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	status := strings.TrimSpace(string(output))
	if status == "" {
		return "clean", nil
	}

	// Count changes
	lines := strings.Split(status, "\n")
	modified := 0
	added := 0
	deleted := 0

	for _, line := range lines {
		if len(line) >= 2 {
			switch line[0] {
			case 'M':
				modified++
			case 'A':
				added++
			case 'D':
				deleted++
			}
		}
	}

	var parts []string
	if modified > 0 {
		parts = append(parts, fmt.Sprintf("%d modified", modified))
	}
	if added > 0 {
		parts = append(parts, fmt.Sprintf("%d added", added))
	}
	if deleted > 0 {
		parts = append(parts, fmt.Sprintf("%d deleted", deleted))
	}

	return strings.Join(parts, ", "), nil
}

// getRecentCommits gets recent commit messages
func getRecentCommits(count int) ([]string, error) {
	cmd := exec.Command("git", "log", "--oneline", fmt.Sprintf("-%d", count))
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var commits []string

	for _, line := range lines {
		if line != "" {
			commits = append(commits, line)
		}
	}

	return commits, nil
}

// getGitRemote gets the remote URL
func getGitRemote() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
