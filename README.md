# Echo SCS Session Middleware

[![GoDoc](https://godoc.org/github.com/canidam/echo-scs-session?status.png)](https://pkg.go.dev/github.com/canidam/echo-scs-session?tab=doc)
[![Go report card](https://goreportcard.com/badge/github.com/canidam/echo-scs-session)](https://goreportcard.com/report/github.com/canidam/echo-scs-session)

All credit must go to [alexedwards](https://github.com/alexedwards) for his great package [scs](https://github.com/alexedwards/scs)

_This is a fork of [spazzymoto/echo-scs-session](https://github.com/spazzymoto/echo-scs-session) package with bug fixes_

This package includes a fix to critical issue [Middleware does not set max-age/expires #1](https://github.com/spazzymoto/echo-scs-session/issues/1).
After spending an hour trying to figure this out, I found the bug in the middleware and decided to fix it under a new repo, 
update the dependency packages and the Go version

The upstream is not maintained anymore.

## Features (As per [SCS](https://github.com/alexedwards/scs))

* Automatic loading and saving of session data via middleware.
* Choice of server-side session stores including PostgreSQL, MySQL, Redis, BadgerDB and BoltDB. Custom session stores are also supported.
* Supports multiple sessions per request, 'flash' messages, session token regeneration, idle and absolute session timeouts, and 'remember me' functionality.
* Easy to extend and customize. Communicate session tokens to/from clients in HTTP headers or request/response bodies.
* Efficient design. Smaller, faster and uses less memory than [gorilla/sessions](https://github.com/gorilla/sessions).


## Instructions

* [Installation](#installation)
* [Basic Use](#basic-use)


### Installation

This package requires Go 1.22 or newer.

```
$ go get github.com/canidam/echo-scs-session
```

### Basic Use
See [scs](https://github.com/alexedwards/scs) for further configuration options for scs sessions

```go
package main

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
	"github.com/canidam/echo-scs-session"
)

var sessionManager *scs.SessionManager

func main() {
	// Initialize a new session manager and configure the session lifetime.
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	e := echo.New()

	// Use the LoadAndSave() middleware.
	e.Use(session.LoadAndSave(sessionManager))

	e.GET("/put", putHandler)
	e.GET("/get", getHandler)

	e.Logger.Fatal(e.Start(":4000"))
}

func putHandler(c echo.Context) error {
	// Store a new key and value in the session data.
	sessionManager.Put(c.Request().Context(), "message", "Hello from a session!")

	return c.String(http.StatusOK, "")
}

func getHandler(c echo.Context) error {
	// Use the GetString helper to retrieve the string value associated with a
	// key. The zero value is returned if the key does not exist.
	msg := sessionManager.GetString(c.Request().Context(), "message")

	return c.String(http.StatusOK, msg)
}
```

```
$ curl -i --cookie-jar cj --cookie cj localhost:4000/put
HTTP/1.1 200 OK
Cache-Control: no-cache="Set-Cookie"
Content-Type: text/plain; charset=UTF-8
Set-Cookie: session=0KumL8V5WYuvZwEQj2IPrYvm-cC3y7m8xQWLhTmxq_U; Path=/; HttpOnly; SameSite=Lax; Expires=Tue, 2 April 2024 09:28:00 GMT;
Vary: Cookie
Date: Mon, 1 April 2024 08:28:00 GMT
Content-Length: 0

$ curl -i --cookie-jar cj --cookie cj localhost:4000/get
HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Date: Thu, 20 May 2021 14:47:40 GMT
Content-Length: 21

Hello from a session!
```
