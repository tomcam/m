package app

import (
	"embed"
	"fmt"
	"github.com/tomcam/m/pkg/app"
	"github.com/tomcam/m/pkg/default"
	"io/fs"
	"os"
	"path/filepath"
)

// The following embeds all files and subdirectories
// from the themes subdirectory of this package into
// the executable. So have the subdirectory available
// at compile time. Then you can run the finished
// executable anywhere and it will display the
// list of files even though the themes directory
// doesn't exist at runtime.

//go:embed themes/*
var themeFiles embed.FS

// embedListDir() displays the filenames in the embedded
// directory named theme.
func embedListDir(files embed.FS) error {
	fs.WalkDir(files, ".", func(srcFilename string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(srcFilename)
		return nil
	})
	return nil
}

//func cp(d fs.DirEntry, source string) {
func cp(stat fs.FileInfo, source string) {
	//return nil

}

// embedDirCpy() copies the files specified  in embed.FS to
// the directory specified in path.
func embedDirCpy(files embed.FS, target string) error {
	fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// TODO: Improve error handling
			return err
		}
		if d.IsDir() {
			//fmt.Printf("%s <dir>", d.Name())
			if path != "." {
				// Todo: compute this using config info
				//path  = filepath.Join("/Users/tom/code/foo/bar", path)
				target := filepath.Join("/Users/tom/code/foo/bar", path)
				err := os.MkdirAll(path, defaults.PublicFilePermissions)
				if err != nil {
					// TODO: Improve error handling
					return err
				}
			}

			return nil
		}
		stat, err := os.Stat(path)
		if !stat.Mode().IsRegular() {
			// TODO: Proper error handling
			return fmt.Errorf("%s can't be copied (error: %v)", path, stat)
		}
		//cp(stat, path)
		return nil
	})
	return nil
}
