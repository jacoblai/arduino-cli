// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
// Copyright 2015 Matthijs Kooijman
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

package test

import (
	"path/filepath"
	"strings"
	"testing"

	paths "github.com/arduino/go-paths-helper"
	bldr "github.com/jacoblai/arduino-cli/arduino/builder"
	"github.com/jacoblai/arduino-cli/arduino/builder/cpp"
	"github.com/jacoblai/arduino-cli/legacy/builder"
	"github.com/jacoblai/arduino-cli/legacy/builder/types"
	"github.com/stretchr/testify/require"
)

func loadPreprocessedSketch(t *testing.T, ctx *types.Context) string {
	res, err := ctx.SketchBuildPath.Join(ctx.Sketch.MainFile.Base() + ".cpp").ReadFile()
	NoError(t, err)
	return string(res)
}

func TestPrototypesAdderBridgeExample(t *testing.T) {
	sketchLocation := paths.New("downloaded_libraries", "Bridge", "examples", "Bridge", "Bridge.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 33 "+quotedSketchLocation+"\nvoid setup();\n#line 46 "+quotedSketchLocation+"\nvoid loop();\n#line 62 "+quotedSketchLocation+"\nvoid process(BridgeClient client);\n#line 82 "+quotedSketchLocation+"\nvoid digitalCommand(BridgeClient client);\n#line 109 "+quotedSketchLocation+"\nvoid analogCommand(BridgeClient client);\n#line 149 "+quotedSketchLocation+"\nvoid modeCommand(BridgeClient client);\n#line 33 "+quotedSketchLocation+"\n")
}

