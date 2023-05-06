package vfacore

import "fmt"

func getBanner() string {

	bannerF := `
 222222222222222        ffffffffffffffff                    
2:::::::::::::::22     f::::::::::::::::f                   
2::::::222222:::::2   f::::::::::::::::::f                  
2222222     2:::::2   f::::::fffffff:::::f                  
            2:::::2   f:::::f       ffffff  aaaaaaaaaaaaa   
            2:::::2   f:::::f               a::::::::::::a  
         2222::::2   f:::::::ffffff         aaaaaaaaa:::::a 
    22222::::::22    f::::::::::::f                  a::::a 
  22::::::::222      f::::::::::::f           aaaaaaa:::::a 
 2:::::22222         f:::::::ffffff         aa::::::::::::a 
2:::::2               f:::::f              a::::aaaa::::::a 
2:::::2               f:::::f             a::::a    a:::::a 
2:::::2       222222 f:::::::f            a::::a    a:::::a 
2::::::2222222:::::2 f:::::::f            a:::::aaaa::::::a 
2::::::::::::::::::2 f:::::::f             a::::::::::aa:::a
22222222222222222222 fffffffff              aaaaaaaaaa  aaaa

Author Vasilev Kirill https://vasilevkirill.ru

Telegram Hook running: %s 
Radius server running: %s
`

	nBanner := fmt.Sprintf(bannerF, configGlobalS.Telegram.WebHookAddress, configGlobalS.Radius.ServerAddress)
	return nBanner
}

func ShowBanner() {
	fmt.Print(getBanner())
}
