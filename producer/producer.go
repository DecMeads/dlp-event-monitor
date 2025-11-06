package producer

import (
	"channel_filter/event"
	"math/rand"
	"time"
)

type Producer struct {
	producer_id   string
	message_count int
}

func NewProducer(producer_id string) *Producer {
	return &Producer{
		producer_id: producer_id,
	}
}

func (p *Producer) createRandomEvent() event.Event {
	users := []string{
		"hr.alice.johnson", "hr.bob.wilson", "finance.carol.davis",
		"finance.david.miller", "marketing.eve.brown", "engineering.frank.garcia",
		"sales.grace.martinez", "legal.henry.taylor", "it.admin.root",
		"contractor.mike.temp", "contractor.sarah.ext", "vendor.john.consultant",
		"temp.lisa.summer", "contractor.alex.qa",
		"exec.ceo.smith", "exec.cfo.jones", "exec.hr.director",
	}

	actions := []string{
		"downloaded", "uploaded_to_cloud", "copied_to_usb", "emailed_external",
		"bulk_downloaded", "accessed", "deleted", "modified", "shared_externally",
		"printed", "screenshot_taken", "copied_to_clipboard",
	}

	resources := []string{
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

	timeVariance := time.Duration(rand.Intn(3600)) * time.Second
	eventTime := time.Now().Add(-timeVariance)

	return event.Event{
		User:       users[rand.Intn(len(users))],
		Action:     actions[rand.Intn(len(actions))],
		Resource:   resources[rand.Intn(len(resources))],
		Timestamp:  eventTime,
		ProducerId: p.producer_id,
	}
}

func (p *Producer) Produce(ch chan<- event.Event) {
	for {
		evt := p.createRandomEvent()
		p.message_count++
		ch <- evt
		time.Sleep(time.Duration(rand.Intn(2000)+500) * time.Millisecond)
	}
}

func (p *Producer) GetMessageCount() int {
	return p.message_count
}
