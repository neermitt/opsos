package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-getter"
	"github.com/neermitt/opsos/pkg/config"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/neermitt/opsos/pkg/utils/fs"
	"github.com/otiai10/copy"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const componentConfigFileName = "component.yaml"

func PrepareComponent(ctx context.Context, componentType string, componentName string) error {
	conf := config.GetConfig(ctx)
	componentPath := GetWorkingDirectory(conf, componentType, componentName)

	componentFs := afero.NewBasePathFs(afero.NewOsFs(), componentPath)

	if exists, err := afero.Exists(componentFs, componentConfigFileName); err != nil {
		return err
	} else if !exists {
		return nil
	}

	componentConfigFileContent, err := afero.ReadFile(componentFs, componentConfigFileName)
	if err != nil {
		return err
	}

	var componentConfig ComponentConfig
	if err = yaml.Unmarshal(componentConfigFileContent, &componentConfig); err != nil {
		return err
	}

	return PrepareComponentBySpec(ctx, componentPath, componentConfig.Spec)
}

func PrepareComponentBySpec(ctx context.Context, componentPath string, spec ComponentSpec) error {

	if err := copyFromSource(ctx, componentPath, spec.Source); err != nil {
		return err
	}
	for _, mixin := range spec.Mixins {
		if err := overrideMixin(ctx, componentPath, componentPath, mixin); err != nil {
			return err
		}
	}

	return nil
}

func copyFromSource(ctx context.Context, destDir string, source VendorComponentSource) error {
	var uri string
	// Parse 'uri' template
	if source.Version != "" {
		t, err := template.New(fmt.Sprintf("source-uri-%s", source.Version)).Parse(source.Uri)
		if err != nil {
			return err
		}

		var tpl bytes.Buffer
		err = t.Execute(&tpl, map[string]string{"Version": source.Version})
		if err != nil {
			return err
		}

		uri = tpl.String()
	} else {
		uri = source.Uri
	}

	matcher := fs.IncludeExcludeMatcher(source.IncludedPaths, source.ExcludedPaths)

	return downloadAndCopy(ctx, uri, ".", destDir, matcher)
}

func overrideMixin(ctx context.Context, componentPath string, destDir string, mixin VendorComponentMixins) error {
	var uri string
	if mixin.Uri == "" {
		return errors.New("'uri' must be specified for each 'mixin' in the 'component.yaml' file")
	}

	if mixin.Filename == "" {
		return errors.New("'filename' must be specified for each 'mixin' in the 'component.yaml' file")
	}

	// Parse 'uri' template
	if mixin.Version != "" {
		t, err := template.New(fmt.Sprintf("mixin-uri-%s", mixin.Version)).Parse(mixin.Uri)
		if err != nil {
			return err
		}

		var tpl bytes.Buffer
		err = t.Execute(&tpl, map[string]string{"Version": mixin.Version})
		if err != nil {
			return err
		}

		uri = tpl.String()
	} else {
		uri = mixin.Uri
	}

	// Check if `uri` is a file path.
	// If it's a file path, check if it's an absolute path.
	// If it's not absolute path, join it with the base path (component dir) and convert to absolute path.
	if absPath, err := utils.JoinAbsolutePathWithPath(componentPath, uri); err == nil {
		uri = absPath
	}

	return downloadAndCopy(ctx, mixin.Uri, mixin.Filename, destDir, fs.NewAllMatcher())
}

func downloadAndCopy(ctx context.Context, url string, subDir string, destDir string, matcher fs.Matcher) error {
	tempDir, err := os.MkdirTemp("", strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		return err
	}

	defer func() {
		if e := os.RemoveAll(tempDir); err == nil {
			err = e
		}
	}()
	// Download the source into the temp folder
	client := &getter.Client{
		Ctx:  ctx,
		Dst:  filepath.Clean(filepath.Join(tempDir, subDir)),
		Src:  url,
		Mode: getter.ClientModeDir,
	}
	if err = client.Get(); err != nil {
		return err
	}

	return copy.Copy(tempDir, destDir, copy.Options{
		PreserveTimes: false,
		PreserveOwner: false,
		Skip: func(src string) (bool, error) {
			if strings.HasSuffix(src, ".git") {
				return true, nil
			}
			trimmedSrc := utils.TrimBasePathFromPath(tempDir+"/", src)
			return !matcher.Match(trimmedSrc), nil
		},
		OnSymlink: func(src string) copy.SymlinkAction {
			return copy.Deep
		},
	})
}