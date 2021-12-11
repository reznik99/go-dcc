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
    A simple library to Decode, Generate and Validate EU-DCC / Vaccine Passports 
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
This package was a project to help me learn better about the protocols used.

Here's why:
* To understand how Vaccine Passports / EU-DCC / EU-DGC / Greenpasses work.
* Analyze what personal information is stored within.
* Understand the security and limitations.




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

1. To generate an example EU-DCC vaccine certificate as QR code using data from `data.json`
   ```go
   import (
       "github.com/reznik99/go-dcc"
   )

   func main() {
       var err = dcc.GenerateQR("./data.json", "DCC-Test.png")
   }
   ```
2. To generate an example EU-DCC vaccine certificate as raw
   ```go
   import (
       "fmt"
       "github.com/reznik99/go-dcc"
   )

   func main() {
       data, err := dcc.Generate("./data.json")

       fmt.Printf("%s", data)

       // Will print `HC1:......`
   }
   ```
3. To decode/read an EU-DCC QR code Vaccine Pass
   ```go
    import (
        "fmt"
        "github.com/reznik99/go-dcc"
        "github.com/ethereum/go-ethereum/common/hexutil"
    )

    func main() {
        payload, rawMsg, _ := dcc.ParseQR("path/to/qr-code.png")
        prettyDCC, _ := json.MarshalIndent(payload, "", "	")

        fmt.Printf("Headers:   %v\n", rawMsg.Headers)
        fmt.Printf("Payload:   %s\n", prettyDCC)
        fmt.Printf("Signature: %s\n", hexutil.Encode(rawMsg.Signature))
    }
   ```



_For more examples, please refer to the [Documentation](https://pkg.go.dev/github.com/reznik99/go-dcc)_

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [x] Decode/Read EU DCC certs
- [x] Encode/Generate EU DCC certs (Valid schema but not valid signature obviously)
- [ ] Verify/Validate EU DCC certs (This is not working quite yet)

See the [open issues](https://github.com/reznik99/go-dcc/Best-README-Template/issues) for a full list of proposed features (and known issues).
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
