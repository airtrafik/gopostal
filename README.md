# gopostal

Go/cgo interface to [libpostal](https://github.com/openvenues/libpostal), a C library for fast international street address parsing and normalization.

Fork of [openvenues/gopostal](https://github.com/openvenues/gopostal) with the following improvements:

- Explicit `Setup()` / `SetupDataDir()` / `Teardown()` lifecycle (no `init()` side effects)
- Configurable data directory via function parameter or `LIBPOSTAL_DATA_DIR` environment variable
- Proper error returns instead of `log.Fatal`
- Go module support

## Usage

### Expand addresses

```go
package main

import (
    "fmt"
    "log"

    expand "github.com/airtrafik/gopostal/expand"
)

func main() {
    if err := expand.Setup(); err != nil {
        log.Fatal(err)
    }
    defer expand.Teardown()

    expansions := expand.ExpandAddress("Quatre-vingt-douze Ave des Champs-Élysées")
    for _, e := range expansions {
        fmt.Println(e)
    }
}
```

### Parse addresses

```go
package main

import (
    "fmt"
    "log"

    parser "github.com/airtrafik/gopostal/parser"
)

func main() {
    if err := parser.Setup(); err != nil {
        log.Fatal(err)
    }
    defer parser.Teardown()

    parsed := parser.ParseAddress("781 Franklin Ave Crown Heights Brooklyn NY 11216 USA")
    for _, c := range parsed {
        fmt.Printf("%s: %s\n", c.Label, c.Value)
    }
}
```

### Custom data directory

Point libpostal at a non-default data directory (e.g., a GCS FUSE mount):

```go
if err := parser.SetupDataDir("/mnt/gcs/libpostal-data"); err != nil {
    log.Fatal(err)
}
```

Or set the `LIBPOSTAL_DATA_DIR` environment variable and call `Setup()`:

```bash
export LIBPOSTAL_DATA_DIR=/mnt/gcs/libpostal-data
```

## Prerequisites

Install the libpostal C library before using gopostal.

**On Mac**
```bash
brew install curl autoconf automake libtool pkg-config
```

**On Ubuntu/Debian**
```bash
sudo apt-get install curl autoconf automake libtool pkg-config
```

**Installing libpostal**

```bash
git clone https://github.com/openvenues/libpostal
cd libpostal
./bootstrap.sh
./configure --datadir=[...some dir with a few GB of space...]
make -j4
sudo make install
sudo ldconfig  # Linux only
```

## Installation

```bash
go get github.com/airtrafik/gopostal/expand
go get github.com/airtrafik/gopostal/parser
```

## Development

```bash
make          # Full build: clean, lint, test
make test     # Run tests only
make lint     # Run linters only
make lint-fix # Auto-fix formatting
make help     # Show all targets
```

## Releasing

Releases use [svu](https://github.com/caarlos0/svu) for semantic versioning based on conventional commits:

```bash
make release-preview  # Show what the next version would be
make release          # Interactive release with confirmation
```
