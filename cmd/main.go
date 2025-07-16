package main

import caddycmd "github.com/caddyserver/caddy/v2/cmd"
import _ "github.com/zxspirit/ddns"

// plug in Caddy modules here
import _ "github.com/caddyserver/caddy/v2/modules/standard"

func main() {
	caddycmd.Main()
}
