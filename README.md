# pack

[![Godoc][godoc-image]][godoc-url]
[![Report][report-image]][report-url]
[![Tests][tests-image]][tests-url]
[![Coverage][coverage-image]][coverage-url]
[![Patreon][patreon-image]][patreon-url]

Packs the assets for your web server.

## Installation

```shell
go get -u github.com/blitzprog/home/...
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

## Coding style

Please take a look at the [style guidelines](https://github.com/akyoto/quality/blob/master/STYLE.md) if you'd like to make a pull request.

## Patrons

| [![Scott Rayapoullé](https://avatars3.githubusercontent.com/u/11772084?s=70&v=4)](https://github.com/soulcramer) |
|---|
| [Scott Rayapoullé](https://github.com/soulcramer) |

Want to see [your own name here](https://www.patreon.com/eduardurbach)?

## Author

| [![Eduard Urbach on Twitter](https://gravatar.com/avatar/16ed4d41a5f244d1b10de1b791657989?s=70)](https://twitter.com/eduardurbach "Follow @eduardurbach on Twitter") |
|---|
| [Eduard Urbach](https://eduardurbach.com) |

[godoc-image]: https://godoc.org/github.com/blitzprog/home?status.svg
[godoc-url]: https://godoc.org/github.com/blitzprog/home
[report-image]: https://goreportcard.com/badge/github.com/blitzprog/home
[report-url]: https://goreportcard.com/report/github.com/blitzprog/home
[tests-image]: https://cloud.drone.io/api/badges/blitzprog/home/status.svg
[tests-url]: https://cloud.drone.io/blitzprog/home
[coverage-image]: https://codecov.io/gh/blitzprog/home/graph/badge.svg
[coverage-url]: https://codecov.io/gh/blitzprog/home
[patreon-image]: https://img.shields.io/badge/patreon-donate-green.svg
[patreon-url]: https://www.patreon.com/eduardurbach
