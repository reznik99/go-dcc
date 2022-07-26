<div id="top"></div>


[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]


<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/reznik99/go-dcc">
    <img src="https://cdn.icon-icons.com/icons2/2699/PNG/512/golang_logo_icon_171073.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">GO-DCC</h3>

  <p align="center">
    A simple library to View, Generate and Verify EU-DCC / Vaccine Passes 
    <br />
    <a href="https://pkg.go.dev/github.com/reznik99/go-dcc"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/reznik99/go-dcc/issues">Report Bug</a>
    ·
    <a href="https://github.com/reznik99/go-dcc/issues">Request Feature</a>
  </p>
</div>


## About This Package

[![Product Name Screen Shot][screenshot]](https://example.com)

This package offers simple API calls to Decode, Encode and Verify Vaccine Passes.
<br>
APIs are subject to changes as improvements are made.

<!-- GETTING STARTED -->
## Getting Started
### Prerequisites

_Go get the module to import it into your project_

1. Download package
    ```sh
    go get -u github.com/reznik99/go-dcc
    ```

### Usage

_Below is some examples of how you can use this package, without explicit error handling for simplicity's sake_

1. To generate a Vaccine Pass using data from `data.json`
   ```go
   import (
       "github.com/reznik99/go-dcc"
   )

   func main() {
        // Base64 of first 8 bytes in Signer Certificate
        kid := "dy8HnMQYOrE="

        // Generate or load Signer Key
        privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

        // Generate/Sign Vaccine Pass and save as QR code
        dcc.GenerateQR(privKey, kid, "/path/to/data.json", "/path/to/new-pass.png")

        // Generate/Sign Vaccine Pass Raw string "HC1:..."
        dcc.Generate(privKey, kid, "/path/to/data.json")
   }
   ```
2. To validate/verify a Vaccine Pass
   ```go
   import (
       "github.com/reznik99/go-dcc"
   )

   func main() {
        // Parse Raw Vaccine Pass
        _, rawMsg, _ := dcc.Parse("/path/to/mypass.txt")
        // Parse QR Code Vaccine Pass
        _, rawMsg, _ := dcc.ParseQR("/path/to/mypass.png")

        // Verify Vaccine Pass signature. note: currently slow! This will fetch the PEM Signer Certificates and KIDs
        valid, _ := dcc.Verify(rawMsg)
        
        fmt.Printf("Vaccine Pass Valid: %t\n", valid)
   }
   ```

3. To decode/read a Vaccine Pass
   ```go
    import (
        "fmt"
        "github.com/reznik99/go-dcc"
        "github.com/ethereum/go-ethereum/common/hexutil"
    )

    func main() {
        // Parse Raw Vaccine Pass
        payload, _, _ := dcc.Parse("/path/to/mypass.txt")
        // Parse QR Code Vaccine Pass
        payload, _, _ := dcc.ParseQR("/path/to/mypass.png")

        // Print contents to STDOUT
        prettyDCC, _ := json.MarshalIndent(payload, "", "	")
        fmt.Printf("Headers:   %v\n", rawMsg.Headers)
        fmt.Printf("Payload:   %s\n", prettyDCC)
        fmt.Printf("Signature: %s\n", hexutil.Encode(rawMsg.Signature))
    }
   ```

_Example JSON data file_
```js
    {
        "name": "JOHN",
        "surname": "DOE",
        "dob": "1996-06-06",
        "issuerCountry": "NZ",
        "issuer": "Ministry of Health, NZ",
        "vaccinationDate": "2021-10-21",
        "doses": 2
    }
```


_For more examples, please refer to the [Documentation](https://pkg.go.dev/github.com/reznik99/go-dcc)_

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] Decode/Read EU DCC certs
- [x] Encode/Generate EU DCC certs (Valid schema but not valid signature obviously)
- [x] Verify/Validate EU DCC certs (This is not working quite yet)
- [ ] Improve KID and Signer Certificate fetching logic or allow user to specify values for performance.

See the [open issues](https://github.com/reznik99/go-dcc/issues) for a full list of proposed features (and known issues).
<br>
<br>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

Francesco Gorini - goras.francesco@gmail.com - https://francescogorini.com

Project Link: [https://github.com/reznik99/go-dcc/](https://github.com/reznik99/go-dcc/)

<p align="right">(<a href="#top">back to top</a>)</p>




<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/reznik99/go-dcc.svg?style=for-the-badge
[contributors-url]: https://github.com/reznik99/go-dcc/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/reznik99/go-dcc.svg?style=for-the-badge
[forks-url]: https://github.com/reznik99/go-dcc/network/members
[stars-shield]: https://img.shields.io/github/stars/reznik99/go-dcc.svg?style=for-the-badge
[stars-url]: https://github.com/reznik99/go-dcc/stargazers
[issues-shield]: https://img.shields.io/github/issues/reznik99/go-dcc?style=for-the-badge
[issues-url]: https://github.com/reznik99/go-dcc/issues
[license-shield]: https://img.shields.io/github/license/reznik99/go-dcc?style=for-the-badge
[license-url]: https://github.com/reznik99/go-dcc/blob/master/LICENSE
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://www.linkedin.com/in/francesco-gorini-b334861a6/
[screenshot]: res/read-me-banner.jpg
