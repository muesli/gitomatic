package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var (
	pull     = flag.Bool("pull", true, "automatically pull changes")
	push     = flag.Bool("push", true, "automatically push changes")
	author   = flag.String("author", "gitomatic", "author name for git commits")
	email    = flag.String("email", "gitomatic@fribbledom.com", "email address for git commits")
	interval = flag.String("interval", "1m", "how often to check for changes")
	privkey  = flag.String("privkey", "~/.ssh/id_rsa", "location of private key used for auth")
	username = flag.String("username", "", "username used for auth")
	password = flag.String("password", "", "password used for auth")
)

func fatal(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func gitAdd(w *git.Worktree, path string) error {
	log.Printf("Adding file to work-tree: %s\n", path)
	_, err := w.Add(path)
	return err
}

func gitRemove(w *git.Worktree, path string) error {
	log.Printf("Removing file from work-tree: %s\n", path)
	_, err := w.Remove(path)
	return err
}

func gitCommit(w *git.Worktree, message string) error {
	log.Printf("Creating commit: %s", message)
	_, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  *author,
			Email: *email,
			When:  time.Now(),
		},
	})
	return err
}

func gitPush(r *git.Repository, auth transport.AuthMethod) error {
	if !gitHasRemote(r) {
		log.Println("Not pushing: no remotes configured.")
		return nil
	}

	log.Println("Pushing changes...")
	return r.Push(&git.PushOptions{
		Auth: auth,
	})
}

func gitPull(r *git.Repository, w *git.Worktree, auth transport.AuthMethod) error {
	if !gitHasRemote(r) {
		log.Println("Not pulling: no remotes configured.")
		return nil
	}

	log.Println("Pulling changes...")
	err := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err == transport.ErrEmptyRemoteRepository {
		return nil
	}
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

func gitHasRemote(r *git.Repository) bool {
	remotes, _ := r.Remotes()
	return len(remotes) > 0
}

func parseAuthArgs() (transport.AuthMethod, error) {
	if len(*username) > 0 {
		return &http.BasicAuth{
			Username: *username,
			Password: *password,
		}, nil
	}

	*privkey, _ = homedir.Expand(*privkey)
	auth, err := ssh.NewPublicKeysFromFile("git", *privkey, "")
	if err != nil {
		return nil, err
	}
	return auth, nil
}

func main() {
	fmt.Println("git-o-matic")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("usage: gitomatic <path>")
		os.Exit(1)
	}

	timeout, err := time.ParseDuration(*interval)
	if err != nil {
		fatal("cannot parse interval: %s\n", err)
	}
	auth, err := parseAuthArgs()
	if err != nil {
		fatal("cannot parse key: %s\n", err)
	}

	path := flag.Args()[0]

	for {
		log.Println("Checking repository:", path)
		r, err := git.PlainOpen(path)
		if err != nil {
			fatal("cannot open repository: %s\n", err)
		}
		w, err := r.Worktree()
		if err != nil {
			fatal("cannot access repository: %s\n", err)
		}

		if *pull {
			err = gitPull(r, w, auth)
			if err != nil {
				fatal("cannot pull from repository: %s\n", err)
			}
		}

		if *push {
			status, err := w.Status()
			if err != nil {
				fatal("cannot retrieve git status: %s\n", err)
			}

			changes := 0
			msg := ""
			for path, s := range status {
				switch s.Worktree {
				case git.Untracked:
					log.Printf("New file detected: %s\n", path)
					err := gitAdd(w, path)
					if err != nil {
						fatal("cannot add file: %s\n", err)
					}

					msg += fmt.Sprintf("Add %s.\n", path)
					changes++

				case git.Modified:
					log.Printf("Modified file detected: %s\n", path)
					err := gitAdd(w, path)
					if err != nil {
						fatal("cannot add file: %s\n", err)
					}

					msg += fmt.Sprintf("Update %s.\n", path)
					changes++

				case git.Deleted:
					log.Printf("Deleted file detected: %s\n", path)
					err := gitRemove(w, path)
					if err != nil {
						fatal("cannot remove file: %s\n", err)
					}

					msg += fmt.Sprintf("Remove %s.\n", path)
					changes++

				default:
					log.Printf("%s %s %s\n", string(s.Worktree), string(s.Staging), path)
				}
			}

			if changes == 0 {
				log.Println("No changes detected.")
			} else {
				err = gitCommit(w, msg)
				if err != nil {
					fatal("cannot commit: %s\n", err)
				}
				err = gitPush(r, auth)
				if err != nil {
					fatal("cannot push: %s\n", err)
				}
			}
		}

		log.Printf("Sleeping until next check in %s...\n", timeout)
		time.Sleep(timeout)
	}
}
