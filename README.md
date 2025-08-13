# svgbb

`svgbb` is a simple, fast, and portable command-line tool for generating SVG-based security payloads. It's designed for bug bounty hunters and security researchers who need to quickly create a set of SVG files for testing vulnerabilities like Server-Side Request Forgery (SSRF), and XML External Entity (XXE) injection.

The tool works by taking a set of SVG templates, injecting a user-provided string (like a URL or a JavaScript payload), and saving the results as new `.svg` files. All templates are embedded directly into the binary, making the tool a single, self-contained executable that you can run anywhere.

## Features

* **Portable**: Single binary with no external dependencies. Templates are embedded.
* **Simple**: Focused on one task: generating SVG payloads from templates.
* **Fast**: Written in Go for high performance.
* **Customizable**: Easily add your own SVG templates to expand your testing capabilities.

## Installation

To install `svgbb`, you can either build it from source or use `go install`:

```
go install github.com/vodafon/svgbb@latest
```

Make sure your Go `bin` directory (e.g., `$HOME/go/bin`) is in your system's `PATH`.

## Usage

The tool has three command-line flags:

| **Flag** | **Description** | **Default** | **Required** |
| :------- | :----------------------------------------------- | :---------- | :----------- |
| `-s`     | The string/payload to inject into the templates. | `""`        | Yes          |
| `-dir`   | The output directory for the generated SVGs.     | `.`         | No           |
| `-prefix`| An optional prefix for the output filenames.     | `""`        | No           |

### Example

Let's generate a set of SVG files for an XSS payload that will try to load a script from a collaborator server.

```
svgbb -s "https://your-collaborator-id.oastify.com" -dir "output" -prefix "ssrf-"
```
