package producer

import "math/rand"

type Role string

const (
	RoleCEO        Role = "CEO"
	RoleEngineer   Role = "Engineer"
	RoleContractor Role = "Contractor"
	RoleHR         Role = "HR"
	RoleFinance    Role = "Finance"
	RoleITAdmin    Role = "IT_Admin"
)

type User struct {
	Name string
	Role Role
}

type MarkovChain struct {
	States      []string
	Transitions map[string][]float64
	Initial     []float64
}

func (mc *MarkovChain) GetNextAction(currentState string) string {
	if currentState == "" || mc.Transitions[currentState] == nil {
		return mc.sampleFromDistribution(mc.Initial)
	}
	probs := mc.Transitions[currentState]
	return mc.sampleFromDistribution(probs)
}

func (mc *MarkovChain) sampleFromDistribution(probs []float64) string {
	r := rand.Float64()
	cumulative := 0.0
	for i, prob := range probs {
		cumulative += prob
		if r <= cumulative {
			return mc.States[i]
		}
	}
	return mc.States[0]
}

func GetRoleMarkovChain(role Role) *MarkovChain {
	switch role {
	case RoleCEO:
		return &MarkovChain{
			States: []string{"emailed_external", "shared_externally", "accessed", "downloaded", "uploaded_to_cloud"},
			Transitions: map[string][]float64{
				"emailed_external":  {0.3, 0.2, 0.2, 0.15, 0.15},
				"shared_externally": {0.25, 0.3, 0.2, 0.15, 0.1},
				"accessed":          {0.2, 0.2, 0.3, 0.2, 0.1},
				"downloaded":        {0.15, 0.15, 0.25, 0.3, 0.15},
				"uploaded_to_cloud": {0.2, 0.2, 0.2, 0.2, 0.2},
			},
			Initial: []float64{0.3, 0.25, 0.2, 0.15, 0.1},
		}

	case RoleEngineer:
		return &MarkovChain{
			States: []string{"downloaded", "accessed", "modified", "bulk_downloaded", "copied_to_clipboard"},
			Transitions: map[string][]float64{
				"downloaded":          {0.2, 0.3, 0.25, 0.15, 0.1},
				"accessed":            {0.25, 0.25, 0.3, 0.1, 0.1},
				"modified":            {0.2, 0.3, 0.3, 0.1, 0.1},
				"bulk_downloaded":     {0.3, 0.2, 0.2, 0.2, 0.1},
				"copied_to_clipboard": {0.25, 0.25, 0.25, 0.15, 0.1},
			},
			Initial: []float64{0.3, 0.3, 0.2, 0.15, 0.05},
		}

	case RoleContractor:
		return &MarkovChain{
			States: []string{"accessed", "downloaded", "emailed_external", "uploaded_to_cloud", "shared_externally"},
			Transitions: map[string][]float64{
				"accessed":          {0.4, 0.3, 0.15, 0.1, 0.05},
				"downloaded":        {0.35, 0.35, 0.15, 0.1, 0.05},
				"emailed_external":  {0.3, 0.3, 0.25, 0.1, 0.05},
				"uploaded_to_cloud": {0.3, 0.3, 0.2, 0.15, 0.05},
				"shared_externally": {0.3, 0.3, 0.2, 0.1, 0.1},
			},
			Initial: []float64{0.4, 0.3, 0.15, 0.1, 0.05},
		}

	case RoleHR:
		return &MarkovChain{
			States: []string{"accessed", "downloaded", "emailed_external", "printed", "shared_externally"},
			Transitions: map[string][]float64{
				"accessed":          {0.3, 0.25, 0.2, 0.15, 0.1},
				"downloaded":        {0.25, 0.3, 0.2, 0.15, 0.1},
				"emailed_external":  {0.2, 0.2, 0.3, 0.15, 0.15},
				"printed":           {0.25, 0.25, 0.2, 0.2, 0.1},
				"shared_externally": {0.2, 0.2, 0.25, 0.15, 0.2},
			},
			Initial: []float64{0.3, 0.25, 0.25, 0.1, 0.1},
		}

	case RoleFinance:
		return &MarkovChain{
			States: []string{"accessed", "downloaded", "emailed_external", "uploaded_to_cloud", "printed"},
			Transitions: map[string][]float64{
				"accessed":          {0.25, 0.3, 0.2, 0.15, 0.1},
				"downloaded":        {0.2, 0.3, 0.25, 0.15, 0.1},
				"emailed_external":  {0.2, 0.25, 0.3, 0.15, 0.1},
				"uploaded_to_cloud": {0.2, 0.2, 0.2, 0.25, 0.15},
				"printed":           {0.25, 0.25, 0.2, 0.15, 0.15},
			},
			Initial: []float64{0.3, 0.3, 0.2, 0.1, 0.1},
		}

	case RoleITAdmin:
		return &MarkovChain{
			States: []string{"accessed", "downloaded", "modified", "deleted", "copied_to_usb"},
			Transitions: map[string][]float64{
				"accessed":      {0.2, 0.25, 0.25, 0.15, 0.15},
				"downloaded":    {0.25, 0.2, 0.25, 0.15, 0.15},
				"modified":      {0.2, 0.25, 0.25, 0.15, 0.15},
				"deleted":       {0.25, 0.25, 0.2, 0.2, 0.1},
				"copied_to_usb": {0.25, 0.25, 0.2, 0.15, 0.15},
			},
			Initial: []float64{0.25, 0.25, 0.25, 0.15, 0.1},
		}

	default:
		// Default to basic access pattern
		return &MarkovChain{
			States: []string{"accessed", "downloaded"},
			Transitions: map[string][]float64{
				"accessed":   {0.5, 0.5},
				"downloaded": {0.5, 0.5},
			},
			Initial: []float64{0.5, 0.5},
		}
	}
}

