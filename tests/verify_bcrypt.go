package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "admin123"
	
	// Hash yang kita gunakan
	hashFromDB := "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.2U8xGTv5MBDj2iqpYi"
	
	fmt.Println("Testing Bcrypt Password Verification")
	fmt.Println("=====================================")
	fmt.Printf("Password: %s\n", password)
	fmt.Printf("Hash: %s\n", hashFromDB)
	fmt.Println()
	
	// Test 1: Verify password dengan hash dari DB
	err := bcrypt.CompareHashAndPassword([]byte(hashFromDB), []byte(password))
	if err != nil {
		fmt.Printf("❌ Verification FAILED: %v\n", err)
		fmt.Println("\nHash mungkin invalid atau password tidak match.")
	} else {
		fmt.Println("✅ Verification SUCCESS")
		fmt.Println("Hash dan password cocok!")
	}
	
	fmt.Println()
	fmt.Println("Testing Alternative Hashes:")
	fmt.Println("============================")
	
	// Alternative hashes to test
	alternativeHashes := []string{
		"$2y$10$9t0eT3bFLvFCwZP1.LFbCueNJ.uXsQQQb7vGlpPp5j9lB7Jl6zYwm", // $2y$ variant
		"$2b$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.2U8xGTv5MBDj2iqpYi", // $2b$ variant
	}
	
	for i, hash := range alternativeHashes {
		fmt.Printf("\nHash %d: %s\n", i+1, hash)
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err != nil {
			fmt.Printf("❌ FAILED: %v\n", err)
		} else {
			fmt.Println("✅ SUCCESS")
		}
	}
	
	fmt.Println()
	fmt.Println("Generating Fresh Hash:")
	fmt.Println("======================")
	
	// Generate fresh hash
	freshHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generating hash: %v\n", err)
	} else {
		fmt.Printf("Fresh Hash: %s\n", string(freshHash))
		
		// Verify fresh hash
		err = bcrypt.CompareHashAndPassword(freshHash, []byte(password))
		if err != nil {
			fmt.Printf("Fresh hash verification FAILED: %v\n", err)
		} else {
			fmt.Println("Fresh hash verification SUCCESS ✅")
		}
	}
}
