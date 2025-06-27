package main

import (
	_ "github.com/miscord-dev/epgstationctl/internal/commands/channels"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/programs"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/recordings"
	_ "github.com/miscord-dev/epgstationctl/internal/commands/rules"
	"github.com/miscord-dev/epgstationctl/internal/commands/root"
)

func main() {
	root.Execute()
}
