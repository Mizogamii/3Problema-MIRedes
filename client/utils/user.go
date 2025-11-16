package utils

import (
	"fmt"
	"strconv"
	
	"pbl/shared"
)

func Cadastro() shared.User {
	var user shared.User

	fmt.Print("Insira o nome do usuário: ")
	user.UserName = ReadLineSafe()
	fmt.Print("Insira a senha desejada: ")
	user.Password = ReadLineSafe()

	return user
}

func Login() shared.User {
	var user shared.User

	fmt.Print("Insira o nome do usuário: ")
	user.UserName = ReadLineSafe()
	print("Insira a sua senha: ")
	user.Password = ReadLineSafe()

	return user
}

func Troca() int {
	fmt.Println("\nDigite o número da carta que deseja trocar: ")
	for {
        input := ReadLineSafe()
        num, err := strconv.Atoi(input)
        if err != nil {
            fmt.Println("Digite um número válido.")
            continue
        }
        return num
    }
}