func GetRoleResources(role Role) []string {
	allResources := []string{
		"customer_database.csv", "employee_ssn_list.xlsx", "payroll_q4_2024.xlsx",
		"credit_card_data.csv", "bank_account_details.pdf", "tax_records_2024.xlsx",
		"confidential_merger_docs.docx", "salary_survey_data.csv",
		"sales_pipeline_q1.xlsx", "client_contract_template.docx", "budget_forecast_2025.xlsx",
		"product_roadmap_internal.pptx", "vendor_pricing_list.csv", "performance_reviews.pdf",
		"marketing_strategy_2024.docx", "competitive_analysis.xlsx",
		"company_logo.png", "public_readme.txt", "meeting_notes_general.docx",
		"training_materials.pdf", "org_chart_public.pdf", "company_handbook.pdf",
		"office_photos.zip", "presentation_template.pptx",
		"https://drive.google.com/sensitive_data", "https://dropbox.com/client_files",
		"\\\\shared_drive\\hr\\confidential", "ftp://vendor.com/upload",
		"email_attachment_financial_report.pdf",
	}

	// Role-specific resource preferences
	roleResources := map[Role][]string{
		RoleCEO: {
			"confidential_merger_docs.docx", "budget_forecast_2025.xlsx",
			"sales_pipeline_q1.xlsx", "competitive_analysis.xlsx",
			"performance_reviews.pdf", "marketing_strategy_2024.docx",
		},
		RoleEngineer: {
			"product_roadmap_internal.pptx", "competitive_analysis.xlsx",
			"client_contract_template.docx", "presentation_template.pptx",
			"training_materials.pdf", "company_handbook.pdf",
		},
		RoleContractor: {
			"training_materials.pdf", "company_handbook.pdf",
			"public_readme.txt", "meeting_notes_general.docx",
			"presentation_template.pptx", "company_logo.png",
		},
		RoleHR: {
			"employee_ssn_list.xlsx", "payroll_q4_2024.xlsx",
			"salary_survey_data.csv", "performance_reviews.pdf",
			"org_chart_public.pdf", "\\\\shared_drive\\hr\\confidential",
		},
		RoleFinance: {
			"payroll_q4_2024.xlsx", "credit_card_data.csv",
			"bank_account_details.pdf", "tax_records_2024.xlsx",
			"budget_forecast_2025.xlsx", "vendor_pricing_list.csv",
		},
		RoleITAdmin: {
			"customer_database.csv", "employee_ssn_list.xlsx",
			"credit_card_data.csv", "bank_account_details.pdf",
			"\\\\shared_drive\\hr\\confidential", "https://drive.google.com/sensitive_data",
		},
	}

	// Return role-specific resources, but allow some access to all resources
	if resources, ok := roleResources[role]; ok {
		// Mix role-specific with general resources (70% role-specific, 30% general)
		return append(resources, allResources...)
	}

	return allResources
}

func GetMaliciousMarkovChain() *MarkovChain {
	return &MarkovChain{
		States: []string{"bulk_downloaded", "copied_to_usb", "emailed_external", "uploaded_to_cloud", "accessed"},
		Transitions: map[string][]float64{
			"bulk_downloaded":   {0.3, 0.25, 0.2, 0.15, 0.1},
			"copied_to_usb":     {0.25, 0.3, 0.2, 0.15, 0.1},
			"emailed_external":  {0.2, 0.2, 0.3, 0.2, 0.1},
			"uploaded_to_cloud": {0.2, 0.2, 0.25, 0.25, 0.1},
			"accessed":          {0.15, 0.15, 0.2, 0.2, 0.3},
		},
		Initial: []float64{0.3, 0.25, 0.2, 0.15, 0.1},
	}
}
