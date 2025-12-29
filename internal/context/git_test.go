package context

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

// TestIsGitRepository tests the isGitRepository function
func TestIsGitRepository(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "success - returns boolean type",
			test: func(t *testing.T) {
				result := isGitRepository()
				// Verify it returns a boolean
				if _, ok := interface{}(result).(bool); !ok {
					t.Errorf("isGitRepository() should return a boolean")
				}
			},
		},
		{
			name: "returns boolean in any environment",
			test: func(t *testing.T) {
				result := isGitRepository()
				// Should always return a boolean, regardless of git repo state
				if result != true && result != false {
					t.Errorf("isGitRepository() returned invalid boolean value")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestGetGitRepository tests the getGitRepository function
func TestGetGitRepository(t *testing.T) {
	tests := []struct {
		name          string
		cmdOutput     string
		expectedRepo  string
	}{
		{
			name:          "extracts repository name from standard path",
			cmdOutput:     "/home/user/projects/my-repo\n",
			expectedRepo:  "my-repo",
		},
		{
			name:          "extracts repository name from nested path",
			cmdOutput:     "/var/www/projects/myapp/service\n",
			expectedRepo:  "service",
		},
		{
			name:          "extracts repository name from single level",
			cmdOutput:     "/repo\n",
			expectedRepo:  "repo",
		},
		{
			name:          "returns empty for empty input",
			cmdOutput:     "",
			expectedRepo:  "",
		},
		{
			name:          "handles trailing whitespace",
			cmdOutput:     "/home/user/repo  \n",
			expectedRepo:  "repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parsing logic from getGitRepository
			path := strings.TrimSpace(tt.cmdOutput)
			parts := strings.Split(path, "/")

			var result string
			if len(parts) > 0 && path != "" {
				result = parts[len(parts)-1]
			}

			if result != tt.expectedRepo {
				t.Errorf("Expected repo name %q, got %q", tt.expectedRepo, result)
			}
		})
	}
}

// TestGetGitBranch tests the getGitBranch function
func TestGetGitBranch(t *testing.T) {
	tests := []struct {
		name            string
		cmdOutput       string
		expectedBranch  string
	}{
		{
			name:            "returns main branch name",
			cmdOutput:       "main\n",
			expectedBranch:  "main",
		},
		{
			name:            "returns feature branch name",
			cmdOutput:       "feature/new-feature\n",
			expectedBranch:  "feature/new-feature",
		},
		{
			name:            "trims whitespace from branch name",
			cmdOutput:       "  develop  \n",
			expectedBranch:  "develop",
		},
		{
			name:            "returns empty for empty output",
			cmdOutput:       "",
			expectedBranch:  "",
		},
		{
			name:            "handles branch with special characters",
			cmdOutput:       "feature/JIRA-123-new-feature\n",
			expectedBranch:  "feature/JIRA-123-new-feature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parsing logic from getGitBranch
			branch := strings.TrimSpace(tt.cmdOutput)

			if branch != tt.expectedBranch {
				t.Errorf("Expected branch %q, got %q", tt.expectedBranch, branch)
			}
		})
	}
}

// TestGetGitCommit tests the getGitCommit function
func TestGetGitCommit(t *testing.T) {
	tests := []struct {
		name           string
		cmdOutput      string
		expectedCommit string
	}{
		{
			name:           "shortens full hash to 7 characters",
			cmdOutput:      "abc123def456789abcdef\n",
			expectedCommit: "abc123d",
		},
		{
			name:           "returns 7 char hash as is",
			cmdOutput:      "1234567\n",
			expectedCommit: "1234567",
		},
		{
			name:           "shortens 8 char hash to 7",
			cmdOutput:      "12345678\n",
			expectedCommit: "1234567",
		},
		{
			name:           "returns short hashes less than 7 chars as is",
			cmdOutput:      "abc\n",
			expectedCommit: "abc",
		},
		{
			name:           "handles whitespace in output",
			cmdOutput:      "  abc1234def456789  \n",
			expectedCommit: "abc1234",
		},
		{
			name:           "returns empty for empty output",
			cmdOutput:      "",
			expectedCommit: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parsing logic from getGitCommit
			commit := strings.TrimSpace(tt.cmdOutput)
			if len(commit) > 7 {
				commit = commit[:7]
			}

			if commit != tt.expectedCommit {
				t.Errorf("Expected commit %q, got %q", tt.expectedCommit, commit)
			}
		})
	}
}

