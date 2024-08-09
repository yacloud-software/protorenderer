package tests

import (
	"golang.conradwood.net/go-easyops/utils"
)

const (
	fake_git_repo_config = `[core]
        repositoryformatversion = 0
        filemode = true
        bare = false
        logallrefupdates = true
[submodule]
        active = .
[remote "origin"]
        url = https://git.conradwood.net/git/test.git
        fetch = +refs/heads/*:refs/remotes/origin/*
[branch "master"]
        remote = origin
        merge = refs/heads/master

`
)

// turn this directory into a fake git repository
func write_fake_git_repo(dir string) {
	err := utils.WriteFileCreateDir(dir+"/.git/config", []byte(fake_git_repo_config))
	utils.Bail("unable to write file", err)
}
