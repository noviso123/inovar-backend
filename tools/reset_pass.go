package main

import (
	"fmt"
	"inovar/lib/shared"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Carregar .env local (executando da raiz do projeto backend)
	if err := godotenv.Load(".env"); err != nil {
		// Tentar carregar do diretório pai caso esteja rodando de tools/
		godotenv.Load("../.env")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("DATABASE_URL não configurada")
	}

	if err := shared.InitDB(); err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	email := "admin@inovar.com"
	newPassword := "123456"

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Erro ao gerar hash: %v", err)
	}

	db := shared.GetDB()
	var user shared.User

	// Verificar se usuário existe
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		log.Printf("Usuário %s não encontrado. Criando...", email)
		user = shared.User{
			Name:         "Admin",
			Email:        email,
			Role:         "admin",
			PasswordHash: string(hashedPassword),
			Active:       true,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Erro ao criar usuário: %v", err)
		}
		fmt.Println("✅ Usuário criado com sucesso!")
	} else {
		// Atualizar senha
		user.PasswordHash = string(hashedPassword)
		if err := db.Save(&user).Error; err != nil {
			log.Fatalf("Erro ao atualizar senha: %v", err)
		}
		fmt.Println("✅ Senha atualizada com sucesso!")
	}
}