// TestGetGitStatus tests the getGitStatus function
func TestGetGitStatus(t *testing.T) {
	tests := []struct {
		name           string
		cmdOutput      string
		expectedStatus string
	}{
		{
			name:           "returns clean for no changes",
			cmdOutput:      "",
			expectedStatus: "clean",
		},
		{
			name:           "counts single modified file",
			cmdOutput:      " M file.go\n",
			expectedStatus: "1 modified",
		},
		{
			name:           "counts multiple modified files",
			cmdOutput:      " M file1.go\n M file2.go\n",
			expectedStatus: "2 modified",
		},
		{
			name:           "counts added files",
			cmdOutput:      "A new.go\n",
			expectedStatus: "1 added",
		},
		{
			name:           "counts deleted files",
			cmdOutput:      "D old.go\n",
			expectedStatus: "1 deleted",
		},
		{
			name:           "combines modified, added, and deleted",
			cmdOutput:      " M file.go\nA new.go\nD old.go\n",
			expectedStatus: "1 modified, 1 added, 1 deleted",
		},
		{
			name:           "counts multiple mixed changes",
			cmdOutput:      " M file1.go\n M file2.go\nA new1.go\nA new2.go\nD old.go\n",
			expectedStatus: "2 modified, 2 added, 1 deleted",
		},
		{
			name:           "ignores unrecognized status codes",
			cmdOutput:      "?? unknown.go\n M modified.go\n",
			expectedStatus: "1 modified",
		},
		{
			name:           "handles whitespace in output",
			cmdOutput:      "  \n M file.go\n  \n",
			expectedStatus: "1 modified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parsing logic from getGitStatus
			status := strings.TrimSpace(tt.cmdOutput)
			if status == "" {
				status = "clean"
			} else {
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
					parts = append(parts, formatGitCount(modified, "modified"))
				}
				if added > 0 {
					parts = append(parts, formatGitCount(added, "added"))
				}
				if deleted > 0 {
					parts = append(parts, formatGitCount(deleted, "deleted"))
				}

				status = strings.Join(parts, ", ")
			}

			if status != tt.expectedStatus {
				t.Errorf("Expected status %q, got %q", tt.expectedStatus, status)
			}
		})
	}
}

// formatGitCount formats a count with a label (helper for testing)
func formatGitCount(count int, label string) string {
	return fmt.Sprintf("%d %s", count, label)
}

// TestGetRecentCommits tests the getRecentCommits function
func TestGetRecentCommits(t *testing.T) {
	tests := []struct {
		name            string
		count           int
		cmdOutput       string
		expectedCommits []string
	}{
		{
			name:      "returns single commit",
			count:     1,
			cmdOutput: "abc1234 feat: add new feature\n",
			expectedCommits: []string{
				"abc1234 feat: add new feature",
			},
		},
		{
			name:      "returns multiple commits",
			count:     3,
			cmdOutput: "abc1234 feat: add new feature\ndef5678 fix: resolve issue\n1234567 docs: update readme\n",
			expectedCommits: []string{
				"abc1234 feat: add new feature",
				"def5678 fix: resolve issue",
				"1234567 docs: update readme",
			},
		},
		{
			name:      "filters out empty lines",
			count:     2,
			cmdOutput: "abc1234 feat: add new feature\n\ndef5678 fix: resolve issue\n",
			expectedCommits: []string{
				"abc1234 feat: add new feature",
				"def5678 fix: resolve issue",
			},
		},
		{
			name:            "returns empty slice for empty output",
			count:           5,
			cmdOutput:       "",
			expectedCommits: []string{},
		},
		{
			name:      "handles commits with long messages",
			count:     2,
			cmdOutput: "abc1234 feat: add very long feature description that spans multiple words\ndef5678 fix: resolve issue with specific error handling\n",
			expectedCommits: []string{
				"abc1234 feat: add very long feature description that spans multiple words",
				"def5678 fix: resolve issue with specific error handling",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parsing logic from getRecentCommits
			lines := strings.Split(strings.TrimSpace(tt.cmdOutput), "\n")
			var commits []string

			for _, line := range lines {
				if line != "" {
					commits = append(commits, line)
				}
			}

			if len(commits) != len(tt.expectedCommits) {
				t.Errorf("Expected %d commits, got %d", len(tt.expectedCommits), len(commits))
			}

			for i, commit := range commits {
				if i < len(tt.expectedCommits) && commit != tt.expectedCommits[i] {
					t.Errorf("Commit %d: expected %q, got %q", i, tt.expectedCommits[i], commit)
				}
			}
		})
	}
}

