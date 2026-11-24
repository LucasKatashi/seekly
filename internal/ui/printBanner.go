package ui

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func PrintBanner(silent bool) {
	if !silent {
		art := `
   ____        __    __    
  / __/__ ___ / /__ / /_ __
 _\ \/ -_) -_)  '_// / // /
/___/\__/\__/_/\_\/_/\_, / 
                    /___/`

		green := color.New(color.FgGreen)
		fmt.Fprintf(os.Stderr, "%s\n	   by: LucasKatashi\n\n", green.Sprint(art))
	}
}
