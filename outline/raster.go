package outline

import (
	"context"
	"fmt"
	"image"
	"os"
	"os/exec"
	"strings"

	"github.com/sfomuseum/go-svg"
)

type RasterizeOptions struct {
	UseBatik  bool
	PathBatik string
}

func RasterizeNative(ctx context.Context, opts *RasterizeOptions, input string) (image.Image, error) {

	r, err := os.Open(input)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s for reading, %w", input, err)
	}

	defer r.Close()

	return svg.Rasterize(ctx, r)
}

func RasterizeBatik(ctx context.Context, opts *RasterizeOptions, input string) (image.Image, error) {

	cmd := "java"

	args := []string{
		"-jar",
		opts.PathBatik,
		input,
	}

	err := exec.CommandContext(ctx, cmd, args...).Run()

	if err != nil {
		return nil, fmt.Errorf("Failed to run batik, %w", err)
	}

	// Why can't I specify the output filename in Batik?
	output := strings.Replace(input, ".svg", ".png", 1)

	defer os.Remove(output)

	r, err := os.Open(output)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s for reading, %w", output, err)
	}

	defer r.Close()

	im, _, err := image.Decode(r)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode image, %w", err)
	}

	return im, nil
}
