# go-dcc
EU Digital Covid Certificate utilities in Go [Create, Validate and Parse EU-DCC/Greenpass] 

*Repo is a work in-progress*

-----

### Purpose:

- To understand how Vaccine Passports / EU-DCC / EU-DGC / Greenpasses work.
- Analyze what personal information is stored within.
- Understand the security and limitations.

-----

### CLI Usage:

`./go-dcc -h` for more help on usages and parameters.<br>
<br>

### Generation
#### Create and Sign Greenpass/EU-DCC with sample data that matches required schema for vaccine passports. <br>
`./go-dcc -gen -in "./data.json" -out "./test.png"` to generate from json data and export DCC as QR code in PNG.<br>
`./go-dcc -gen -raw -in "./data.json" -out "./test.txt"` to generate from json data and export DCC as txt.<br>
<br>

### Validation
#### Validate Greenpass/EU-DCC's. *todo:*<br>
`./go-dcc -verify -in "./eu-dcc.png"` to validate DCC from QR code.<br>
`./go-dcc -verify -raw -in "./eu-dcc.txt"` to validate raw DCC from txt file.<br>
<br>

### Inspection
#### Parse/Print contents of Greenpass/EU-DCC <br>
`./go-dcc -info -in "./eu-dcc.png"` to read/parse a QR Code DCC.<br>
`./go-dcc -info -raw -in "./eu-dcc.txt"` to read/parse a raw txt file DCC.<br>
<br>

-----

### Build instructions
`build.ps1 linux|windows` to build executable.
