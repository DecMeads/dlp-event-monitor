package producer

// GenerateCorporateUsers creates users for each role according to the hierarchy
func GenerateCorporateUsers() []User {
	users := []User{}

	// 1 CEO
	users = append(users, User{Name: "Sarah Chen", Role: RoleCEO})

	// 5 Engineers
	engineers := []string{
		"Frank Garcia", "Priya Patel", "Marcus Johnson",
		"Emily Rodriguez", "David Kim",
	}
	for _, name := range engineers {
		users = append(users, User{Name: name, Role: RoleEngineer})
	}

	// 2 Contractors
	contractors := []string{
		"Mike Temp", "Alex QA",
	}
	for _, name := range contractors {
		users = append(users, User{Name: name, Role: RoleContractor})
	}

	// 3 HR
	hr := []string{
		"Alice Johnson", "Bob Wilson", "Lisa Martinez",
	}
	for _, name := range hr {
		users = append(users, User{Name: name, Role: RoleHR})
	}

	// 3 Finance
	finance := []string{
		"Carol Davis", "David Miller", "Jennifer White",
	}
	for _, name := range finance {
		users = append(users, User{Name: name, Role: RoleFinance})
	}

	// 1 IT Admin
	users = append(users, User{Name: "Root Admin", Role: RoleITAdmin})

	return users
}
