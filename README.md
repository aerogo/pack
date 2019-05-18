# pack

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Patreon][patreon-image]][patreon-url]

Packs the assets for your web server.

## Installation

```shell
go get -u github.com/aerogo/pack/...
```

## Usage

![pack usage](docs/usage.gif)

Run `pack` in your project directory. It will scan your project directory recursively to compile `.pixy`, `.scarlet` and `.js` files resulting in a `components` package in your root directory. You can then import the `components` package in your project to access all of your assets.

### Performance

Pack uses parallel compilation via job queues and is therefore extremely fast, much faster than the popular [webpack](https://github.com/webpack/webpack).

## Components

Since `components` is a generated directory you should list this directory in your `.gitignore` file.

### CSS

```go
import "github.com/YOUR_ORG/YOUR_REPO/components/css"
```

```go
css.Bundle()
```

Returns the CSS bundle which is a string of CSS containing all styles.

### JS

```go
import "github.com/YOUR_ORG/YOUR_REPO/components/js"
```

```go
js.Bundle()
```

Returns the JS bundle which is a string of JS containing all scripts.

### Templates

```go
import "github.com/YOUR_ORG/YOUR_REPO/components"
```

Templates are registered as public functions in the `components` package and can be called directly. All components are global, thus you can call a component from one file in another file without any import directives. Components return an HTML `string` but they use a single `bytes.Buffer` via pooling and streaming under the hood, which is extremely fast.

## Style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.

## Sponsors

| [![Scott Rayapoullé](https://avatars3.githubusercontent.com/u/11772084?s=70&v=4)](https://github.com/soulcramer) | [![Eduard Urbach](https://avatars2.githubusercontent.com/u/438936?s=70&v=4)](https://twitter.com/eduardurbach) |
| --- | --- |
| [Scott Rayapoullé](https://github.com/soulcramer) | [Eduard Urbach](https://eduardurbach.com) |

Want to see [your own name here?](https://www.patreon.com/eduardurbach)

[godoc-image]: https://godoc.org/github.com/aerogo/pack?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/pack
[report-image]: https://goreportcard.com/badge/github.com/aerogo/pack
[report-url]: https://goreportcard.com/report/github.com/aerogo/pack
[tests-image]: https://cloud.drone.io/api/badges/aerogo/pack/status.svg
[tests-url]: https://cloud.drone.io/aerogo/pack
[coverage-image]: https://codecov.io/gh/aerogo/pack/graph/badge.svg
[coverage-url]: https://codecov.io/gh/aerogo/pack
[patreon-image]: https://img.shields.io/badge/patreon-donate-green.svg
[patreon-url]: https://www.patreon.com/eduardurbach
