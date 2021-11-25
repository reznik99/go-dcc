# go-dcc
EU Digital Covid Certificate utilities in Go [Create, Validate and Parse Green-Pass/EU-DCC] 

*Repo is a work in-progress*

-----

### Purpose:

- To further understand how Vaccine Passports / EU-DCC / EU-DGC / Greenpasses work.
- Analyze what personal information is store within.
- Understand the security and limitations.

-----

### CLI Usage:

`./go-dcc -h` for more help on usages and parameters

#### Create and Sign Greenpass/EU-DCC with sample data that matches required schema for vaccine passports. <br>
`./go-dcc -gen -in "./data.json" -out "eu-dcc.png"`


#### Validate Greenpass/EU-DCC's <br>
`./go-dcc -verify -in "./eu-dcc.png"`
`./go-dcc -verify -raw -in "./eu-dcc.txt"`


#### Parse/Print contents of Greenpass/EU-DCC <br>
`./go-dcc -info -in "./eu-dcc.png"`
`./go-dcc -info -raw -in "./eu-dcc.txt"`

-----

### Build instructions
`build.ps1 linux|windows` to build executable.
