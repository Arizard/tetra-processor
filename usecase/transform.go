package usecase

import (
	"github.com/arizard/tetra"
)

// TransformCSV is a use case which takes a Tetra and a Tesseract config,
// transforms a CSV, and saves the result.
func TransformCSV(
	tetraCfgGetter func() tetra.Config,
	stringCSVLoader func() string,
	stringCSVSaver func(string),
) {
	csv := stringCSVLoader()
	tetraConfig := tetraCfgGetter()

	newCSV := tetra.TransformCSV(tetraConfig, csv)

	stringCSVSaver(newCSV)
}
