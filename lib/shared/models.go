package shared

import (
	"time"
)

// User model
type User struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Name               string    `json:"name"`
	Email              string    `json:"email" gorm:"unique;not null"`
	PasswordHash       string    `json:"-" gorm:"column:password_hash"`
	Role               string    `json:"role"` // admin, manager, technician
	CompanyID          *uint     `json:"companyId" gorm:"column:company_id"`
	Active             bool      `json:"active" gorm:"default:true"`
	MustChangePassword bool      `json:"mustChangePassword" gorm:"column:must_change_password;default:false"`
	AvatarURL          *string   `json:"avatarUrl" gorm:"column:avatar_url"`
	Phone              *string   `json:"phone"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// Client model
type Client struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	CNPJ      string    `json:"cnpj" gorm:"unique"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CompanyID uint      `json:"companyId" gorm:"column:company_id"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Request model (Service Request)
type Request struct {
	ID                uint       `json:"id" gorm:"primaryKey"`
	Title             string     `json:"title" gorm:"not null"`
	Description       string     `json:"description"`
	Status            string     `json:"status" gorm:"default:pending"` // pending, in_progress, completed
	Priority          string     `json:"priority"`                      // low, medium, high
	ClientID          uint       `json:"clientId" gorm:"column:client_id"`
	CompanyID         uint       `json:"companyId" gorm:"column:company_id"`
	CreatedBy         uint       `json:"createdBy" gorm:"column:created_by"`
	AssignedTo        *uint      `json:"assignedTo" gorm:"column:assigned_to"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	Observation       string     `json:"observation"`
	MaterialsUsed     string     `json:"materialsUsed"`
	NextMaintenanceAt *time.Time `json:"nextMaintenanceAt"`
	ScheduledAt       *time.Time `json:"scheduledAt"`
	PreventiveDone    bool       `json:"preventiveDone"`
}

type Equipment struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	SerialNumber string    `json:"serialNumber"`
	ClientID     uint      `json:"clientId" gorm:"column:client_id"`
	Active       bool      `json:"active" gorm:"default:true"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ChecklistItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RequestID   uint      `json:"requestId" gorm:"column:request_id"`
	Description string    `json:"description"`
	Checked     bool      `json:"checked" gorm:"default:false"`
	Observation string    `json:"observation"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type AgendaEntry struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Title        string    `json:"title"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	TechnicianID uint      `json:"technicianId" gorm:"column:technician_id"`
	RequestID    *uint     `json:"requestId" gorm:"column:request_id"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Company struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	LogoURL   string    `json:"logoUrl" gorm:"column:logo_url"`
	Website   string    `json:"website"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Transaction struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // income, expense
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	RequestID   *uint     `json:"requestId" gorm:"column:request_id"`
	CreatedAt   time.Time `json:"createdAt"`
}

type FiscalConfig struct {
	ID                  uint   `json:"id" gorm:"primaryKey"`
	CompanyID           uint   `json:"companyId" gorm:"column:company_id"`
	CertificatePath     string `json:"certificatePath" gorm:"column:certificate_path"`
	CertificatePassword string `json:"certificatePassword" gorm:"column:certificate_password"`
	Environment         string `json:"environment"` // homologacao, producao
	RegimeTributario    string `json:"regimeTributario" gorm:"column:regime_tributario"`
}

type AuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"userId" gorm:"column:user_id"`
	Entity    string    `json:"entity"`
	Action    string    `json:"action"` // CREATE, UPDATE, DELETE
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"createdAt"`
}

type Setting struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Key   string `json:"key" gorm:"unique;not null"`
	Value string `json:"value"`
}

// OrcamentoItem model (Budget)
type OrcamentoItem struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	RequestID  uint      `json:"requestId" gorm:"column:request_id"`
	Descricao  string    `json:"descricao"`
	Quantidade float64   `json:"quantidade"`
	ValorUnit  float64   `json:"valorUnit" gorm:"column:valor_unit"`
	Tipo       string    `json:"tipo"` // peca, mao_de_obra
	Aprovado   bool      `json:"orcamentoAprovado" gorm:"column:aprovado;default:false"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Assinatura model
type Assinatura struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RequestID uint      `json:"requestId" gorm:"column:request_id"`
	Tipo      string    `json:"tipo"`                        // cliente, tecnico
	DataURL   string    `json:"assinatura" gorm:"type:text"` // Base64 image
	CreatedAt time.Time `json:"createdAt"`
}

// NFSe model (Invoice)
type NFSe struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RequestID    uint      `json:"requestId" gorm:"column:request_id"`
	Numero       string    `json:"numero"`
	Status       string    `json:"status"` // emitido, cancelado, erro
	PDFURL       string    `json:"pdfUrl" gorm:"column:pdf_url"`
	XMLURL       string    `json:"xmlUrl" gorm:"column:xml_url"`
	MotivoCancel string    `json:"motivoCancel" gorm:"column:motivo_cancel"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