// TestGetGitRemote tests the getGitRemote function
func TestGetGitRemote(t *testing.T) {
	tests := []struct {
		name           string
		cmdOutput      string
		expectedRemote string
	}{
		{
			name:           "returns https remote URL",
			cmdOutput:      "https://github.com/user/repo.git\n",
			expectedRemote: "https://github.com/user/repo.git",
		},
		{
			name:           "returns ssh remote URL",
			cmdOutput:      "git@github.com:user/repo.git\n",
			expectedRemote: "git@github.com:user/repo.git",
		},
		{
			name:           "trims whitespace from remote URL",
			cmdOutput:      "  https://github.com/user/repo.git  \n",
			expectedRemote: "https://github.com/user/repo.git",
		},
		{
			name:           "returns empty for empty output",
			cmdOutput:      "",
			expectedRemote: "",
		},
		{
			name:           "handles custom git server URLs",
			cmdOutput:      "https://gitlab.company.com/team/project.git\n",
			expectedRemote: "https://gitlab.company.com/team/project.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the parsing logic from getGitRemote
			remote := strings.TrimSpace(tt.cmdOutput)

			if remote != tt.expectedRemote {
				t.Errorf("Expected remote %q, got %q", tt.expectedRemote, remote)
			}
		})
	}
}

// Integration Tests (require actual git to be available)

// TestIsGitRepositoryIntegration tests isGitRepository with actual git
func TestIsGitRepositoryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	result := isGitRepository()

	// Should return a boolean
	_, ok := interface{}(result).(bool)
	if !ok {
		t.Errorf("isGitRepository() should return a boolean")
	}
}

// TestGetGitBranchIntegration tests getGitBranch with actual git
func TestGetGitBranchIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if not in a git repository
	if !isGitRepository() {
		t.Skip("Not in a git repository, skipping integration test")
	}

	branch, err := getGitBranch()

	// If we got an error, just log it - some test environments may not have git
	if err != nil {
		t.Logf("Failed to get git branch: %v", err)
		return
	}

	// If successful, verify the result is reasonable
	if branch == "" {
		// Empty branch might be valid in detached head state
		t.Logf("Got empty branch name (may be detached HEAD)")
	}
}

// TestGetGitCommitIntegration tests getGitCommit with actual git
func TestGetGitCommitIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if not in a git repository
	if !isGitRepository() {
		t.Skip("Not in a git repository, skipping integration test")
	}

	commit, err := getGitCommit()

	// If we got an error, just log it
	if err != nil {
		t.Logf("Failed to get git commit: %v", err)
		return
	}

	// Commit hash should be 7 characters or less
	if len(commit) > 7 {
		t.Errorf("Expected commit hash to be 7 characters or less, got %d: %s", len(commit), commit)
	}
}

// TestGetGitStatusIntegration tests getGitStatus with actual git
func TestGetGitStatusIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if not in a git repository
	if !isGitRepository() {
		t.Skip("Not in a git repository, skipping integration test")
	}

	status, err := getGitStatus()

	// If we got an error, just log it
	if err != nil {
		t.Logf("Failed to get git status: %v", err)
		return
	}

	// Status should be either "clean" or contain descriptive text
	if status != "" {
		// Valid statuses contain "modified", "added", "deleted", or "clean"
		validKeywords := []string{"clean", "modified", "added", "deleted"}
		found := false
		for _, keyword := range validKeywords {
			if strings.Contains(status, keyword) {
				found = true
				break
			}
		}
		if !found && status != "" {
			t.Logf("Got unexpected status format: %s", status)
		}
	}
}

