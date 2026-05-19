package collector

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Entry struct {
	File    string
	Title   string
	URLPath string
}

type Group struct {
	Folder  string
	Entries []Entry
}

func Collect(root string) ([]Group, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	byFolder := map[string][]Entry{}
	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		for _, part := range strings.Split(rel, string(filepath.Separator)) {
			if strings.HasPrefix(part, ".") {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		if d.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") {
			return nil
		}

		folder := filepath.ToSlash(filepath.Dir(rel))
		if folder == "." {
			folder = ""
		}
		urlPath := "/" + filepath.ToSlash(rel)
		byFolder[folder] = append(byFolder[folder], Entry{
			File:    filepath.Base(path),
			Title:   ExtractTitle(path),
			URLPath: urlPath,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	folders := make([]string, 0, len(byFolder))
	for f := range byFolder {
		folders = append(folders, f)
	}
	sort.Strings(folders)

	groups := make([]Group, 0, len(folders))
	for _, f := range folders {
		entries := byFolder[f]
		sort.Slice(entries, func(i, j int) bool { return entries[i].URLPath < entries[j].URLPath })
		groups = append(groups, Group{Folder: f, Entries: entries})
	}
	return groups, nil
}

func ExtractTitle(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return stem(path)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimLeft(line, "# "))
		}
	}
	return stem(path)
}

func stem(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}
