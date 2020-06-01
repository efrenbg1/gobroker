package linux
import (
	"os"
	"syscall"
	"fmt"
)

func TCP() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error tcp limit! ", err)
	}
	rLimit.Max = 20000
	rLimit.Cur = 20000
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error setting tcp limit!", err)
		os.Exit(3)
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error setting tcp limit!", err)
		os.Exit(3)
	}
	fmt.Println("tcp limit set to:", rLimit)
}
