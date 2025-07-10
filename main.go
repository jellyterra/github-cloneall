// Copyright 2025 Jelly Terra <jellyterra@proton.me>
// Use of this source code form is governed under the MIT License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	err := _main()
	if err != nil {
		fmt.Println(err)
	}
}

func _main() error {
	var (
		cloneBaseUrl = flag.String("h", "https://github.com", "Clone base URL.")
		username     = flag.String("u", "", "GitHub username.")
		deepClone    = flag.Bool("deep", false, "Deeply clone.")

		retryDelay   = flag.Int("d", 3, "Retry delay.")
		retryUplimit = flag.Int("r", 3, "Retry uplimit.")

		gitArgs []string
	)
	flag.Parse()

	if !*deepClone {
		gitArgs = append(gitArgs, "--depth", "1")
	}

	if *username == "" {
		fmt.Println("Missing username.")
		return nil
	}

	resp, err := http.Get("https://api.github.com/users/" + *username + "/repos")
	if err != nil {
		return err
	}

	var repos []Repo

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(respBody, &repos)
	if err != nil {
		return err
	}

	for _, r := range repos {
		fmt.Println(r.FullName)
	}
	reposLen := len(repos)

	for i, r := range repos {
		cloneUrl := *cloneBaseUrl + "/" + r.FullName

		_, err := os.Stat(r.Name)
		if !os.IsNotExist(err) {
			continue
		}

		for retry := 0; retry != *retryUplimit; retry++ {
			fmt.Println("\n>>> [", i+1, "/", reposLen, "] (", retry, "/", *retryUplimit, ") Cloning", cloneUrl)

			cmd := exec.Command("git", append([]string{"clone", cloneUrl}, gitArgs...)...)
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout

			err = cmd.Run()
			if err != nil {
				time.Sleep(time.Duration(*retryDelay) * time.Second)
				continue
			}

			break
		}
	}

	return nil
}
