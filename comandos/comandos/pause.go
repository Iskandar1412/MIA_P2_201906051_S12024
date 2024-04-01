package comandos

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Pause() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Command Pause: > ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input != "" {
		fmt.Println("Ok, dont press any key exept enter -> Continue the process...")
	}
	fmt.Println("Continue the process...")
}
