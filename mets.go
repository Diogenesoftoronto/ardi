package main

import (
	"archive/tar"
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	z7 "github.com/bodgit/sevenzip"
)

func CopyMets(path string, dst string) (*os.File, error) {
	var tmp *os.File
	ext := filepath.Ext(path)

	// Need to look for a file that includes a file with a mets.
	// The mets will have a uuid: mets.<uuid>.xml
	// If this is a tar file or zip or 7z we need to extract and read from that directory before going to the data directory
	//TODO: There might be an option to create a unified interface for archives (tar, zip, and 7zip)
	switch ext {
	case ZIP:
		archive, err := zip.OpenReader(path)
		if err != nil {
			return nil, err
		}
		// You can defer and handle and error by wrapping a function in an anonymous function. This way we can have defer blocks!
		defer archive.Close()

		for _, f := range archive.File {
			if strings.Contains(f.Name, "METS") {
				tmp, err = os.Create(filepath.Join(dst, filepath.Base(f.Name)))
				if err != nil {
					return nil, err
				}
				file, err := f.Open()
				if err != nil {
					return tmp, err
				}

				if _, err := io.Copy(tmp, file); err != nil {
					log.Fatal("Could not copy the mets")
				}
				break
			}
		}
	case Z7:
		// A bit of code duplication here, I wonder if this really is the best way
		archive, err := z7.OpenReader(path)
		if err != nil {
			return nil, err
		}
		defer archive.Close()

		for _, f := range archive.File {
			if strings.Contains(f.Name, "METS") {
				tmp, err = os.Create(filepath.Join(dst, filepath.Base(f.Name)))
				if err != nil {
					return nil, err
				}
				file, err := f.Open()
				if err != nil {
					return tmp, err
				}
				if _, err := io.Copy(tmp, file); err != nil {
					log.Fatal("Could not copy the mets")
				}
				break
			}
		}
	case TAR:
		r, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer r.Close()
		archive := tar.NewReader(r)
		for {
			h, err := archive.Next()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			if strings.Contains(h.Name, "mets") {
				tmp, err = os.Create(filepath.Join(dst, filepath.Base(h.Name)))
				if err != nil {
					return nil, err
				}
				if _, err := io.Copy(tmp, archive); err != nil {
					return tmp, err
				}
				break
			}
		}
	default:
		// In the default case it could be that the folder path we were sent was not compressed or that it is in a format that is not recognized.
		log.Fatal("Currently only compressed files are supported")
	}
	// Now that we are done copying all the mets files to the temp directory we can finally work on them!
	return tmp, nil
}
