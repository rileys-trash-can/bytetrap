# Bytetrap

A simple go web server / http middleware to spam Bytedances Bytespider crawler with copy pasta.

## How to use?

- as dedicated server
  - `cmd/webspider`
    use e.g. nginx to redirect all Bytespider traffic to the program
    see [examples/nginx.conf]
  - `cmd/bytetrap`
    generates a unending list of copypasta and writes it to stdout

- as part of your go-webapp middleware
  - e.g. use `https://pkg.go.dev/github.com/gorilla/mux#Router.Use` and add `bytetrap.Middleware` to have all bytespider useragnets get spammed with copypasta