func TestPrototypesAdderSketchWithIfDef(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("SketchWithIfDef", "SketchWithIfDef.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessed := LoadAndInterpolate(t, filepath.Join("SketchWithIfDef", "SketchWithIfDef.preprocessed.txt"), ctx)
	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Equal(t, preprocessed, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderBaladuino(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("Baladuino", "Baladuino.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessed := LoadAndInterpolate(t, filepath.Join("Baladuino", "Baladuino.preprocessed.txt"), ctx)
	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Equal(t, preprocessed, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderCharWithEscapedDoubleQuote(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("CharWithEscapedDoubleQuote", "CharWithEscapedDoubleQuote.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessed := LoadAndInterpolate(t, filepath.Join("CharWithEscapedDoubleQuote", "CharWithEscapedDoubleQuote.preprocessed.txt"), ctx)
	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Equal(t, preprocessed, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderIncludeBetweenMultilineComment(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("IncludeBetweenMultilineComment", "IncludeBetweenMultilineComment.ino"), "arduino:sam:arduino_due_x_dbg")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessed := LoadAndInterpolate(t, filepath.Join("IncludeBetweenMultilineComment", "IncludeBetweenMultilineComment.preprocessed.txt"), ctx)
	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Equal(t, preprocessed, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderLineContinuations(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("LineContinuations", "LineContinuations.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessed := LoadAndInterpolate(t, filepath.Join("LineContinuations", "LineContinuations.preprocessed.txt"), ctx)
	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Equal(t, preprocessed, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderStringWithComment(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("StringWithComment", "StringWithComment.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessed := LoadAndInterpolate(t, filepath.Join("StringWithComment", "StringWithComment.preprocessed.txt"), ctx)
	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Equal(t, preprocessed, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithStruct(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("SketchWithStruct", "SketchWithStruct.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessed := LoadAndInterpolate(t, filepath.Join("SketchWithStruct", "SketchWithStruct.preprocessed.txt"), ctx)
	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	obtained := strings.Replace(preprocessedSketch, "\r\n", "\n", -1)
	// ctags based preprocessing removes the space after "dostuff", but this is still OK
	// TODO: remove this exception when moving to a more powerful parser
	preprocessed = strings.Replace(preprocessed, "void dostuff (A_NEW_TYPE * bar);", "void dostuff(A_NEW_TYPE * bar);", 1)
	obtained = strings.Replace(obtained, "void dostuff (A_NEW_TYPE * bar);", "void dostuff(A_NEW_TYPE * bar);", 1)
	require.Equal(t, preprocessed, obtained)
}

func TestPrototypesAdderSketchWithConfig(t *testing.T) {
	sketchLocation := paths.New("sketch_with_config", "sketch_with_config.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 13 "+quotedSketchLocation+"\nvoid setup();\n#line 17 "+quotedSketchLocation+"\nvoid loop();\n#line 13 "+quotedSketchLocation+"\n")

	preprocessed := LoadAndInterpolate(t, filepath.Join("sketch_with_config", "sketch_with_config.preprocessed.txt"), ctx)
	require.Equal(t, preprocessed, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchNoFunctionsTwoFiles(t *testing.T) {
	sketchLocation := paths.New("sketch_no_functions_two_files", "sketch_no_functions_two_files.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, paths.New("sketch_no_functions_two_files", "sketch_no_functions_two_files.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	mergedSketch := loadPreprocessedSketch(t, ctx)
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Equal(t, mergedSketch, preprocessedSketch) // No prototypes added
}

func TestPrototypesAdderSketchNoFunctions(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("sketch_no_functions", "sketch_no_functions.ino"), "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	sketchLocation := paths.New("sketch_no_functions", "sketch_no_functions.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())
	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	mergedSketch := loadPreprocessedSketch(t, ctx)
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Equal(t, mergedSketch, preprocessedSketch) // No prototypes added
}

func TestPrototypesAdderSketchWithDefaultArgs(t *testing.T) {
	sketchLocation := paths.New("sketch_with_default_args", "sketch_with_default_args.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 4 "+quotedSketchLocation+"\nvoid setup();\n#line 7 "+quotedSketchLocation+"\nvoid loop();\n#line 1 "+quotedSketchLocation+"\n")
}

func TestPrototypesAdderSketchWithInlineFunction(t *testing.T) {
	sketchLocation := paths.New("sketch_with_inline_function", "sketch_with_inline_function.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")

	expected := "#line 1 " + quotedSketchLocation + "\nvoid setup();\n#line 2 " + quotedSketchLocation + "\nvoid loop();\n#line 4 " + quotedSketchLocation + "\nshort unsigned int testInt();\n#line 8 " + quotedSketchLocation + "\nstatic int8_t testInline();\n#line 12 " + quotedSketchLocation + "\n__attribute__((always_inline)) uint8_t testAttribute();\n#line 1 " + quotedSketchLocation + "\n"
	obtained := preprocessedSketch
	// ctags based preprocessing removes "inline" but this is still OK
	// TODO: remove this exception when moving to a more powerful parser
	expected = strings.Replace(expected, "static inline int8_t testInline();", "static int8_t testInline();", -1)
	obtained = strings.Replace(obtained, "static inline int8_t testInline();", "static int8_t testInline();", -1)
	// ctags based preprocessing removes "__attribute__ ....." but this is still OK
	// TODO: remove this exception when moving to a more powerful parser
	expected = strings.Replace(expected, "__attribute__((always_inline)) uint8_t testAttribute();", "uint8_t testAttribute();", -1)
	obtained = strings.Replace(obtained, "__attribute__((always_inline)) uint8_t testAttribute();", "uint8_t testAttribute();", -1)
	require.Contains(t, obtained, expected)
}

func TestPrototypesAdderSketchWithFunctionSignatureInsideIFDEF(t *testing.T) {
	sketchLocation := paths.New("sketch_with_function_signature_inside_ifdef", "sketch_with_function_signature_inside_ifdef.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 1 "+quotedSketchLocation+"\nvoid setup();\n#line 3 "+quotedSketchLocation+"\nvoid loop();\n#line 15 "+quotedSketchLocation+"\nint8_t adalight();\n#line 1 "+quotedSketchLocation+"\n")
}

func TestPrototypesAdderSketchWithUSBCON(t *testing.T) {
	sketchLocation := paths.New("sketch_with_usbcon", "sketch_with_usbcon.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := &types.Context{
		HardwareDirs:         paths.NewPathList(filepath.Join("..", "hardware"), "downloaded_hardware"),
		BuiltInToolsDirs:     paths.NewPathList("downloaded_tools"),
		OtherLibrariesDirs:   paths.NewPathList("libraries"),
		BuiltInLibrariesDirs: paths.New("downloaded_libraries"),
		Verbose:              true,
	}
	ctx = prepareBuilderTestContext(t, ctx, sketchLocation, "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 5 "+quotedSketchLocation+"\nvoid ciao();\n#line 10 "+quotedSketchLocation+"\nvoid setup();\n#line 15 "+quotedSketchLocation+"\nvoid loop();\n#line 5 "+quotedSketchLocation+"\n")
}

func TestPrototypesAdderSketchWithTypename(t *testing.T) {
	sketchLocation := paths.New("sketch_with_typename", "sketch_with_typename.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := &types.Context{
		HardwareDirs:         paths.NewPathList(filepath.Join("..", "hardware"), "downloaded_hardware"),
		BuiltInLibrariesDirs: paths.New("libraries", "downloaded_libraries"),
		BuiltInToolsDirs:     paths.NewPathList("downloaded_tools"),
		Verbose:              true,
	}
	ctx = prepareBuilderTestContext(t, ctx, sketchLocation, "arduino:avr:leonardo")
	defer cleanUpBuilderTestContext(t, ctx)

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	expected := "#line 6 " + quotedSketchLocation + "\nvoid setup();\n#line 10 " + quotedSketchLocation + "\nvoid loop();\n#line 12 " + quotedSketchLocation + "\ntypename Foo<char>::Bar func();\n#line 6 " + quotedSketchLocation + "\n"
	obtained := preprocessedSketch
	// ctags based preprocessing ignores line with typename
	// TODO: remove this exception when moving to a more powerful parser
	expected = strings.Replace(expected, "#line 12 "+quotedSketchLocation+"\ntypename Foo<char>::Bar func();\n", "", -1)
	obtained = strings.Replace(obtained, "#line 12 "+quotedSketchLocation+"\ntypename Foo<char>::Bar func();\n", "", -1)
	require.Contains(t, obtained, expected)
}

func TestPrototypesAdderSketchWithIfDef2(t *testing.T) {
	sketchLocation := paths.New("sketch_with_ifdef", "sketch_with_ifdef.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:yun")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 5 "+quotedSketchLocation+"\nvoid elseBranch();\n#line 9 "+quotedSketchLocation+"\nvoid f1();\n#line 10 "+quotedSketchLocation+"\nvoid f2();\n#line 12 "+quotedSketchLocation+"\nvoid setup();\n#line 14 "+quotedSketchLocation+"\nvoid loop();\n#line 5 "+quotedSketchLocation+"\n")

	expectedSource := LoadAndInterpolate(t, filepath.Join("sketch_with_ifdef", "sketch.preprocessed.txt"), ctx)
	require.Equal(t, expectedSource, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithIfDef2SAM(t *testing.T) {
	sketchLocation := paths.New("sketch_with_ifdef", "sketch_with_ifdef.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:sam:arduino_due_x_dbg")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 2 "+quotedSketchLocation+"\nvoid ifBranch();\n#line 9 "+quotedSketchLocation+"\nvoid f1();\n#line 10 "+quotedSketchLocation+"\nvoid f2();\n#line 12 "+quotedSketchLocation+"\nvoid setup();\n#line 14 "+quotedSketchLocation+"\nvoid loop();\n#line 2 "+quotedSketchLocation+"\n")

	expectedSource := LoadAndInterpolate(t, filepath.Join("sketch_with_ifdef", "sketch.preprocessed.SAM.txt"), ctx)
	require.Equal(t, expectedSource, strings.Replace(preprocessedSketch, "\r\n", "\n", -1))
}

func TestPrototypesAdderSketchWithConst(t *testing.T) {
	sketchLocation := paths.New("sketch_with_const", "sketch_with_const.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:uno")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "#include <Arduino.h>\n#line 1 "+quotedSketchLocation+"\n")
	require.Contains(t, preprocessedSketch, "#line 1 "+quotedSketchLocation+"\nvoid setup();\n#line 2 "+quotedSketchLocation+"\nvoid loop();\n#line 4 "+quotedSketchLocation+"\nconst __FlashStringHelper* test();\n#line 6 "+quotedSketchLocation+"\nconst int test3();\n#line 8 "+quotedSketchLocation+"\nvolatile __FlashStringHelper* test2();\n#line 10 "+quotedSketchLocation+"\nvolatile int test4();\n#line 1 "+quotedSketchLocation+"\n")
}

func TestPrototypesAdderSketchWithDosEol(t *testing.T) {
	ctx := prepareBuilderTestContext(t, nil, paths.New("eol_processing", "eol_processing.ino"), "arduino:avr:uno")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))
	// only requires no error as result
}

func TestPrototypesAdderSketchWithSubstringFunctionMember(t *testing.T) {
	sketchLocation := paths.New("sketch_with_class_and_method_substring", "sketch_with_class_and_method_substring.ino")
	quotedSketchLocation := cpp.QuoteString(Abs(t, sketchLocation).String())

	ctx := prepareBuilderTestContext(t, nil, sketchLocation, "arduino:avr:uno")
	defer cleanUpBuilderTestContext(t, ctx)

	ctx.Verbose = true

	var _err error
	commands := []types.Command{
		&builder.ContainerSetupHardwareToolsLibsSketchAndProps{},
		types.BareCommand(func(ctx *types.Context) error {
			ctx.LineOffset, _err = bldr.PrepareSketchBuildPath(ctx.Sketch, ctx.SourceOverride, ctx.SketchBuildPath)
			return _err
		}),
		&builder.ContainerFindIncludes{},
	}
	for _, command := range commands {
		err := command.Run(ctx)
		NoError(t, err)
	}
	NoError(t, builder.PreprocessSketch(ctx))

	preprocessedSketch := loadPreprocessedSketch(t, ctx)
	require.Contains(t, preprocessedSketch, "class Foo {\nint blooper(int x) { return x+1; }\n};\n\nFoo foo;\n\n#line 7 "+quotedSketchLocation+"\nvoid setup();")
}
