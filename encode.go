package jpnic

import (
	"fmt"
)

func Marshal(input WebTransaction) (string, error) {
	var result string

	// JPNIC Network Information
	result += fmt.Sprintf("WORK_KIND_ID=%s\nIPADDR=%s\nNETWRK_NM=%s\nINFRA_USR_KIND=%s\nORG_NM_JP1=%s\nORG_NM_JP2=%s\nORG_NM_JP3=%s\nORG_NM1=%s\nORG_NM2=%s\nORG_NM3=%s\nZIPCODE=%s\nADDR_JP1=%s\nADDR_JP2=%s\nADDR_JP3=%s\nADDR1=%s\nADDR2=%s\nADDR3=%s\nABUSE=%s\nRYAKUSYO=%s\nNMSRV=%s\nNTFY_MAIL=%s\nPLAN_DATA=%s\nDELI_NO=%s\nRTN_DATE=%s\n",
		input.Network.KindID,
		input.Network.IPAddress,
		input.Network.NetworkName,
		input.Network.InfraUserKind,
		input.Network.OrgJP1,
		input.Network.OrgJP2,
		input.Network.OrgJP3,
		input.Network.Org1,
		input.Network.Org2,
		input.Network.Org3,
		input.Network.ZipCode,
		input.Network.AddrJP1,
		input.Network.AddrJP2,
		input.Network.AddrJP3,
		input.Network.Addr1,
		input.Network.Addr2,
		input.Network.Addr3,
		input.Network.Abuse,
		input.Network.Ryakusyo,
		input.Network.NameServer,
		input.Network.NotifyEmail,
		input.Network.Plan,
		input.Network.DeliNo,
		input.Network.ReturnDate,
	)

	// JPNIC Admin User
	result += fmt.Sprintf("ADM_JPNIC_HDL=%s\nADM_GNAME_JP=%s\nADM_GNAME=%s\nADM_EMAIL=%s\nADM_ORG_NM_JP1=%s\nADM_ORG_NM_JP2=%s\nADM_ORG_NM_JP3=%s\nADM_ORG_NM1=%s\nADM_ORG_NM2=%s\nADM_ORG_NM3=%s\nADM_ZIPCODE=%s\nADM_ADDR_JP1=%s\nADM_ADDR_JP2=%s\nADM_ADDR_JP3=%s\nADM_ADDR1=%s\nADM_ADDR2=%s\nADM_ADDR3=%s\nADM_DIVISION_JP=%s\nADM_DIVISION=%s\nADM_PHONE=%s\nADM_FAX=%s\nADM_NTFY_MAIL=%s\n",
		input.AdminUser.JPNICHandle,
		input.AdminUser.NameJP,
		input.AdminUser.Name,
		input.AdminUser.Email,
		input.AdminUser.OrgJP1,
		input.AdminUser.OrgJP2,
		input.AdminUser.OrgJP3,
		input.AdminUser.Org1,
		input.AdminUser.Org2,
		input.AdminUser.Org3,
		input.AdminUser.ZipCode,
		input.AdminUser.AddrJP1,
		input.AdminUser.AddrJP2,
		input.AdminUser.AddrJP3,
		input.AdminUser.Addr1,
		input.AdminUser.Addr2,
		input.AdminUser.Addr3,
		input.AdminUser.DivisionJP,
		input.AdminUser.Division,
		input.AdminUser.Phone,
		input.AdminUser.Fax,
		input.AdminUser.NotifyMail,
	)

	// JPNIC Tech User
	for count, tech := range input.TechUsers {
		techCount := count + 1
		result += fmt.Sprintf("TECH%d_JPNIC_HDL=%s\nTECH%d_GNAME_JP=%s\nTECH%d_GNAME=%s\nTECH%d_EMAIL=%s\nTECH%d_ORG_NM_JP1=%s\nTECH%d_ORG_NM_JP2=%s\nTECH%d_ORG_NM_JP3=%s\nTECH%d_ORG_NM1=%s\nTECH%d_ORG_NM2=%s\nTECH%d_ORG_NM3=%s\nTECH%d_ZIPCODE=%s\nTECH%d_ADDR_JP1=%s\nTECH%d_ADDR_JP2=%s\nTECH%d_ADDR_JP3=%s\nTECH%d_ADDR1=%s\nTECH%d_ADDR2=%s\nTECH%d_ADDR3=%s\nTECH%d_DIVISION_JP=%s\nTECH%d_DIVISION=%s\nTECH%d_PHONE=%s\nTECH%d_FAX=%s\nTECH%d_NTFY_MAIL=%s\n",
			techCount, tech.JPNICHandle,
			techCount, tech.NameJP,
			techCount, tech.Name,
			techCount, tech.Email,
			techCount, tech.OrgJP1,
			techCount, tech.OrgJP2,
			techCount, tech.OrgJP3,
			techCount, tech.Org1,
			techCount, tech.Org2,
			techCount, tech.Org3,
			techCount, tech.ZipCode,
			techCount, tech.AddrJP1,
			techCount, tech.AddrJP2,
			techCount, tech.AddrJP3,
			techCount, tech.Addr1,
			techCount, tech.Addr2,
			techCount, tech.Addr3,
			techCount, tech.DivisionJP,
			techCount, tech.Division,
			techCount, tech.Phone,
			techCount, tech.Fax,
			techCount, tech.NotifyMail)
	}

	// Etc
	result += fmt.Sprintf("CERT_ID=%s\nPSWD=%s\n", input.Etc.CertID, input.Etc.Password)

	return result, nil
}
