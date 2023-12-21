package outline

import (
	"context"
	"fmt"
	"image"
	"log"
	"os/exec"
	"strconv"
)

type TraceOptions struct {
	Precision int
	Speckle   int
}

func Trace(ctx context.Context, input string, output string, trace_opts *TraceOptions, raster_opts *RasterizeOptions) (image.Image, error) {

	log.Println("VTRACER")

	err := Vtrace(ctx, input, output, trace_opts)

	if err != nil {
		return nil, fmt.Errorf("Failed to run vtracer, %w", err)
	}

	log.Println("RASTERIZE")

	if raster_opts.UseBatik {
		return RasterizeBatik(ctx, raster_opts, output)
	}

	// This is very (very) slow
	return RasterizeNative(ctx, raster_opts, output)
}

func Vtrace(ctx context.Context, input string, output string, opts *TraceOptions) error {

	precision := opts.Precision
	speckle := opts.Speckle

	cmd := "vtracer"

	args := []string{
		"-i",
		input,
		"-o",
		output,
		"--color_precision",
		strconv.Itoa(precision),
		"--filter_speckle",
		strconv.Itoa(speckle),
	}

	return exec.CommandContext(ctx, cmd, args...).Run()
}
