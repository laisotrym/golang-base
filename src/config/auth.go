package config

func AccessibleRoles() map[string][]string {
	// const bankServicePath = "/app.safeweb.v1.AgentAccessService/"
	const exposeServicePath = "/app.safeweb.v1.UserService/"

	return map[string][]string{
		exposeServicePath + "Get":    {"admin"},
		exposeServicePath + "Update": {"admin"},
		exposeServicePath + "Delete": {"admin"},
		exposeServicePath + "List":   {"admin"},

		// bankServicePath + "ListAll":   {"admin"},
		// bankServicePath + "Find":      {"admin", "user"},
		// bankServicePath + "Get":       {"admin", "user"},
	}
}
