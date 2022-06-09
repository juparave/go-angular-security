## go-server

Create go module

    $ mkdir go-app && cd go-app
    $ go mod init go-app

Install fiber v2

    $ go get github.com/gofiber/fiber/v2

## Installing go-app dependencies

    $ go mod download
    $ go mod tidy

## Dev Cycle

Codespace URL: https://juparave-go-angular-security-pgvrxv62rvj5-4200.githubpreview.dev/

### Autorefresh server with (Air)[https://github.com/cosmtrek/air]

Install air executable in ~/bin/ and inside `go-app` run:

    $ air
