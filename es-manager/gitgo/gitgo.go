package gitgo


import (
	"fmt"
	"time"

	"gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/config"
    "gopkg.in/src-d/go-git.v4/plumbing/object"
    "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func CloneRepository(url, path string) {
	fmt.Println("Cloning repository...")

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL: url,
	})

	if err != nil {
		fmt.Printf("Error cloning repository: %v\n", err)
        return
	}

	fmt.Printf("Repository cloned successfully to %s\n", path)
}

func CheckForChanges(path string) error {
	fmt.Println("Checking for changes...")

	// Open the repository
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return fmt.Errorf("Error opening repository: %v", err)
	}

	// Fetch the latest changes
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("Error Fetching changes: %v", err)
	}

	// Get the worktree
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("Error getting worktree: %v", err)
	}

	// Add all files
	_, err = w.Add(".")
	if err != nil {
		return fmt.Errorf("Error adding files: %v", err)
	}

	// Check if there are any changes
	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("Error getting status: %v", err)
	}

	// If there are changes, commit and push
	if !status.IsClean() {
		commitMessage := "Adding new secrets"
		_, err := w.Commit(commitMessage, &git.CommitOptions{
			Author: &object.Signature{
				Name:  os.Getenv("GIT_USER"),
				Email: os.Getenv("GIT_EMAIL"),
				When:  time.Now(),
			},
		})
		if err != nil {
			return fmt.Errorf("Error committing changes: %v", err)
		}
        
        err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username:  os.Getenv("GIT_USER"), // Replace with your GitHub username or personal access token
				Password: os.Getenv("GIT_PASSWORD"), // Replace with your GitHub password or leave empty if using a personal access token
			},
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			fmt.Printf("Error pulling changes: %v\n", err)
		}
		
        err = repo.Push(&git.PushOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username:  os.Getenv("GIT_USER"), // Replace with your GitHub username or personal access token
                Password: os.Getenv("GIT_PASSWORD"),
				// Replace with your GitHub password or leave empty if using a personal access token
			},
			RefSpecs: []config.RefSpec{config.RefSpec(fmt.Sprintf("refs/heads/testing-secrets:refs/heads/testing-secrets"))},
		})
		if err != nil {
			return fmt.Errorf("Error pushing changes: %v", err)
		}

		fmt.Println("Changes committed and pushed.")
	} else {
		fmt.Println("No changes detected.")
	}

	return nil
}
