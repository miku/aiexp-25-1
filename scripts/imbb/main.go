package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogleman/gg"
)

// BoundingBox represents the coordinates from the JSON file
type BoundingBox struct {
	Title string  `json:"Title"`
	XMin  float64 `json:"x_min"`
	YMin  float64 `json:"y_min"`
	XMax  float64 `json:"x_max"`
	YMax  float64 `json:"y_max"`
}

func main() {
	// Parse command line arguments
	imageFile := flag.String("image", "", "Path to the image file (jpg or png)")
	jsonFile := flag.String("json", "", "Path to JSON file with bounding boxes")
	outputFile := flag.String("output", "output.png", "Path for the output image")
	lineWidth := flag.Float64("width", 2.0, "Width of the bounding box lines")
	flag.Parse()

	if *imageFile == "" || *jsonFile == "" {
		fmt.Println("Please provide both image and JSON file paths")
		flag.PrintDefaults()
		return
	}

	// Read the image file
	img, err := readImage(*imageFile)
	if err != nil {
		fmt.Printf("Error reading image: %v\n", err)
		return
	}

	// Read the JSON file
	boxes, err := readBoundingBoxes(*jsonFile)
	if err != nil {
		fmt.Printf("Error reading JSON: %v\n", err)
		return
	}

	// Create a new drawing context
	dc := gg.NewContextForImage(img)

	// Set drawing properties
	dc.SetRGBA(1, 0, 0, 1) // Red color
	dc.SetLineWidth(*lineWidth)

	// Draw bounding boxes
	for _, box := range boxes {
		width := box.XMax - box.XMin
		height := box.YMax - box.YMin

		// Draw rectangle
		dc.DrawRectangle(box.XMin, box.YMin, width, height)
		dc.Stroke()

		// Optional: Draw title above the box
		// Uncomment if you want to display titles
		/*
			dc.SetRGBA(1, 1, 1, 0.7) // White background for text
					dc.DrawRectangle(box.XMin, box.YMin-20, dc.MeasureString(box.Title)+10, 20)
							dc.Fill()
									dc.SetRGBA(0, 0, 0, 1) // Black text
											dc.DrawString(box.Title, box.XMin+5, box.YMin-5)
		*/
	}

	// Save the result
	err = dc.SavePNG(*outputFile)
	if err != nil {
		fmt.Printf("Error saving output image: %v\n", err)
		return
	}

	fmt.Printf("Successfully created %s with %d bounding boxes\n", *outputFile, len(boxes))
}

// readImage reads an image file (JPG or PNG)
func readImage(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the image based on file extension
	var img image.Image
	ext := strings.ToLower(filepath.Ext(filePath))

	if ext == ".jpg" || ext == ".jpeg" {
		img, err = jpeg.Decode(file)
	} else if ext == ".png" {
		img, err = png.Decode(file)
	} else {
		return nil, fmt.Errorf("unsupported image format: %s", ext)
	}

	return img, err
}

// readBoundingBoxes reads the JSON file containing bounding box coordinates
func readBoundingBoxes(filePath string) ([]BoundingBox, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var boxes []BoundingBox
	err = json.Unmarshal(data, &boxes)
	if err != nil {
		return nil, err
	}

	return boxes, nil
}
