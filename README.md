# X431 Log File to CSV Converter

This program reads LAUNCH diagnostic X431 log files and converts them into a CSV file for easier usage on various other tools.
It converts extracted data into a structured CSV format.

Both `.x431` and `.dzx` files are supported.

> Looking for a desktop version? See [x431-to-excel](https://github.com/ucukertz/x431-to-excel), a GUI app that does the same conversion but outputs `.xlsx` instead of `.csv`.

## Prerequisites
 - Go installed on your machine.
 - An X431 file (`.x431` or `.dzx`) to parse and convert.

## Build
Clone the repository and navigate into the directory:

`git clone https://github.com/ucukertz/x431-to-csv.git && cd x431-to-csv`

Build the program:

`go build -o x431-to-csv`

## Usage
Run the program from your terminal or command prompt with the following format:

`./x431-to-csv <input-file-path>`

Suppose you want to convert `example.x431`. Run the following:

`./x431-to-csv example.x431`

It will generate `example.x431.csv` in the same directory as your input file.

`.dzx` files work the same way:

`./x431-to-csv example.dzx`

This generates `example.dzx.csv` next to the input file.