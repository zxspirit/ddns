package main

import (
	"github.com/caddyserver/caddy/v2/cmd"
)
import _ "github.com/zxspirit/ddns"
import _ "github.com/caddyserver/caddy/v2/modules/standard"

func main() {
	caddycmd.Main()
}