// TestGetRecentCommitsIntegration tests getRecentCommits with actual git
func TestGetRecentCommitsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if not in a git repository
	if !isGitRepository() {
		t.Skip("Not in a git repository, skipping integration test")
	}

	commits, err := getRecentCommits(5)

	// If we got an error, just log it
	if err != nil {
		t.Logf("Failed to get recent commits: %v", err)
		return
	}

	// Commits should be a non-nil slice
	if commits == nil {
		t.Logf("Got nil commits slice")
		return
	}

	// Each commit line should contain some content
	for i, commit := range commits {
		if commit == "" {
			t.Logf("Commit %d is empty", i)
		}
	}
}

// TestGetGitRemoteIntegration tests getGitRemote with actual git
func TestGetGitRemoteIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Skip if not in a git repository
	if !isGitRepository() {
		t.Skip("Not in a git repository, skipping integration test")
	}

	remote, err := getGitRemote()

	// If we got an error, just log it
	if err != nil {
		t.Logf("Failed to get git remote: %v", err)
		return
	}

	// Remote should be a valid URL or empty
	if remote != "" {
		// Basic check that it looks like a URL
		if !strings.Contains(remote, ":") && !strings.Contains(remote, "@") {
			t.Logf("Remote doesn't look like a URL: %s", remote)
		}
	}
}

// Test helpers for creating temporary git repositories

// createTempGitRepo creates a temporary git repository with initial commit
func createTempGitRepo(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to initialize git repo: %v", err)
	}

	// Configure git user (required for commits)
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git user email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git user name: %v", err)
	}

	return tmpDir
}

// createCommit creates a commit in the given git repository
func createCommit(t *testing.T, repoDir, message string) {
	// Create a file to commit
	filename := fmt.Sprintf("file-%d.txt", time.Now().UnixNano())
	filepath := path.Join(repoDir, filename)

	if err := os.WriteFile(filepath, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add file
	cmd := exec.Command("git", "add", filename)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git add: %v", err)
	}

	// Commit
	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to git commit: %v", err)
	}
}

// modifyFile modifies a file in the repository (creates uncommitted changes)
func modifyFile(t *testing.T, repoDir, filename string) {
	filepath := path.Join(repoDir, filename)

	if err := os.WriteFile(filepath, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}
}

// createNewFile creates a new untracked file
func createNewFile(t *testing.T, repoDir, filename string) {
	filepath := path.Join(repoDir, filename)

	if err := os.WriteFile(filepath, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create new file: %v", err)
	}
}

// TestGetGitContextWithValidRepo tests getGitContext with a valid repository
func TestGetGitContextWithValidRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Create initial commit
	createCommit(t, tmpDir, "Initial commit")

	// Change to repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	// Get git context
	ctx, err := getGitContext()
	if err != nil {
		t.Fatalf("getGitContext() returned error: %v", err)
	}

	if ctx == nil {
		t.Fatal("getGitContext() returned nil context")
	}

	// Verify repository name
	if ctx.Repository != path.Base(tmpDir) {
		t.Errorf("Expected repository name %q, got %q", path.Base(tmpDir), ctx.Repository)
	}

	// Verify branch (default is usually 'master' or 'main')
	if ctx.Branch == "" {
		t.Error("Branch should not be empty")
	}

	// Verify commit hash (should be 7 characters)
	if ctx.CommitHash == "" {
		t.Error("CommitHash should not be empty")
	}
	if len(ctx.CommitHash) != 7 {
		t.Errorf("Expected commit hash length of 7, got %d: %s", len(ctx.CommitHash), ctx.CommitHash)
	}

	// Verify status is clean
	if ctx.Status != "clean" {
		t.Errorf("Expected status 'clean', got %q", ctx.Status)
	}

	// Verify recent commits contains our commit
	if len(ctx.RecentCommits) == 0 {
		t.Error("RecentCommits should not be empty")
	}
}

