package main

import (
	"fmt"
	"inovar/lib/shared"
	"log"
	"os"
)

func main() {
	// Force the production URL
	os.Setenv("DATABASE_URL", "postgresql://postgres:Inovar2026!Secure@db.bxbupbnjcingfvjszrau.supabase.co:5432/postgres")

	if err := shared.InitDB(); err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	db := shared.GetDB()
	var users []shared.User

	if err := db.Find(&users).Error; err != nil {
		log.Fatalf("Erro ao listar usuários: %v", err)
	}

	fmt.Println("Lista de Usuários:")
	fmt.Println("--------------------------------------------------")
	for _, u := range users {
		fmt.Printf("ID: %v | Nome: %s | Email: %s | Role: %s | Active: %v\n", u.ID, u.Name, u.Email, u.Role, u.Active)
	}
	fmt.Println("--------------------------------------------------")
}
