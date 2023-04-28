package main

import (
	"github.com/bep/godartsass"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

func generateCSS() error {
	input, err := os.ReadFile("./static/scss/style.scss")

	if err != nil {
		return errors.Wrap(err, "failed to read scss file")
	}

	currentDir, err := filepath.Abs("./static/scss")

	if err != nil {
		return errors.Wrap(err, "failed to get current directory")
	}

	transpiler, err := godartsass.Start(godartsass.Options{})

	if err != nil {
		return errors.Wrap(err, "failed to start sass transpiler")
	}

	defer func() {
		_ = transpiler.Close()
	}()

	comp, err := transpiler.Execute(godartsass.Args{
		Source:       string(input),
		OutputStyle:  godartsass.OutputStyleCompressed,
		IncludePaths: []string{currentDir},
	})

	if err != nil {
		return errors.Wrap(err, "failed to compile scss")
	}

	if err = os.WriteFile("./static/style.min.css", []byte(comp.CSS), os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to write css file")
	}

	return nil
}

func compileJS(from, to string, mangle bool) error {
	input, err := os.ReadFile(from)

	if err != nil {
		return err
	}

	result := api.Transform(string(input), api.TransformOptions{
		MinifyIdentifiers: mangle,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
	})

	return os.WriteFile(to, result.Code, os.ModePerm)
}

func generateJS() error {
	funcs := []struct {
		from   string
		mangle bool
	}{
		{from: "fts_exec.js", mangle: false},
		{from: "theme_selector.js", mangle: true},
		{from: "toggle_menu.js", mangle: true},
		{from: "owl_wasm.js", mangle: false},
	}

	for _, meta := range funcs {
		minJS := strings.TrimSuffix(meta.from, ".js") + ".min.js"

		if err := compileJS("./static/js/"+meta.from, "./static/"+minJS, meta.mangle); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	var err error

	if err = generateCSS(); err != nil {
		panic(err)
	}

	if err = generateJS(); err != nil {
		panic(err)
	}
}
