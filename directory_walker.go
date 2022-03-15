package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type DirectoryWalkState struct {
	remainingTargetPath []string
	pathWalked          string
	pathOptions         []string
}

func (state DirectoryWalkState) DoWalk(folderToWalk string) DirectoryWalkState {
	os.Chdir(folderToWalk)
	return DirectoryWalkState{
		remainingTargetPath: state.remainingTargetPath,
		pathWalked:          filepath.Join(state.pathWalked, folderToWalk),
	}
}

func (state DirectoryWalkState) PrepareWalk() DirectoryWalkState {
	return DirectoryWalkState{
		remainingTargetPath: state.remainingTargetPath[1:],
		pathWalked:          state.pathWalked,
		pathOptions:         getCandidates(state.remainingTargetPath[0]),
	}
}

func InitDirectoryWalk(targetPath string) DirectoryWalkState {
	return DirectoryWalkState{
		remainingTargetPath: getSegments(targetPath),
		pathWalked:          "",
	}
}

func getSegments(path string) []string {
	var folder, file = filepath.Split(path)
	var result []string
	if len(folder) > 0 {
		result = getSegments(strings.TrimRight(folder, "\\/"))
	}
	if len(file) > 0 {
		result = append(result, file)
	}
	return result
}

func getCandidates(folder string) []string {
	if folder == "." {
		return []string{"."}
	}
	if folder == ".." {
		return []string{".."}
	}
	var localFolders, _ = os.ReadDir(".")
	candidates := make([]string, 0)
	var regexMatcher = func(target string) bool {
		match, _ := regexp.Match(folder, []byte(target))
		return match
	}
	for _, entry := range localFolders {
		if entry.Name() == folder {
			return []string{entry.Name()}
		}
		if strings.Contains(entry.Name(), folder) || regexMatcher(entry.Name()) {
			candidates = append(candidates, entry.Name())
		}
	}
	return candidates

}
