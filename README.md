# 🏘 Lodgify Dumper

Lodgify Dumper is a tiny application thrown together in a few hours for a friend who needed to export a few values from Lodgify and import them into an Excel spreadsheet.
It is by no means robust and should not be used in any production enviornments (why would you need such a specific application anyway!?).


## 🧰 Build

Being that Go compiles for various operating systems painlessly it should be fairly easy to compile and run this on your OS of choice.

#### Windows
```bash
GOOS=windows GOARCH=amd64 go build
```

#### Linux
```bash
GOOS=linux GOARCH=amd64 go build
```

## ⚙ Usage

```bash
./dump_bookings.exe path_to_excel.xlsx sheet_name
```

## 🧪 Tests

Tests are pretty barebones, really only ensures that the happy path works. 

## 🛠 Contributing

If for some reason another human actually finds this useful and wants to contribute, please feel fre to open a PR.