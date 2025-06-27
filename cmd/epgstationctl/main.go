package main

import (
	_ "github.com/miscord-dev/epgstationctl/internal/commands/channels"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/programs"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/recordings"
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/rules"
)

func main() {
	root.Execute()
}
