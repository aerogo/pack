# pack

[![Go report card][goreportcard-image]][goreportcard-url]
[![Travis build][travis-image]][travis-url]
[![Code coverage][codecov-image]][codecov-url]
[![License][license-image]][license-url]

Packs the assets for your web server.

## Installation

```
go install github.com/aerogo/pack
```

## Usage

![pack usage](docs/usage.gif)

Run `pack` in your project directory. It will scan your project directory recursively to compile `.pixy`, `.scarlet` and `.js` files resulting in a `components` package in your root directory. You can then import the `components` package in your project to access all of your assets.

### Performance

Pack uses parallel compilation via job queues and is therefore extremely fast, much faster than the popular [webpack](https://github.com/webpack/webpack).

78 [pixy](https://github.com/aerogo/pack) templates, 64 [scarlet](https://github.com/aerogo/scarlet) styles and 30 scripts can be compiled in [less than 60 milliseconds](https://gist.github.com/akyoto/878ec0dfbcb4e2d7759c4119e004b68c). For comparison, webpack needs about 50 milliseconds for a single `Hello World` script.

## Components

Since `components` is a generated directory you should list this directory in your `.gitignore` file.

### CSS

```go
import "github.com/.../.../components/css"
```

```go
css.Bundle()
```

Returns the CSS bundle which is a string of CSS containing all styles.

### JS

```go
import "github.com/.../.../components/js"
```

```go
js.Bundle()
```

Returns the JS bundle which is a string of JS containing all scripts.

### Templates

```go
import "github.com/.../.../components"
```

Templates are registered as public functions in the `components` package and can be called directly. All components are global, thus you can call a component from one file in another file without any import directives. Components return an HTML `string` but they use a single `bytes.Buffer` via pooling and streaming under the hood, which is extremely fast.

[godoc-image]: https://godoc.org/github.com/aerogo/pack?status.svg
[godoc-url]: https://godoc.org/github.com/aerogo/pack
[goreportcard-image]: https://goreportcard.com/badge/github.com/aerogo/pack
[goreportcard-url]: https://goreportcard.com/report/github.com/aerogo/pack
[travis-image]: https://travis-ci.org/aerogo/pack.svg?branch=master
[travis-url]: https://travis-ci.org/aerogo/pack
[codecov-image]: https://codecov.io/gh/aerogo/pack/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/aerogo/pack
[license-image]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: https://github.com/aerogo/pack/blob/master/LICENSE
