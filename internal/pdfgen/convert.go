package pdfgen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertToPDF(inputPath, outputPath string, workerID int) error {
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	profileDir := fmt.Sprintf("/tmp/lo_profile_%d", workerID)

	cmd := exec.Command("soffice",
		"--headless",
		"-env:UserInstallation=file://"+profileDir,
		"--convert-to", "pdf",
		"--outdir", outputPath,
		inputPath,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("converting to pdf: %w", err)
	}

	log.Printf("pdf generated for: %s", filepath.Base(inputPath))
	return nil
}
