# goathlon

![LICENSE-shield][license-shield-url] ![WORKFLOW-status][workflow-status-url]

A prototype system for processing and analyzing biathlon competition events.

## 🔧 Installation & Setup

### Prerequisites

- **Go**: v1.20 or newer.
- **Make** (optional, for simplified commands)

### Usage

#### Basic setup

```bash
git clone https://github.com/artem-burashnikov/goathlon.git
cd goathlon
```

#### Running program

You need to complete a few steps in order to run the program:

1. `CONFIG_PATH` environment varaible must be set.
2. Input stream must be provided.

```bash
CONFIG_PATH="test/testdata/single/config.json" go run . < test/testdata/single/events
```

Take a look at [examples](/examples/README.md).

#### Running tests

```bash
make test
```

#### Generate coverage report

```bash
make coverage
```

## 📜 License

The project is licensed under an MIT License.

<!---->
[license-shield-url]: https://img.shields.io/github/license/artem-burashnikov/goathlon?style=for-the-badge&color=blue
[workflow-status-url]: https://img.shields.io/github/actions/workflow/status/artem-burashnikov/goathlon/.github%2Fworkflows%2Fci.yaml?style=for-the-badge&color=lightgreen
