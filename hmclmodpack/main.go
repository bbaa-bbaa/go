package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Manifest struct {
	Name        string          `json:"name"`
	Author      string          `json:"author"`
	Version     string          `json:"version"`
	Description string          `json:"description"`
	FileApi     string          `json:"fileApi"`
	Files       []ManifestFile  `json:"files"`
	Addons      []ManifestAddon `json:"addons"`
}
type ManifestFile struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
}
type ManifestAddon struct {
	Id      string `json:"id"`
	Version string `json:"version"`
}

func main() {
	var dat Manifest
	buf, _ := ioutil.ReadFile("./modpack/server-manifest.json")
	json.Unmarshal(buf, &dat)
	dat.Files = []ManifestFile{}
	dat.Files = gethash(dat.Files)
	outputjson, _ := json.MarshalIndent(dat, "", " ")
	ioutil.WriteFile("./modpack/server-manifest.json", outputjson, 0644)
	fmt.Println("处理完成")
}

func gethash(files []ManifestFile) []ManifestFile {
	path, _ := os.Getwd()
	path += "/modpack/overrides/"
	filepath.WalkDir(path, func(fpath string, d fs.DirEntry, err error) error {
		hash := sha1.New()
		buf, err := ioutil.ReadFile(fpath)
		if err != nil {
			// is Directory
			return nil
		}
		hash.Write(buf)
		shash := fmt.Sprintf("%x", hash.Sum(nil))
		fpath = strings.Replace(fpath, path, "", -1)
		files = append(files, ManifestFile{Path: fpath, Hash: shash})
		fmt.Printf("%s: %s\n", fpath, shash)
		return nil
	})
	return files
}
