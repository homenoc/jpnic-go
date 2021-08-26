package jpnic

type Result struct {
	Err           error
	ResultErr     []error
	RecepNo       string
	AdmJPNICHdl   string
	Tech1JPNICHdl string
	Tech2JPNICHdl string
}

type WebTransaction struct {
	Network   Network    `json:"network"`
	AdminUser AdminUser  `json:"admin_user"`
	TechUsers []TechUser `json:"tech_users"`
	Etc       Etc        `json:"etc"`
}

type Network struct {
	KindID        string `json:"kind_id"`
	IPAddress     string `json:"ip_address"`
	NetworkName   string `json:"network_name"`
	InfraUserKind string `json:"infra_user_kind"`
	OrgJP1        string `json:"org_jp_1"`
	OrgJP2        string `json:"org_jp_2"`
	OrgJP3        string `json:"org_jp_3"`
	Org1          string `json:"org_1"`
	Org2          string `json:"org_2"`
	Org3          string `json:"org_3"`
	ZipCode       string `json:"zip_code"`
	AddrJP1       string `json:"addr_jp_1"`
	AddrJP2       string `json:"addr_jp_2"`
	AddrJP3       string `json:"addr_jp_3"`
	Addr1         string `json:"addr_1"`
	Addr2         string `json:"addr_2"`
	Addr3         string `json:"addr_3"`
	Abuse         string `json:"abuse"`
	Ryakusyo      string `json:"ryakusho"`
	NameServer    string `json:"name_server"`
	NotifyEmail   string `json:"notify_email"`
	Plan          string `json:"plan"`
	DeliNo        string `json:"deli_no"`
	ReturnDate    string `json:"return_date"`
}

type AdminUser struct {
	JPNICHandle string `json:"jpnic_handle"`
	NameJP      string `json:"name_jp"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	OrgJP1      string `json:"org_jp_1"`
	OrgJP2      string `json:"org_jp_2"`
	OrgJP3      string `json:"org_jp_3"`
	Org1        string `json:"org_1"`
	Org2        string `json:"org_2"`
	Org3        string `json:"org_3"`
	ZipCode     string `json:"zip_code"`
	AddrJP1     string `json:"addr_jp_1"`
	AddrJP2     string `json:"addr_jp_2"`
	AddrJP3     string `json:"addr_jp_3"`
	Addr1       string `json:"addr_1"`
	Addr2       string `json:"addr_2"`
	Addr3       string `json:"addr_3"`
	DivisionJP  string `json:"division_jp"`
	Division    string `json:"division"`
	Phone       string `json:"phone"`
	Fax         string `json:"fax"`
	NotifyMail  string `json:"notify_mail"`
}

type TechUser struct {
	JPNICHandle string `json:"jpnic_handle"`
	NameJP      string `json:"name_jp"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	OrgJP1      string `json:"org_jp_1"`
	OrgJP2      string `json:"org_jp_2"`
	OrgJP3      string `json:"org_jp_3"`
	Org1        string `json:"org_1"`
	Org2        string `json:"org_2"`
	Org3        string `json:"org_3"`
	ZipCode     string `json:"zip_code"`
	AddrJP1     string `json:"addr_jp_1"`
	AddrJP2     string `json:"addr_jp_2"`
	AddrJP3     string `json:"addr_jp_3"`
	Addr1       string `json:"addr_1"`
	Addr2       string `json:"addr_2"`
	Addr3       string `json:"addr_3"`
	DivisionJP  string `json:"division_jp"`
	Division    string `json:"division"`
	Phone       string `json:"phone"`
	Fax         string `json:"fax" `
	NotifyMail  string `json:"notify_mail"`
}

type Etc struct {
	CertID   string `json:"cert_id"`
	Password string `json:"password"`
}
