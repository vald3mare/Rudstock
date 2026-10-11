package services

/*func Register(username, phone, password, email string) (authToken string, refreshToken string, err error) {
	formatEmail := strings.ToLower(strings.TrimSpace(email))
	checkEmail := strings.Contains(formatEmail, "@")
	if !checkEmail {
		return "", "", fmt.Errorf("invalid email format")
	}

	if password == "" {
		return "", "", fmt.Errorf("password cannot be empty")
	}

	// хеширование пароля с использованием argon2id
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	userRole := h.svc.CreateUser(username, phone, hash, formatEmail)

	authToken, err := h.svc.GenerateAuthToken(userRole)
	if err != nil {
		return fmt.Errorf("failed to generate auth token: %v", err)
	}
	refreshToken, err := h.svc.GenerateRefreshToken(userRole)
	if err != nil {
		return fmt.Errorf("failed to generate refresh token: %v", err)
	}

	if email != "" && checkEmail {
		err = h.svc.SendVerificationEmail(formatEmail, authToken)

		log.Printf("Verification email sent to %s: %v", formatEmail, err)
		if err != nil {
			return fmt.Errorf("failed to send verification email: %v", err)
		}
	}

	// Return the generated tokens or any other relevant information
	// For example, you might want to return the authToken and refreshToken
	fmt.Printf("Auth Token: %s\n", authToken)
	fmt.Printf("Refresh Token: %s\n", refreshToken)

	return authToken, refreshToken, nil
}