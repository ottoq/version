# version

Golang version package generation tool used to standardize reporting compile time version information across golang applications.

###  ENV Variables

The following environment variables are expected to be set in your build pipeline:

* `CI_PIPELINE_ID` - set by Gitlab CI for current pipeline
* `CI_BUILD_ID` - set by Gitlab CI for current build

### Example output:

```
ID:          DEV
Description: 5db53ac5f6-dirty
Hostname:    foobar123.local
Go Runtime:  devel +5dd7108 Tue Dec 6 16:30:47 2016 -0800
```

## Installation

You can install this either as a `govendor'd` dependency, or as a submodule:

### govendor

```
govendor fetch github.com/ottoq/version
govendor fetch github.com/ottoq/version/gen
```

### Submodule

Add this repository to your project as a submodule. From the top level of your git repository execute these commands:
```
mkdir -p vendor/github.com/ottoq
git submodule add https://github.com/ottoq/version.git vendor/github.com/ottoq/version
```

## Project setup

Edit your `main` package file adding the following generate directive somewhere close to the top:
```
//go:generate go run vendor/github.com/ottoq/version/gen/gen.go
```

Import the `version` package anywhere you need to query your application's version:
```
import (
    "github.com/ottoq/version"
)
```

Utilize the `version` package - call `version.String()` to find the compile time application version.

__On each new release build__, usually handled for you via CI Run `go generate` - this generates your build version data. Not rerunning `generate` before release builds will report old version labels.
