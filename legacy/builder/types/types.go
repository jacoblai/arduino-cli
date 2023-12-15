// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package types

import (
	"fmt"

	paths "github.com/arduino/go-paths-helper"
	"github.com/jacoblai/arduino-cli/arduino/libraries"
	"github.com/jacoblai/arduino-cli/arduino/sketch"
)

type SourceFile struct {
	// Sketch or Library pointer that this source file lives in
	Origin interface{}
	// Path to the source file within the sketch/library root folder
	RelativePath *paths.Path
}

func (f *SourceFile) Equals(g *SourceFile) bool {
	return f.Origin == g.Origin &&
		f.RelativePath.EqualsTo(g.RelativePath)
}

// Create a SourceFile containing the given source file path within the
// given origin. The given path can be absolute, or relative within the
// origin's root source folder
func MakeSourceFile(ctx *Context, origin interface{}, path *paths.Path) (*SourceFile, error) {
	if path.IsAbs() {
		var err error
		path, err = sourceRoot(ctx, origin).RelTo(path)
		if err != nil {
			return nil, err
		}
	}
	return &SourceFile{Origin: origin, RelativePath: path}, nil
}

// Return the build root for the given origin, where build products will
// be placed. Any directories inside SourceFile.RelativePath will be
// appended here.
func buildRoot(ctx *Context, origin interface{}) *paths.Path {
	switch o := origin.(type) {
	case *sketch.Sketch:
		return ctx.SketchBuildPath
	case *libraries.Library:
		return ctx.LibrariesBuildPath.Join(o.DirName)
	default:
		panic("Unexpected origin for SourceFile: " + fmt.Sprint(origin))
	}
}

// Return the source root for the given origin, where its source files
// can be found. Prepending this to SourceFile.RelativePath will give
// the full path to that source file.
func sourceRoot(ctx *Context, origin interface{}) *paths.Path {
	switch o := origin.(type) {
	case *sketch.Sketch:
		return ctx.SketchBuildPath
	case *libraries.Library:
		return o.SourceDir
	default:
		panic("Unexpected origin for SourceFile: " + fmt.Sprint(origin))
	}
}

func (f *SourceFile) SourcePath(ctx *Context) *paths.Path {
	return sourceRoot(ctx, f.Origin).JoinPath(f.RelativePath)
}

func (f *SourceFile) ObjectPath(ctx *Context) *paths.Path {
	return buildRoot(ctx, f.Origin).Join(f.RelativePath.String() + ".o")
}

func (f *SourceFile) DepfilePath(ctx *Context) *paths.Path {
	return buildRoot(ctx, f.Origin).Join(f.RelativePath.String() + ".d")
}

type LibraryResolutionResult struct {
	Library          *libraries.Library
	NotUsedLibraries []*libraries.Library
}

type Command interface {
	Run(ctx *Context) error
}

type BareCommand func(ctx *Context) error

func (cmd BareCommand) Run(ctx *Context) error {
	return cmd(ctx)
}
