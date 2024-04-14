// Copyright 2024 Jelly Terra
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
)

func main() {
	err := _main()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func _main() error {
	username := flag.String("u", "", "GitHub username")
	flag.Parse()

	if *username == "" {
		println("missing username")
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

	_ = json.Unmarshal(respBody, &repos)

	for _, r := range repos {
		println("Cloning", r.FullName)
		err = exec.Command("git", "clone", r.CloneUrl).Run()
		if err != nil {
			return err
		}
	}

	return nil
}
