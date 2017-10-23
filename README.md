# pack

Packs the assets for your web server.

## Performance

Pack uses parallel compilation via job queues and is therefore extremely fast, much faster than the popular [webpack](https://github.com/webpack/webpack).

78 [pixy](https://github.com/aerogo/pixy) templates, 64 [scarlet](https://github.com/aerogo/scarlet) styles and 30 scripts can be compiled in [less than 60 milliseconds](https://gist.github.com/blitzprog/878ec0dfbcb4e2d7759c4119e004b68c).

For comparison, webpack needs about 50 milliseconds for a single `Hello World` script.

## Installation

```
go get -u github.com/aerogo/pack
go install github.com/aerogo/pack
```

## Usage

Run `pack` in your project directory. It will generate a `components` package in your root directory which you can import in your project.

## Components

Since `components` is a generated directory you should have this directory in your `.gitignore` file.

### CSS

```go
components.CSS()
```

Returns the CSS bundle which is a string of CSS containing all styles.

### JS

```go
components.JS()
```

Returns the JS bundle which is a string of JS containing all scripts.

### Templates

Templates are registered as public functions in the `components` package and can be called directly. All components are global, thus you can call a component from one file in another file without any import directives. Components return an HTML `string` but they use a single `bytes.Buffer` via pooling and streaming under the hood, which is extremely fast.