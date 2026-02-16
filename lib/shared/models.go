package shared

import (
	"time"
)

// User model
type User struct {
	ID                 string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name               string    `json:"name"`
	Email              string    `json:"email" gorm:"unique;not null"`
	PasswordHash       string    `json:"-" gorm:"column:password_hash"`
	Role               string    `json:"role"` // ADMIN_SISTEMA, PRESTADOR, TECNICO, CLIENTE
	CompanyID          *string   `json:"companyId" gorm:"column:company_id"`
	Active             bool      `json:"active" gorm:"default:true"`
	MustChangePassword bool      `json:"mustChangePassword" gorm:"column:must_change_password;default:false"`
	AvatarURL          *string   `json:"avatarUrl" gorm:"column:avatar_url"`
	Phone              *string   `json:"phone"`
	CreatedAt          time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}

// Client model
type Client struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name      string    `json:"name" gorm:"not null"`
	CNPJ      string    `json:"cnpj" gorm:"unique"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CompanyID string    `json:"companyId" gorm:"column:company_id"`
	Active    bool      `json:"active" gorm:"default:true"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Client) TableName() string {
	return "clientes"
}

// Request model (Service Request)
type Request struct {
	ID                string     `json:"id" gorm:"primaryKey;type:uuid"`
	Numero            int        `json:"numero" gorm:"column:numero"`
	Title             string     `json:"title" gorm:"not null"`
	Description       string     `json:"description"`
	Status            string     `json:"status" gorm:"default:ABERTA"` // ABERTA, PENDENTE, etc.
	Priority          string     `json:"priority"`                     // BAIXA, MEDIA, etc.
	ClientID          string     `json:"clientId" gorm:"column:client_id"`
	ClientName        string     `json:"clientName" gorm:"column:client_name"`
	CompanyID         string     `json:"companyId" gorm:"column:company_id"`
	CreatedBy         string     `json:"createdBy" gorm:"column:created_by"`
	AssignedTo        *string    `json:"assignedTo" gorm:"column:assigned_to"`
	Observation       string     `json:"observation" gorm:"column:observation"`
	MaterialsUsed     string     `json:"materialsUsed" gorm:"column:materials_used"`
	NextMaintenanceAt *time.Time `json:"nextMaintenanceAt" gorm:"column:next_maintenance_at"`
	ScheduledAt       *time.Time `json:"scheduledAt" gorm:"column:scheduled_at"`
	SlaLimit          *time.Time `json:"slaLimit" gorm:"column:sla_limit"`
	ConfirmedAt       *time.Time `json:"confirmedAt" gorm:"column:confirmed_at"`
	ConfirmedBy       *string    `json:"confirmedBy" gorm:"column:confirmed_by"`
	LockedBy          *string    `json:"lockedBy" gorm:"column:locked_by"`
	LockedAt          *time.Time `json:"lockedAt" gorm:"column:locked_at"`
	ValorOrcamento    float64    `json:"valorOrcamento" gorm:"column:valor_orcamento"`
	OrcamentoAprovado bool       `json:"orcamentoAprovado" gorm:"column:orcamento_aprovado"`
	AssinaturaCliente string     `json:"assinaturaCliente" gorm:"column:assinatura_cliente"`
	AssinaturaTecnico string     `json:"assinaturaTecnico" gorm:"column:assinatura_tecnico"`
	DataAssinatura    *time.Time `json:"dataAssinatura" gorm:"column:data_assinatura"`
	PreventiveDone    bool       `json:"preventiveDone" gorm:"column:preventive_done"`
	CreatedAt         time.Time  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         time.Time  `json:"updatedAt" gorm:"column:updated_at"`
}

func (Request) TableName() string {
	return "solicitacoes"
}

type Equipment struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name         string    `json:"name" gorm:"not null"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	SerialNumber string    `json:"serialNumber" gorm:"column:serial_number"`
	ClientID     string    `json:"clientId" gorm:"column:client_id"`
	Active       bool      `json:"active" gorm:"default:true"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Equipment) TableName() string {
	return "equipamentos"
}

