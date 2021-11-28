package services

import (
	"bufio"
	"fmt"
	"os"
)

func PrintLogo() {
	fmt.Println("██████╗ ██╗     ██╗████████╗███████╗███████╗██╗  ██╗ █████╗ ██████╗ ███████╗")
	fmt.Println("██╔══██╗██║     ██║╚══██╔══╝╚══███╔╝██╔════╝██║  ██║██╔══██╗██╔══██╗██╔════╝")
	fmt.Println("██████╔╝██║     ██║   ██║     ███╔╝ ███████╗███████║███████║██████╔╝█████╗  ")
	fmt.Println("██╔══██╗██║     ██║   ██║    ███╔╝  ╚════██║██╔══██║██╔══██║██╔══██╗██╔══╝  ")
	fmt.Println("██████╔╝███████╗██║   ██║   ███████╗███████║██║  ██║██║  ██║██║  ██║███████╗")
	fmt.Println("╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝")
}

func ReadStdInLine() *string {
	stdReader := bufio.NewReader(os.Stdin)
	line, err := stdReader.ReadString('\n')
	if err != nil {
		return nil
	}
	return &line
}
