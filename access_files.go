package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Replace this with the path to your folder containing CSV files
	csvFolder := "/Users/ankitbali/Desktop/abc"

	// Create a parent folder for organizing output folders
	outputFolder := filepath.Join(csvFolder, "output")
	err := os.Mkdir(outputFolder, 0755)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating output folder: %v\n", err)
		return
	}

	// List files in the CSV folder
	files, err := os.ReadDir(csvFolder)
	if err != nil {
		fmt.Printf("Error reading CSV folder: %v\n", err)
		return
	}

	// Loop through files and create folders
	for _, file := range files {
		// Check if the file is a CSV file
		if strings.HasSuffix(file.Name(), ".csv") {
			// Extract the file name without extension
			fileName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

			// Create the folder path by joining outputFolder and fileName
			folderPath := filepath.Join(outputFolder, fileName)

			// Create the folder
			err := os.Mkdir(folderPath, 0755)
			if err != nil {
				fmt.Printf("Error creating folder %s: %v\n", fileName, err)
				continue
			}

			// Open the CSV file
			csvFilePath := filepath.Join(csvFolder, file.Name())
			csvFile, err := os.Open(csvFilePath)
			if err != nil {
				fmt.Printf("Error opening CSV file %s: %v\n", file.Name(), err)
				continue
			}
			defer csvFile.Close()

			// Read the CSV file
			reader := csv.NewReader(csvFile)
			var caseNumberIndex int
			var header []string

			// Read the header row to find the "Case Number" column index
			header, err = reader.Read()
			if err != nil {
				fmt.Printf("Error reading header from CSV file %s: %v\n", file.Name(), err)
				continue
			}

			// Find the index of the "Case Number" column
			for i, col := range header {
				if col == "Case Number" {
					caseNumberIndex = i
					break
				}
			}

			// Regular expression to match the format some-number/some-number
			formatRegex := regexp.MustCompile(`^(\d+)/(\d+)$`)

			// Read each row and extract values in the specified format
			line := 2
			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					fmt.Printf("Error reading row %d from CSV file %s: %v\n", line, file.Name(), err)
					break
				}

				// Check if the row has the expected number of fields
				if len(row) != len(header) {
					fmt.Printf("Skipping row %d in CSV file %s: wrong number of fields\n", line, file.Name())
					continue
				}

				// Get the Case Number value from the row
				caseNumberValue := row[caseNumberIndex]

				// Check if the Case Number value matches the expected format
				matches := formatRegex.FindStringSubmatch(caseNumberValue)
				if len(matches) == 3 {
					value1 := matches[1]
					value2 := matches[2]

					// Look for files in three separate folders
					for _, subfolder := range []string{"folderA", "folderB", "folderC"} {
						// Function to recursively search for files with extracted values
						searchFiles := func(dir string) error {
							return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
								if err != nil {
									return err
								}

								if !info.IsDir() {
									// Check if the file name contains both values
									if strings.Contains(info.Name(), value1) && strings.Contains(info.Name(), value2) {
										// Move the file to the created folder
										err := os.Rename(path, filepath.Join(folderPath, info.Name()))
										if err != nil {
											fmt.Printf("Error moving file %s to folder %s: %v\n", info.Name(), fileName, err)
										} else {
											fmt.Printf("Moved file %s to folder %s\n", info.Name(), fileName)
										}
									}
								}
								return nil
							})
						}

						// Start searching for files in the subfolder
						subfolderPath := filepath.Join(csvFolder, subfolder)
						err := searchFiles(subfolderPath)
						if err != nil {
							fmt.Printf("Error searching files in folder %s: %v\n", subfolder, err)
						}
					}
				} else {
					fmt.Printf("Skipping row %d in CSV file %s: invalid Case Number format\n", line, file.Name())
				}

				line++
			}
		}
	}
}
