package converter

import (
	"backend/src/service/messages"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"mime/multipart"
	"os"
	"regexp"
)

func ConvertFromXLSX(file multipart.File) ([]messages.Message, error) {
	tempFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		return nil, fmt.Errorf("error saving the file: %v", err)
	}

	f, err := excelize.OpenFile(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("error opening the Excel file: %v", err)
	}
	defer f.Close()

	sheet := f.GetSheetList()[0]
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("error reading the Excel file: %v", err)
	}

	pattern, _ := regexp.Compile(`<[^>]*>|&nbsp;`)
	result := make([]messages.Message, len(rows)-1)
	for idx, row := range rows[1:] {
		if row[0] != "" {
			mes := messages.Message{
				UserID:      row[0],
				SubmitDate:  row[1],
				MessageText: DropHTML(row[2], pattern),
			}
			result[idx] = mes
		} else {
			break
		}
	}

	return result, nil
}

func DropHTML(message string, pattern *regexp.Regexp) string {
	return pattern.ReplaceAllString(message, "")
}
