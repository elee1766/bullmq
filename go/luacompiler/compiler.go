package luacompiler

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"strings"
)

type Workspace struct {
	rawFiles      map[string][]byte
	compiledFiles map[string][]byte
}

func LoadWorkspace(root fs.FS) (*Workspace, error) {
	s := &Workspace{
		compiledFiles: map[string][]byte{},
		rawFiles:      map[string][]byte{},
	}
	err := fs.WalkDir(root, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".lua") {
			return nil
		}
		name := strings.TrimSuffix(path, ".lua")
		contents, err := fs.ReadFile(root, path)
		if err != nil {
			return err
		}
		s.rawFiles[name] = contents
		return nil
	})
	if err != nil {
		return nil, err
	}
	for k := range s.rawFiles {
		ans, err := resolveImports(k, s.rawFiles)
		if err != nil {
			return nil, fmt.Errorf("fail compile %s: %w", k, err)
		}
		s.compiledFiles[k] = ans
	}

	return s, nil
}

func (s *Workspace) Load(name string) ([]byte, error) {
	val, ok := s.compiledFiles[name]
	if !ok {
		return nil, fmt.Errorf("script %s not found", name)
	}
	return val, nil
}

func (s *Workspace) All() map[string][]byte {
	return s.compiledFiles
}

func resolveImports(fileName string, files map[string][]byte) ([]byte, error) {
	var resolved int
	var err error
	step, ok := files[fileName]
	if !ok {
		return nil, fmt.Errorf("file %s not found", fileName)
	}

	for {
		step, resolved, err = resolveImportStep(fileName, step, files)
		if err != nil {
			return nil, err
		}
		if resolved == 0 {
			return step, nil
		}
	}
}

func resolveImportStep(fileName string, fileBytes []byte, files map[string][]byte) ([]byte, int, error) {
	basePath, _ := path.Split(fileName)
	var n int
	newFile := new(bytes.Buffer)
	newFile.Grow(len(fileBytes))
	rd := bufio.NewScanner(bytes.NewBuffer(fileBytes))
	for rd.Scan() {
		orig := rd.Bytes()
		cur := bytes.TrimSpace(orig)
		if bytes.HasPrefix(cur, []byte("--- @include ")) {
			n = n + 1
			importedFile := string(bytes.Trim(bytes.TrimSpace(bytes.TrimPrefix(cur, []byte("--- @include "))), `"`))
			// find the file
			fullImportPath := path.Join(basePath, importedFile)
			importedBytes, err := resolveImports(fullImportPath, files)
			if err != nil {
				return nil, 0, err
			}
			newFile.WriteString(fmt.Sprintf("-- Begin Import %s --\r\n", importedFile))
			newFile.Write(importedBytes)
			newFile.WriteString(fmt.Sprintf("-- End   Import %s --\r\n\r\n", importedFile))
		} else {
			newFile.Write(orig)
			newFile.Write([]byte("\r\n"))
		}
	}
	return newFile.Bytes(), n, nil
}
