package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// KML structures
type KML struct {
	XMLName  xml.Name `xml:"kml"`
	Xmlns    string   `xml:"xmlns,attr"`
	Document Document `xml:"Document"`
}

type Document struct {
	Name        string       `xml:"name"`
	Description string       `xml:"description"`
	Styles      []Style      `xml:"Style"`
	StyleMaps   []StyleMap   `xml:"StyleMap"`
	Folders     []Folder     `xml:"Folder"`
}

type Style struct {
	ID        string    `xml:"id,attr"`
	IconStyle IconStyle `xml:"IconStyle"`
}

type IconStyle struct {
	Icon Icon `xml:"Icon"`
}

type Icon struct {
	Href string `xml:"href"`
}

type StyleMap struct {
	ID    string `xml:"id,attr"`
	Pairs []Pair `xml:"Pair"`
}

type Pair struct {
	Key      string `xml:"key"`
	StyleUrl string `xml:"styleUrl"`
}

type Folder struct {
	Name       string      `xml:"name"`
	Placemarks []Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
	StyleUrl    string `xml:"styleUrl"`
	Point       Point  `xml:"Point"`
}

type Point struct {
	Coordinates string `xml:"coordinates"`
}

// Define new styles
var newStyles = []Style{
	{
		ID: "placemark-red",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-red.png"},
		},
	},
	{
		ID: "placemark-blue",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-blue.png"},
		},
	},
	{
		ID: "placemark-purple",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-purple.png"},
		},
	},
	{
		ID: "placemark-yellow",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-yellow.png"},
		},
	},
	{
		ID: "placemark-pink",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-pink.png"},
		},
	},
	{
		ID: "placemark-brown",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-brown.png"},
		},
	},
	{
		ID: "placemark-green",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-green.png"},
		},
	},
	{
		ID: "placemark-orange",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-orange.png"},
		},
	},
	{
		ID: "placemark-deeppurple",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-deeppurple.png"},
		},
	},
	{
		ID: "placemark-lightblue",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-lightblue.png"},
		},
	},
	{
		ID: "placemark-cyan",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-cyan.png"},
		},
	},
	{
		ID: "placemark-teal",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-teal.png"},
		},
	},
	{
		ID: "placemark-lime",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-lime.png"},
		},
	},
	{
		ID: "placemark-deeporange",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-deeporange.png"},
		},
	},
	{
		ID: "placemark-gray",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-gray.png"},
		},
	},
	{
		ID: "placemark-bluegray",
		IconStyle: IconStyle{
			Icon: Icon{Href: "https://omaps.app/placemarks/placemark-bluegray.png"},
		},
	},
}

// Create mapping from old style patterns to new style IDs
func createStyleMapping() map[string]string {
	// Map old style patterns to new styles based on color
	// This mapping is based on the color codes in the original styles
	return map[string]string{
		"0288D1": "placemark-blue",      // Blue
		"0097A7": "placemark-cyan",      // Cyan
		"097138": "placemark-teal",      // Teal/Dark Green
		"558B2F": "placemark-green",     // Green
		"673AB7": "placemark-purple",    // Purple
		"795548": "placemark-brown",     // Brown
		"880E4F": "placemark-deeppurple", // Deep Purple/Pink
		"F9A825": "placemark-yellow",    // Yellow
		"FF5252": "placemark-red",       // Red
	}
}

func findNewStyleForOld(oldStyleID string, colorMapping map[string]string) string {
	// Extract color code from old style ID
	for colorCode, newStyle := range colorMapping {
		if strings.Contains(oldStyleID, colorCode) {
			return newStyle
		}
	}
	// Default to blue if no match found
	return "placemark-blue"
}

func updatePlacemarkStyles(folders []Folder, styleMapping map[string]string) {
	for i := range folders {
		for j := range folders[i].Placemarks {
			oldStyleUrl := folders[i].Placemarks[j].StyleUrl
			// Remove # prefix
			oldStyleID := strings.TrimPrefix(oldStyleUrl, "#")
			
			if newStyleID, exists := styleMapping[oldStyleID]; exists {
				folders[i].Placemarks[j].StyleUrl = "#" + newStyleID
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run kml_style_replacer.go <input.kml> <output.kml>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Read input file
	data, err := ioutil.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse KML
	var kml KML
	err = xml.Unmarshal(data, &kml)
	if err != nil {
		fmt.Printf("Error parsing XML: %v\n", err)
		os.Exit(1)
	}

	// Create color mapping
	colorMapping := createStyleMapping()
	
	// Build style mapping from old StyleMap IDs to new style IDs
	styleMapping := make(map[string]string)
	
	// First, map all StyleMap IDs to their corresponding new styles
	for _, styleMap := range kml.Document.StyleMaps {
		newStyleID := findNewStyleForOld(styleMap.ID, colorMapping)
		styleMapping[styleMap.ID] = newStyleID
	}

	// Replace styles with new ones
	kml.Document.Styles = newStyles
	
	// Remove StyleMaps as they're no longer needed
	kml.Document.StyleMaps = []StyleMap{}

	// Update all placemark style references
	updatePlacemarkStyles(kml.Document.Folders, styleMapping)

	// Marshal back to XML
	output, err := xml.MarshalIndent(kml, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling XML: %v\n", err)
		os.Exit(1)
	}

	// Add XML declaration
	finalOutput := xml.Header + string(output)

	// Write output file
	err = ioutil.WriteFile(outputFile, []byte(finalOutput), 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", inputFile, outputFile)
	fmt.Println("Style mapping summary:")
	for oldStyle, newStyle := range styleMapping {
		fmt.Printf("  %s -> %s\n", oldStyle, newStyle)
	}
}
