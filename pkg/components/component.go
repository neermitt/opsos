package components

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-getter"
	v1 "github.com/neermitt/opsos/api/v1"
	"github.com/neermitt/opsos/pkg/utils"
	"github.com/neermitt/opsos/pkg/utils/fs"
	"github.com/otiai10/copy"
	"github.com/spf13/afero"
)

const componentConfigFileName = "component.yaml"

type PrepareComponentOptions struct {
	DryRun bool
}

func PrepareComponent(ctx context.Context, componentPath string, dstDir string, options PrepareComponentOptions) error {
	log.Printf("[INFO] Begins Init component %s", componentPath)

	defer func() {
		log.Printf("[INFO] Ends Init component %s", componentPath)
	}()

	componentFs := afero.NewBasePathFs(afero.NewOsFs(), componentPath)

	if exists, err := afero.Exists(componentFs, componentConfigFileName); err != nil {
		return err
	} else if !exists {
		return nil
	}

	f, err := componentFs.Open(componentConfigFileName)
	if err != nil {
		return err
	}
	component, err := ReadComponent(f)
	if err != nil {
		return err
	}

	err = PrepareComponentBySpec(ctx, componentPath, dstDir, component.Spec, options)
	if err != nil {
		log.Printf("[Error] Init component failed %s, error :%v", componentPath, err)
	}
	return err
}

func PrepareComponentBySpec(ctx context.Context, componentPath string, dstDir string, spec v1.ComponentSpec, options PrepareComponentOptions) error {
	if err := copyFromSource(ctx, dstDir, spec.Source, options); err != nil {
		return err
	}
	for _, mixin := range spec.Mixins {
		if err := overrideMixin(ctx, componentPath, dstDir, mixin, options); err != nil {
			return err
		}
	}

	return nil
}

func copyFromSource(ctx context.Context, destDir string, source v1.ComponentSource, options PrepareComponentOptions) error {
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
	log.Printf("[INFO] Copying from source %s", uri)
	matcher := fs.IncludeExcludeMatcher(source.IncludedPaths, source.ExcludedPaths)

	return downloadAndCopy(ctx, getter.ClientModeDir, uri, ".", destDir, matcher, options)
}

func overrideMixin(ctx context.Context, componentPath string, destDir string, mixin v1.ComponentMixins, options PrepareComponentOptions) error {
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

	return downloadAndCopy(ctx, getter.ClientModeFile, uri, mixin.Filename, destDir, fs.NewAllMatcher(), options)
}

func downloadAndCopy(ctx context.Context, mode getter.ClientMode, url string, subDir string, destDir string, matcher fs.Matcher, options PrepareComponentOptions) error {
	if options.DryRun {
		return nil
	}
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
		Mode: mode,
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
			if utils.IsDir(src) {
				return false, nil
			}
			trimmedSrc := utils.TrimBasePathFromPath(tempDir+"/", src)
			return !matcher.Match(trimmedSrc), nil
		},
		OnSymlink: func(src string) copy.SymlinkAction {
			return copy.Deep
		},
	})
}

func ReadComponent(r io.Reader) (*v1.Component, error) {
	var component v1.Component

	err := utils.DecodeYaml(r, &component)
	if err != nil {
		return nil, err
	}
	err = validateComponent(component)
	if err != nil {
		return nil, err
	}
	return &component, err
}

func validateComponent(component v1.Component) error {
	if component.ApiVersion != "opsos/v1" || component.Kind != "Component" {
		return fmt.Errorf("no resource found of type %s/%s", component.ApiVersion, component.Kind)
	}
	return nil
}
