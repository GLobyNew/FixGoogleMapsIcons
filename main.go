package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
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

// Define the order of colors to use
var colorOrder = []string{
	"placemark-blue",
	"placemark-cyan",
	"placemark-teal",
	"placemark-lime",
	"placemark-green",
	"placemark-yellow",
	"placemark-orange",
	"placemark-deeporange",
	"placemark-red",
	"placemark-pink",
	"placemark-purple",
	"placemark-deeppurple",
	"placemark-brown",
	"placemark-gray",
	"placemark-bluegray",
	"placemark-lightblue",
}

// Extract unique color codes from style IDs
func extractColorCodes(styles []Style, styleMaps []StyleMap) []string {
	colorCodes := make(map[string]bool)
	
	// Extract from regular styles
	for _, style := range styles {
		// Extract color code from style ID (e.g., "icon-1602-0288D1-normal" -> "0288D1")
		parts := strings.Split(style.ID, "-")
		for _, part := range parts {
			if len(part) == 6 && isHexColor(part) {
				colorCodes[part] = true
			}
		}
	}
	
	// Extract from style maps
	for _, styleMap := range styleMaps {
		// Extract color code from style map ID
		parts := strings.Split(styleMap.ID, "-")
		for _, part := range parts {
			if len(part) == 6 && isHexColor(part) {
				colorCodes[part] = true
			}
		}
	}
	
	// Convert to slice
	result := make([]string, 0, len(colorCodes))
	for code := range colorCodes {
		result = append(result, code)
	}
	
	// Sort for consistent ordering
	sort.Strings(result)
	return result
}

func isHexColor(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

func createSequentialMapping(colorCodes []string) map[string]string {
	mapping := make(map[string]string)
	
	for i, code := range colorCodes {
		// Use modulo to cycle through colors if we have more codes than colors
		colorIndex := i % len(colorOrder)
		mapping[code] = colorOrder[colorIndex]
	}
	
	return mapping
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

	// Extract unique color codes from the original styles
	colorCodes := extractColorCodes(kml.Document.Styles, kml.Document.StyleMaps)
	
	// Create sequential color mapping
	colorMapping := createSequentialMapping(colorCodes)
	
	fmt.Println("Color mapping:")
	for code, color := range colorMapping {
		fmt.Printf("  %s -> %s\n", code, color)
	}
	
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
