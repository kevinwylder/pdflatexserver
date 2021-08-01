package latex

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type SourceDirectory struct {
	sourcePath string
	pdflatex   string
}

const (
	DefaultCompilePath = "/tmp/compile"
	DefaultTexLivePath = "/usr/local/texlive"
)

func NewSourceDirectory(source string) (*SourceDirectory, error) {
	var err error
	source, err = filepath.Abs(source)
	if err != nil {
		return nil, fmt.Errorf("failed to make '%s' absolute: %w", source, err)
	}
	pdflatex, err := exec.LookPath("pdflatex")
	if err != nil {
		log.Println("`pdflatex` not found in path")
		pdflatex, err = lookForTexLive()
		if err != nil {
			return nil, fmt.Errorf("Failed to find `pdflatex`: %w", err)
		}
		log.Printf("Using %s\n", pdflatex)
	}
	d := &SourceDirectory{
		sourcePath: source,
		pdflatex:   pdflatex,
	}
	err = os.MkdirAll(DefaultCompilePath, 0o777)
	if err != nil {
		return nil, fmt.Errorf("Failed to create tmp compile path '%s', %w", DefaultCompilePath, err)
	}
	return d, nil
}

func lookForTexLive() (string, error) {
	machineName, err := getMachineTexName()
	if err != nil {
		return "", fmt.Errorf("Failed to get machine name: %w (is this linux/macos?)", err)
	}
	pattern := fmt.Sprintf("%s/*/bin/%s/pdflatex", DefaultTexLivePath, machineName)
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return "", fmt.Errorf("Failed to get texlive path (is it installed?) %s", pattern)
	}

	sort.Strings(matches)
	return matches[len(matches)-1], nil
}

func getMachineTexName() (string, error) {
	uname, err := exec.Command("uname").Output()
	if err != nil {
		return "", fmt.Errorf("Failed to get machine uname: %w", err)
	}

	unameM, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get uname -m: %w", err)
	}

	return fmt.Sprintf(
		"%s-%s",
		strings.Trim(string(unameM), "\n"),
		strings.Trim(strings.ToLower(string(uname)), "\n"),
	), nil
}

type CompilerError struct {
	bytes.Buffer
}

func (e *CompilerError) Error() string {
	return e.String()
}

func (d *SourceDirectory) PdfCompile(ctx context.Context, tex string) (string, error) {
	compileDir, err := ioutil.TempDir("/tmp/compile", "pdflatex")
	if err != nil {
		return "", fmt.Errorf("Failed to make compile directory: %w", err)
	}
	defer os.RemoveAll(compileDir)

	// compile
	cmd := exec.CommandContext(ctx,
		d.pdflatex,
		"-halt-on-error",
		"-output-directory", compileDir,
		tex,
	)
	var output CompilerError
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		return "", &output
	}

	// move the pdf
	outputDir, _ := path.Split(tex)
	outputs, err := ioutil.ReadDir(compileDir)
	if err != nil {
		return "", fmt.Errorf("Scan output path '%s': %w", compileDir, err)
	}
	for _, output := range outputs {
		name := output.Name()
		if path.Ext(name) != ".pdf" {
			continue
		}

		src := filepath.Join(compileDir, name)
		dst := filepath.Join(outputDir, name)
		os.Remove(dst)
		err = exec.CommandContext(ctx, "mv", src, dst).Run()
		if err != nil {
			return "", fmt.Errorf("Failed to move PDF `mv %s %s`: %w", src, dst, err)
		}
		return dst, nil
	}
	return "", fmt.Errorf("PDF not found in '%s'", compileDir)
}