type ChecklistItem struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid"`
	RequestID     string    `json:"requestId" gorm:"column:solicitacao_id"`
	EquipamentoID string    `json:"equipamentoId" gorm:"column:equipamento_id"`
	Description   string    `json:"description"`
	Checked       bool      `json:"checked" gorm:"default:false"`
	Observation   string    `json:"observation"`
	CheckedByID   string    `json:"checkedById" gorm:"column:checked_by_id"`
	CheckedByName string    `json:"checkedByName" gorm:"column:checked_by_name"`
	CheckedAt     time.Time `json:"checkedAt" gorm:"column:checked_at"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (ChecklistItem) TableName() string {
	return "checklists"
}

type AgendaEntry struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid"`
	Title        string    `json:"title"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	TechnicianID string    `json:"technicianId" gorm:"column:technician_id"`
	RequestID    *string   `json:"requestId" gorm:"column:request_id"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (AgendaEntry) TableName() string {
	return "agenda"
}

type Company struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	LogoURL   string    `json:"logoUrl" gorm:"column:logo_url"`
	Website   string    `json:"website"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Company) TableName() string {
	return "empresas" // Likely empresas in Portuguese
}

type Transaction struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // income, expense
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	RequestID   *string   `json:"requestId" gorm:"column:request_id"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (Transaction) TableName() string {
	return "financeiro"
}

type FiscalConfig struct {
	ID                  string `json:"id" gorm:"primaryKey;type:uuid"`
	CompanyID           string `json:"companyId" gorm:"column:company_id"`
	CertificatePath     string `json:"certificatePath" gorm:"column:certificate_path"`
	CertificatePassword string `json:"certificatePassword" gorm:"column:certificate_password"`
	Environment         string `json:"environment"` // homologacao, producao
	RegimeTributario    string `json:"regimeTributario" gorm:"column:regime_tributario"`
}

func (FiscalConfig) TableName() string {
	return "fiscal_config"
}

type AuditLog struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	UserID    string    `json:"userId" gorm:"column:user_id"`
	Entity    string    `json:"entity"`
	Action    string    `json:"action"` // CREATE, UPDATE, DELETE
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

type Setting struct {
	ID    string `json:"id" gorm:"primaryKey;type:uuid"`
	Key   string `json:"key" gorm:"unique;not null"`
	Value string `json:"value"`
}

func (Setting) TableName() string {
	return "configuracoes"
}

// OrcamentoItem model (Budget)
type OrcamentoItem struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid"`
	RequestID  string    `json:"requestId" gorm:"column:request_id"`
	Descricao  string    `json:"descricao"`
	Quantidade float64   `json:"quantidade"`
	ValorUnit  float64   `json:"valorUnit" gorm:"column:valor_unit"`
	Tipo       string    `json:"tipo"` // peca, mao_de_obra
	Aprovado   bool      `json:"orcamentoAprovado" gorm:"column:aprovado;default:false"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (OrcamentoItem) TableName() string {
	return "orcamentos"
}

// Assinatura model
type Assinatura struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid"`
	RequestID string    `json:"requestId" gorm:"column:request_id"`
	Tipo      string    `json:"tipo"`                        // cliente, tecnico
	DataURL   string    `json:"assinatura" gorm:"type:text"` // Base64 image
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
}

func (Assinatura) TableName() string {
	return "assinaturas"
}

// NFSe model (Invoice)
type NFSe struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid"`
	RequestID    string    `json:"requestId" gorm:"column:request_id"`
	Numero       string    `json:"numero"`
	Status       string    `json:"status"` // emitido, cancelado, erro
	PDFURL       string    `json:"pdfUrl" gorm:"column:pdf_url"`
	XMLURL       string    `json:"xmlUrl" gorm:"column:xml_url"`
	MotivoCancel string    `json:"motivoCancel" gorm:"column:motivo_cancel"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (NFSe) TableName() string {
	return "nfse"
}
