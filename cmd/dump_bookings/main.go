package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ishmaelavila/lodgify_dumper/internal/lodgify"
	"github.com/joho/godotenv"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

var sugarLogger *zap.SugaredLogger

func main() {
	err := godotenv.Load()
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugarLogger = logger.Sugar()

	args := os.Args[1:]

	if len(args) < 2 || len(args) > 2 {
		sugarLogger.Fatalf("usage: dump_bookings path_to_file.xlsx sheet_name")
	}

	path := args[0]
	sheetName := args[1]

	if err != nil {
		sugarLogger.Fatal("Error loading .env file")
	}

	bookings := getBookings()
	writeToExcelFile(path, sheetName, bookings)

}

func writeToExcelFile(path string, sheetName string, bookings []lodgify.Booking) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	i := 0
	for i = 0; i < len(bookings); i++ {
		booking := bookings[i]
		row := i + 2
		propertyNameCell := "A" + strconv.Itoa(row)
		checkInDateCell := "B" + strconv.Itoa(row)
		totalAmountcell := "C" + strconv.Itoa(row)
		f.SetCellValue(sheetName, propertyNameCell, booking.PropertyName)
		f.SetCellValue(sheetName, checkInDateCell, booking.Arrival)
		f.SetCellValue(sheetName, totalAmountcell, booking.TotalAmount)
	}

	f.Save()
	sugarLogger.Infof("wrote %d rows to sheet", i)
}

func getBookings() []lodgify.Booking {

	lodgifyAPIKey := os.Getenv("LODGIFY_API_KEY")

	if lodgifyAPIKey == "" {
		sugarLogger.Fatal("LODGIFY_API_KEY must not be empty")
	}
	lodgifyBaseURL := os.Getenv("LODGIFY_BASE_URL")

	if lodgifyBaseURL == "" {
		sugarLogger.Fatal("LODGIFY_BASE_URL must not be empty")
	}

	args := lodgify.LodgifyClientArgs{
		BaseURL: lodgifyBaseURL,
		APIKey:  lodgifyAPIKey,
		Logger:  sugarLogger,
	}

	sugarLogger.Info("creating lodgify client")
	client, err := lodgify.NewClient(args)

	if err != nil {
		sugarLogger.Fatalf("error initializing lodgify client %w", err)
	}

	bookings, err := client.GetBookings()

	if err != nil {
		sugarLogger.Fatalf("error retreiving bookings from lodgify: %v", err)
	}

	return bookings
}
