package jpnic

type Result struct {
	Err           error
	ResultErr     []error
	RecepNo       string
	AdmJPNICHdl   string
	Tech1JPNICHdl string
	Tech2JPNICHdl string
	Response      string
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

type InfoDetail struct {
	IPAddress            string `json:"ip_address"`
	Ryakusho             string `json:"ryakusho"`
	Type                 string `json:"type"`
	InfraUserKind        string `json:"infra_user_kind"`
	NetworkName          string `json:"network_name"`
	Org                  string `json:"org"`
	OrgEn                string `json:"org_en"`
	PostCode             string `json:"post_code"`
	Address              string `json:"address"`
	AddressEn            string `json:"address_en"`
	AdminJPNICHandle     string `json:"admin_jpnic_handle"`
	AdminJPNICHandleLink string `json:"admin_jpnic_handle_link"`
	TechJPNICHandle      string `json:"tech_jpnic_handle"`
	TechJPNICHandleLink  string `json:"tech_jpnic_handle_link"`
	NameServer           string `json:"name_server"`
	DSRecord             string `json:"ds_record"`
	NotifyAddress        string `json:"notify_address"`
	DeliNo               string `json:"deli_no"`
	RecepNo              string `json:"recep_no"`
	AssignDate           string `json:"assign_date"`
	ReturnDate           string `json:"return_date"`
	UpdateDate           string `json:"update_date"`
}

type JPNICHandleDetail struct {
	IsJPNICHandle bool   `json:"is_jpnic_handle"`
	JPNICHandle   string `json:"jpnic_handle"`
	Name          string `json:"name"`
	NameEn        string `json:"name_en"`
	Email         string `json:"email"`
	Org           string `json:"org"`
	OrgEn         string `json:"org_en"`
	Division      string `json:"division"`
	DivisionEn    string `json:"division_en"`
	Title         string `json:"title"`
	TitleEn       string `json:"title_en"`
	Tel           string `json:"tel"`
	Fax           string `json:"fax"`
	NotifyAddress string `json:"notify_address"`
	UpdateDate    string `json:"update_date"`
}

type InfoIPv4 struct {
	IPAddress   string `json:"ip_address"`
	DetailLink  string `json:"detail_link"`
	Size        string `json:"size"`
	NetworkName string `json:"network_name"`
	AssignDate  string `json:"assign_date"`
	ReturnDate  string `json:"return_date"`
	OrgName     string `json:"org_name"`
	Ryakusho    string `json:"ryakusho"`
	RecepNo     string `json:"recep_no"`
	DeliNo      string `json:"deli_no"`
	Type        string `json:"type"`
	KindID      string `json:"kind_id"`
}

type InfoIPv6 struct {
	IPAddress   string `json:"ip_address"`
	DetailLink  string `json:"detail_link"`
	NetworkName string `json:"network_name"`
	AssignDate  string `json:"assign_date"`
	ReturnDate  string `json:"return_date"`
	OrgName     string `json:"org_name"`
	Ryakusho    string `json:"ryakusho"`
	RecepNo     string `json:"recep_no"`
	DeliNo      string `json:"deli_no"`
	KindID      string `json:"kind_id"`
}

type ReturnIPv6List struct {
	NetworkID     string `json:"network_id"`
	IPAddress     string `json:"ip_address"`
	NetworkName   string `json:"network_name"`
	InfraUserKind string `json:"infra_user_kind"`
	AssignDate    string `json:"assign_date"`
}