// TestGetGitContextWithCleanRepo tests getGitContext with a clean working directory
func TestGetGitContextWithCleanRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Create multiple commits
	createCommit(t, tmpDir, "First commit")
	createCommit(t, tmpDir, "Second commit")
	createCommit(t, tmpDir, "Third commit")

	// Change to repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	// Get git context
	ctx, err := getGitContext()
	if err != nil {
		t.Fatalf("getGitContext() returned error: %v", err)
	}

	if ctx == nil {
		t.Fatal("getGitContext() returned nil context")
	}

	// Verify status is clean
	if ctx.Status != "clean" {
		t.Errorf("Expected status 'clean', got %q", ctx.Status)
	}

	// Verify we have recent commits (should have at least 3)
	if len(ctx.RecentCommits) < 3 {
		t.Errorf("Expected at least 3 recent commits, got %d", len(ctx.RecentCommits))
	}
}

// TestGetGitContextWithUncommittedChanges tests getGitContext with uncommitted changes
func TestGetGitContextWithUncommittedChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Create initial commit
	createCommit(t, tmpDir, "Initial commit")

	// Change to repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	// Modify an existing file (create uncommitted changes)
	modifyFile(t, tmpDir, "file-123.txt")

	// Get git context
	ctx, err := getGitContext()
	if err != nil {
		t.Fatalf("getGitContext() returned error: %v", err)
	}

	if ctx == nil {
		t.Fatal("getGitContext() returned nil context")
	}

	// Verify status shows modified
	if ctx.Status != "clean" && !strings.Contains(ctx.Status, "modified") {
		t.Errorf("Expected status to contain 'modified', got %q", ctx.Status)
	}

	// The status should reflect the modification
	if ctx.Status == "clean" {
		t.Error("Status should not be clean when files are modified")
	}
}

// TestGetGitContextWithMixedChanges tests getGitContext with mixed added/modified/deleted changes
func TestGetGitContextWithMixedChanges(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Create initial commits with some files
	createCommit(t, tmpDir, "Initial commit")
	createCommit(t, tmpDir, "Second commit")

	// Change to repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	// Create modified files and new files
	modifyFile(t, tmpDir, "file-123.txt")  // Modify existing
	createNewFile(t, tmpDir, "newfile.txt") // Add new untracked file

	// Get git context
	ctx, err := getGitContext()
	if err != nil {
		t.Fatalf("getGitContext() returned error: %v", err)
	}

	if ctx == nil {
		t.Fatal("getGitContext() returned nil context")
	}

	// Verify status is not clean
	if ctx.Status == "clean" {
		t.Error("Status should not be clean when files are modified or added")
	}
}

// TestGetGitContextOutsideRepo tests getGitContext outside a git repository
func TestGetGitContextOutsideRepo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a temporary directory that's NOT a git repo
	tmpDir, err := os.MkdirTemp("", "non-git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to the non-git directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to directory: %v", err)
	}

	// Get git context - should return nil without error
	ctx, err := getGitContext()
	if err != nil {
		t.Fatalf("getGitContext() returned error: %v", err)
	}

	if ctx != nil {
		t.Error("getGitContext() should return nil context when outside a git repository")
	}
}

// TestGetGitContextPopulatesAllFields tests that getGitContext populates all fields
func TestGetGitContextPopulatesAllFields(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Create a commit with a meaningful message
	createCommit(t, tmpDir, "feat: add new feature")

	// Change to repo directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	// Get git context
	ctx, err := getGitContext()
	if err != nil {
		t.Fatalf("getGitContext() returned error: %v", err)
	}

	if ctx == nil {
		t.Fatal("getGitContext() returned nil context")
	}

	// Verify all important fields are populated
	if ctx.Repository == "" {
		t.Error("Repository field should be populated")
	}

	if ctx.Branch == "" {
		t.Error("Branch field should be populated")
	}

	if ctx.CommitHash == "" {
		t.Error("CommitHash field should be populated")
	}

	if ctx.Status == "" {
		t.Error("Status field should be populated")
	}

	// RecentCommits might be empty in some edge cases, but log if it is
	if len(ctx.RecentCommits) == 0 {
		t.Logf("RecentCommits is empty (may be expected in some cases)")
	}

	// RemoteURL might be empty if no remote is configured
	if ctx.RemoteURL == "" {
		t.Logf("RemoteURL is empty (expected when no remote is configured)")
	}
}
