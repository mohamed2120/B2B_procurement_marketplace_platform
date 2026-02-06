package main

import (
	"fmt"
	"log"

	"github.com/b2b-platform/identity-service/models"
	"github.com/b2b-platform/identity-service/repository"
	"github.com/b2b-platform/shared/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := database.GetDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)

	// Create default roles
	roles := []models.Role{
		{Name: "super_admin", Description: "Super Administrator", IsSystem: true},
		{Name: "admin", Description: "Platform Administrator", IsSystem: true},
		{Name: "procurement_manager", Description: "Procurement Manager (Buyer Company)", IsSystem: true},
		{Name: "requester", Description: "Requester (Buyer Company)", IsSystem: true},
		{Name: "buyer", Description: "Buyer", IsSystem: true},
		{Name: "approver", Description: "Approver", IsSystem: true},
		{Name: "supplier", Description: "Supplier", IsSystem: true},
		{Name: "catalog_admin", Description: "Catalog Administrator", IsSystem: true},
		{Name: "equipment_manager", Description: "Equipment Manager", IsSystem: true},
		{Name: "company_admin", Description: "Company Administrator", IsSystem: true},
	}

	roleMap := make(map[string]uuid.UUID)
	for _, role := range roles {
		existing, err := roleRepo.GetByName(role.Name)
		if err == nil {
			roleMap[role.Name] = existing.ID
			fmt.Printf("Role %s already exists\n", role.Name)
			continue
		}

		if err := roleRepo.Create(&role); err != nil {
			log.Printf("Failed to create role %s: %v", role.Name, err)
			continue
		}
		roleMap[role.Name] = role.ID
		fmt.Printf("Created role: %s\n", role.Name)
	}

	// Create demo tenant
	demoTenantID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	// Create demo users
	password := "demo123456"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	users := []struct {
		email     string
		firstName string
		lastName  string
		roleNames []string
	}{
		{"admin@demo.com", "Platform", "Admin", []string{"admin"}},
		{"buyer.requester@demo.com", "Requester", "User", []string{"requester"}},
		{"buyer.procurement@demo.com", "Procurement", "Manager", []string{"procurement_manager"}},
		{"supplier@demo.com", "Supplier", "User", []string{"supplier"}},
	}

	for _, u := range users {
		// Check if user exists
		_, err := userRepo.GetByEmail(demoTenantID, u.email)
		if err == nil {
			fmt.Printf("User %s already exists\n", u.email)
			continue
		}

		user := &models.User{
			TenantID:     demoTenantID,
			Email:        u.email,
			PasswordHash: string(hashedPassword),
			FirstName:    u.firstName,
			LastName:     u.lastName,
			IsActive:     true,
			IsVerified:   true,
		}

		if err := userRepo.Create(user); err != nil {
			log.Printf("Failed to create user %s: %v", u.email, err)
			continue
		}

		// Assign roles
		for _, roleName := range u.roleNames {
			if roleID, ok := roleMap[roleName]; ok {
				if err := userRepo.AssignRole(user.ID, roleID, demoTenantID); err != nil {
					log.Printf("Failed to assign role %s to user %s: %v", roleName, u.email, err)
				}
			}
		}

		fmt.Printf("Created user: %s (password: %s)\n", u.email, password)
	}

	fmt.Println("Seeding completed!")
}
