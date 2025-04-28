package main

import (
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Ensure the user provided a file path argument
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <input-file-path>", filepath.Base(os.Args[0]))
	}
	inputFile := os.Args[1]
	outputFile := inputFile + ".csv"

	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open file %s: %v", inputFile, err)
	}
	defer f.Close()

	buffer := make([]byte, 128)

	// Read channel count
	f.Seek(0x134, io.SeekStart)
	f.Read(buffer[:1])
	columnCount := int(buffer[0]) / 4

	// Skip to first length field
	f.Seek(0x0c, io.SeekStart)
	f.Read(buffer[:4])
	var32 := int32(binary.LittleEndian.Uint32(buffer[:4]))
	f.Seek(int64(var32), io.SeekCurrent)

	// Skip headers
	for i := 0; i < 8; i++ {
		f.Read(buffer[:2])
		var16 := int(binary.LittleEndian.Uint16(buffer[:2]))
		f.Seek(int64(var16-2), io.SeekCurrent)
	}

	getFileSize := func(f *os.File) int64 {
		stat, _ := f.Stat()
		return stat.Size()
	}

	// Read point values
	var pointValues []string
	for {
		if pos, _ := f.Seek(0, io.SeekCurrent); pos == getFileSize(f) {
			break
		}

		if _, err := f.Read(buffer[:2]); err != nil {
			break
		}
		var16 := int(binary.LittleEndian.Uint16(buffer[:2]))

		if _, err := f.Read(buffer[:var16-2]); err != nil {
			break
		}

		value := string(buffer[:var16-3])
		pointValues = append(pointValues, value)
	}

	// Read column headers
	columnNames := make([]string, columnCount+1)
	columnNames[0] = "Num"

	f.Seek(0x138, io.SeekStart)
	for i := 0; i < columnCount; i++ {
		f.Read(buffer[:4])
		index := int(binary.LittleEndian.Uint16(buffer[:2]))
		if index != 0 && index-0x09 < len(pointValues) {
			columnNames[i+1] = fmt.Sprintf("%d. %s", i+1, pointValues[index-0x09])
		}
	}

	for i := 0; i < columnCount; i++ {
		f.Read(buffer[:4])
		index := int(binary.LittleEndian.Uint16(buffer[:2]))
		if index != 0 && index-0x09 < len(pointValues) {
			columnNames[i+1] = fmt.Sprintf("%s (%s)", columnNames[i+1], pointValues[index-0x09])
		}
	}

	// Prepare to read data
	f.Seek(0x11c, io.SeekStart)
	f.Read(buffer[:2])
	var16 := int(binary.LittleEndian.Uint16(buffer[:2]))

	f.Seek(int64(var16+8), io.SeekStart)
	f.Read(buffer[:8])
	recordsCount := int(binary.LittleEndian.Uint32(buffer[:4]))

	totalRows := (recordsCount / 4) / columnCount
	var rows [][]string

	for i := 0; i < totalRows; i++ {
		row := make([]string, columnCount+1)
		row[0] = fmt.Sprintf("%d", i+1)

		for j := 0; j < columnCount; j++ {
			f.Read(buffer[:4])
			index := int(binary.LittleEndian.Uint16(buffer[:2])) - 0x09
			if index >= 0 && index < len(pointValues) {
				row[j+1] = pointValues[index]
			} else {
				row[j+1] = "0"
			}
		}
		rows = append(rows, row)
	}

	// Write to CSV
	outFile, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", outputFile, err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	writer.Write(columnNames)
	writer.WriteAll(rows)

	fmt.Printf("CSV export complete: %s\n", filepath.Base(outputFile))
}
