# gohislip

Go library with HiSLIP implementation.

* [About](#user-content-about)
* [Installation](#user-content-installation)
* [Resource](#user-content-resource)
    * [Connect](#user-content-connect)
    * [Write](#user-content-write)
    * [Query](#user-content-query)
* [Examples](#user-content-examples)
* [Contributing](#user-content-contributing)
* [License](#user-content-license)
* [Authors](#user-content-authors)

### About

HiSLIP — is a TCP-based protocol for remote control of Test & Measurement instruments (such as oscilloscopes, power supplies, multimeters, spectrum analyzers etc.)

Protocol of HiSLIP, you can see 
[here](https://web.archive.org/web/20130127105446/http://www.ivifoundation.org/downloads/Class%20Specifications/IVI-6.1_HiSLIP-1.1-2011-02-24.pdf).

### Installation

Via go get command:

```bash
go get https://github.com/0z0sk0/gohislip
```

### Resource

First of all, you need to create a resource:

```go
resource := gohislip.NewHislipResource()
```

In source, it's just an embedded type of internal `resource.Resource`, which contains session data.

#### Connect

To connection you need to use `Connect` method, which supports the patterns of HiSLIP address:

```go
err := resource.Connect("TCPIP0::127.0.0.1::hislip0::INSTR") // default 4880 port
err := resource.Connect("TCPIP0::127.0.0.1::hislip0,4881::INSTR")
```

#### Write

```go
err := resource.Write("SYSTem:PRESet")
if err != nil {
  ...
}
```

#### Query

```go
response, err := resource.Query("SERVice:PORT:COUNt?")
if err != nil {
  ...
}
```

### Examples

See `/examples` folder.

### Contributing

gohislip is open source. If you want to help out with the project please feel free to join in.

All contributions (bug reports, code, doc, ideas, etc.) are welcome.

Please use the github issue tracker and pull request features.

### License 
gohislip includes code covered by the MIT license.

For license details please see the LICENSE file.

### Authors
Created by 0Z0SK0 <annnatoliy@icloud.com>

See the AUTHORS file for full list of contributors.
