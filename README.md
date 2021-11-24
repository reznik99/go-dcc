# go-dcc
EU Digital Covid Certificate utilities in Go [Create, Validate and Parse Green-Pass/EU-DCC] 
<br>
<br>
_Repo work in-progress_
<br>
<br>
#### CLI Usage:
######Create and Sign Greenpass/EU-DCC with sample data that matches required schema for vaccine passports. <br>
`./go-dcc -gen -in data.json -out eu-dcc.png`


######Validate Greenpass/EU-DCC's <br>
`./go-dcc -verify -in eu-dcc.png`
`./go-dcc -verifyRaw -in eu-dcc.txt`


######Parse/Print contents of Greenpass/EU-DCC <br>
`./go-dcc -info -in eu-dcc.png`
`./go-dcc -infoRaw -in eu-dcc.txt`
