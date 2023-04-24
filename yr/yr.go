package yr

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/kristianvv/funtemps/conv"
)

const (
	inputFilename  = "kjevik-temp-celsius-20220318-20230318.csv"
	outputFilename = "kjevik-temp-fahr-20220318-20230318.csv"
)

const footerText = "Data is based on validated data (per 18.03.2023)(CC BY 4.0) from Meteorologisk institutt (MET); endret av Kristian Våg"

func ConvTemperature() {
	if _, err := os.Stat(outputFilename); err == nil {
		fmt.Printf("File '%s' already exists. Do you want to regenerate the file? (y/n): ", outputFilename)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			answer := strings.ToLower(scanner.Text())
			if answer == "y" || answer == "yes" {
				break
			} else if answer == "n" || answer == "no" {
				fmt.Println("Exiting...")
				return
			} else {
				fmt.Print("Invalid answer. Do you want to regenerate the file? (y/n): ")
			}
		}
	}

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := csv.NewWriter(bufio.NewWriter(outputFile))
	defer writer.Flush()

	// Write header row
	if scanner.Scan() {
		writer.Write(strings.Split(scanner.Text(), ";"))
	}

	// Convert and write data rows
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ";")
		if len(fields) < 4 {
			continue
		}

		celsius, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}

		fahrenheit := conv.CelsiusToFahrenheit(celsius)
		fields[3] = strconv.FormatFloat(fahrenheit, 'f', 2, 64)
		writer.Write(fields)
	}

	// Write footer row
	writer.Write(strings.Split(footerText, ";"))

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func AverageTemp() float64 {
	file, err := os.Open(inputFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	sum := 0.0
	count := 0.0

	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ";")
		if len(fields) < 4 {
			continue
		}

		temperature, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			continue
		}

		sum += temperature
		count++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	average := sum / count
	fmt.Printf("Average temperature for the period is: %.2f°C\n", average)

	return average
}

func ProcessLine(line string) string {
	fields := strings.Split(line, ";")
	if len(fields) < 4 {
		return ""
	}
	celsius, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return ""
	}
	fahrenheit := conv.CelsiusToFahrenheit(celsius)
	fields[3] = strconv.FormatFloat(fahrenheit, 'f', 2, 64)
	return strings.Join(fields, ";")
}
