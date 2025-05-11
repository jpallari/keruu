# Keruu

Aggregating RSS/Atom feeds to a single HTML page.

## Installation

First, you need to install [Go](https://golang.org/dl/) version 1.23 or higher.
After that, you can use `go get` to install Keruu:

    $ go get go.lepovirta.org/keruu

The executable should now be in path `$GOPATH/bin/keruu` (or `~/go/bin/keruu`).

Alternatively, if you have [Docker](https://docker.com/) installed, you can run the Docker image from the following location:

```
ghcr.io/jpallari/keruu:main
```

## Usage

Keruu accepts the following CLI flags:

* `-config`: Path to the configuration file (default: read from STDIN)
* `-output`: Path to the HTML output file (default: write to STDOUT)
* `-help`: Displays how to use Keruu

## Configuration

Keruu is configured using YAML. Here's all the configurations accepted by Keruu:

* `feeds`: A list of RSS/Atom feeds to aggregate. At least one feed must be provided.
  * `name` (optional): Name of the feed
  * `url`: URL for the feed
  * `exclude` (optional): A list of regular expression patterns to match against the feed post titles.
    If a post title matches any of the expressions, the post is excluded from the HTML output.
  * `include` (optional): A list of regular expression patterns to match against the feed post titles.
    Only posts that match the expressions are included in the HTML output unless they match the expressions in the `exclude` list.
* `fetch` (optional): A map containing the following configurations
  * `httpTimeout` (optional): Duration for how long to wait for a single feed fetch
  * `propagateErrors` (optional): When set to `true`, the program will return an error code when feed fetching or parsing fails for one or more feeds. This can be useful when you want to catch feed errors in scripts.
* `aggregation` (optional):
  * `title` (optional): Title to use in the HTML output
  * `description` (optional): Description to use in the HTML output
  * `maxPosts` (optional): Maximum number of posts to include in the HTML output
  * `grouping` (optional): Group posts in the HTML output. These options are available.
    * `monthly`: Group posts by month (default)
    * `weekly`: Group posts by week number of the year
    * number (e.g. `10`): Group posts into groups of given size
    * `none` or empty string: No grouping at all
  * `css` (optional): Custom CSS for the HTML output
* `links` (optional): A list of links to generate per feed item.
  * `name`: A name to display for the link
  * `url`: An URL pattern to use for the link.
    In the pattern, `$TITLE` will be replaced with the feed item title,
    and `$URL` will be replaced by the feed item link.

## Example

See the `example/config.yaml` directory for an example configuration.

See [lepo-keruu](https://github.com/jpallari/lepo-keruu) for an example of Keruu scheduled via GitHub Actions.

## GitHub Action

You can run Keruu via GitHub Actions like this:

```yaml
jobs:
  build:
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Keruu
        uses: jpallari/keruu@main
        with:
          config-path: config.yaml
          output-path: index.html
```

The action accepts two parameters:

* `config-path`: Path to Keruu config
* `output-path`: Path to Keruu output (HTML file)

## Docker

Image tag: `ghcr.io/jpallari/keruu:main`

The minimal image is optimized for the use in the command-line and scripting.
Besides the Keruu tool, it only includes the bare minimum system dependencies.
Running the container runs the Keruu tool directly and the container parameters are passed to the tool.

Example:

```
docker run -v $(pwd):/workspace:z \
  ghrc.io/jpallari/keruu:main \
  -config /workspace/config.yaml \
  -output /workspace/output.html
```

## License

GNU General Public License v3.0

See LICENSE file for more information.
