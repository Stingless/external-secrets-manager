package main


import (
	"fmt"
	"time"
    "os"
    "external-secrets-manager/esgenerator"

    "github.com/joho/godotenv"
	"gopkg.in/src-d/go-git.v4"
    "gopkg.in/src-d/go-git.v4/config"
    "gopkg.in/src-d/go-git.v4/plumbing"
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
    
    branch := os.Getenv("ES_REPOSITORY_BRANCH")
    if branch == "" {
        fmt.Println("GIT_BRANCH environment variable not set. Defaulting to 'main'.")
        branch = "main" // Default to 'master' branch if not set
    }
    err = w.Checkout(&git.CheckoutOptions{
        Branch: plumbing.ReferenceName("refs/heads/" + branch ),
        Create: true,
        Force:  true,
    })
    if err != nil {
        return fmt.Errorf("Error checking out branch: %v\n", err)
    }

    fmt.Printf("Switched to branch: %s\n", branch)

    esgenerator.EsGenerator()
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
 
        err = repo.Push(&git.PushOptions{
			RemoteName: "origin",
			Auth: &http.BasicAuth{
				Username:  os.Getenv("GIT_USER"), // Replace with your GitHub username or personal access token
                Password: os.Getenv("GIT_PASSWORD"),
				// Replace with your GitHub password or leave empty if using a personal access token
			},
			RefSpecs: []config.RefSpec{config.RefSpec(fmt.Sprintf("refs/heads/"+branch+":refs/heads/"+branch))},
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

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        fmt.Println("Error loading .env file:", err)
    }
    CloneRepository(os.Getenv("ES_GIT_REPOSITORY"),"vault-es")
    CheckForChanges("vault-es")
    fmt.Println("Done !")

}
