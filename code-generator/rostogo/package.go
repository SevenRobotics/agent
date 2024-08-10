package rostogo

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bluenviron/goroslib/v2/pkg/conversion"
)

func ImportPackage(pkgDir string, outputDir string) error {
	pkgName := filepath.Base(pkgDir)
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory %s: %v", outputDir, err)
	}
	return conversion.ImportPackage(pkgName, pkgDir, outputDir)
}
