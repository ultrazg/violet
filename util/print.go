package util

import "fmt"

func P(h, p, v string) {
	fmt.Printf(`
       _       _      _   
      (_)     | |    | |  
__   ___  ___ | | ___| |_ 
\ \ / / |/ _ \| |/ _ \ __|
 \ V /| | (_) | |  __/ |_ 
  \_/ |_|\___/|_|\___|\__|
                                                        
v%s                           
server started on %s:%s
`, v, h, p)
}
