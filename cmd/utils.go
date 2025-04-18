package cmd

import "github.com/pterm/pterm"

func printBanner() {
	ivrit := `
   ___                   _____          _                  
  / _ \ _ __   ___ _ __ |  ___|__  __ _| |_ _   _ _ __ ___ 
 | | | | '_ \ / _ \ '_ \| |_ / _ \/ _` + "`" + ` | __| | | | '__/ _ \
 | |_| | |_) |  __/ | | |  _|  __/ (_| | |_| |_| | | |  __/
  \___/| .__/ \___|_| |_|_|  \___|\__,_|\__|\__,_|_|  \___|
       |_|                                                 
                                                    CLI    
`

	pterm.Println(ivrit)
	pterm.Println()
	pterm.Printf("version: %s | compiled: %s\n", pterm.LightGreen(Version), pterm.LightGreen(Date))
	pterm.Println(pterm.Cyan("🔗 https://openfeature.dev | https://github.com/open-feature/cli"))
}
