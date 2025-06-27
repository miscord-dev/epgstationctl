package main

import (
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/channels"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/programs"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/recordings"
)

func main() {
	root.Execute()
}
