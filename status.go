package jpnic

// StatusCodeは+1000をしているため、注意が必要
const (
	// 追加
	IPv4Register = 10
	IPv4Edit     = 11
	IPv6Register = 20
	IPv6Edit     = 21

	Infra         = 1
	User          = 2
	Reassignment1 = 3
	Reassignment2 = 4

	// エラーコードやステータスコード
	IPRegistrySystemError                = 1010
	NoCertificateOrUnableToGetMentorCode = 1020
	AuthorityError                       = 1021
	ApplicationProcessingError           = 1030
	InadequateParameters                 = 1099

	MissingRequiredFieldsError                    = 1001
	ExceedsStringError                            = 1002
	ViolationOfTypeError                          = 1003
	InadequateContentFormatError                  = 1004
	InadequateContentExistenceError               = 1005
	InadequateContentMultipleNetworksNotSpecified = 1006
	InadequateContentEtc                          = 1007

	NetWorkAndKindIDError        = 1101
	NetworkAndIPAddressError     = 1102
	NetworkAndNetworkNameError   = 1103
	NetworkAndInfraUserKindError = 1104
	NetworkAndOrgJP1Error        = 1105
	NetworkAndOrgJP2Error        = 1106
	NetworkAndOrgJP3Error        = 1107
	NetworkAndOrg1Error          = 1108
	NetworkAndOrg2Error          = 1109
	NetworkAndOrg3Error          = 1110
	NetworkAndZipCodeError       = 1111
	NetworkAndAddrJP1Error       = 1112
	NetworkAndAddrJP2Error       = 1113
	NetworkAndAddrJP3Error       = 1114
	NetworkAndAddr1Error         = 1115
	NetworkAndAddr2Error         = 1116
	NetworkAndAddr3Error         = 1117
	NetworkAndAbuseError         = 1118
	NetworkAndRyakusyoError      = 1119
	NetworkAndNameServerError    = 1120
	NetworkAndNotifyEmailError   = 1121
	NetworkAndPlanError          = 1122
	NetworkAndDeliNoError        = 1123
	NetworkAndReturnDateError    = 1124
	AdminAndJPNICHandleError     = 1200
	AdminAndNameJPError          = 1201
	AdminAndNameError            = 1202
	AdminAndEmailError           = 1203
	AdminAndOrgJP1Error          = 1204
	AdminAndOrgJP2Error          = 1205
	AdminAndOrgJP3Error          = 1206
	AdminAndOrg1Error            = 1207
	AdminAndOrg2Error            = 1208
	AdminAndOrg3Error            = 1209
	AdminAndZipCodeError         = 1210
	AdminAndAddrJP1Error         = 1211
	AdminAndAddrJP2Error         = 1212
	AdminAndAddrJP3Error         = 1213
	AdminAndAddr1Error           = 1214
	AdminAndAddr2Error           = 1215
	AdminAndAddr3Error           = 1216
	AdminAndDivisionJPError      = 1217
	AdminAndDivisionError        = 1218
	AdminAndPhoneError           = 1219
	AdminAndFaxError             = 1220
	AdminAndNotifyMailError      = 1221
	Tech1AndJPNICHandleError     = 1300
	Tech1AndNameJPError          = 1301
	Tech1AndNameError            = 1302
	Tech1AndEmailError           = 1303
	Tech1AndOrgJP1Error          = 1304
	Tech1AndOrgJP2Error          = 1305
	Tech1AndOrgJP3Error          = 1306
	Tech1AndOrg1Error            = 1307
	Tech1AndOrg2Error            = 1308
	Tech1AndOrg3Error            = 1309
	Tech1AndZipCodeError         = 1310
	Tech1AndAddrJP1Error         = 1311
	Tech1AndAddrJP2Error         = 1312
	Tech1AndAddrJP3Error         = 1313
	Tech1AndAddr1Error           = 1314
	Tech1AndAddr2Error           = 1315
	Tech1AndAddr3Error           = 1316
	Tech1AndDivisionJPError      = 1317
	Tech1AndDivisionError        = 1318
	Tech1AndPhoneError           = 1319
	Tech1AndFaxError             = 1320
	Tech1AndNotifyMailError      = 1321
	Tech2AndJPNICHandleError     = 1400
	Tech2AndNameJPError          = 1401
	Tech2AndNameError            = 1402
	Tech2AndEmailError           = 1403
	Tech2AndOrgJP1Error          = 1404
	Tech2AndOrgJP2Error          = 1405
	Tech2AndOrgJP3Error          = 1406
	Tech2AndOrg1Error            = 1407
	Tech2AndOrg2Error            = 1408
	Tech2AndOrg3Error            = 1409
	Tech2AndZipCodeError         = 1410
	Tech2AndAddrJP1Error         = 1411
	Tech2AndAddrJP2Error         = 1412
	Tech2AndAddrJP3Error         = 1413
	Tech2AndAddr1Error           = 1414
	Tech2AndAddr2Error           = 1415
	Tech2AndAddr3Error           = 1416
	Tech2AndDivisionJPError      = 1417
	Tech2AndDivisionError        = 1418
	Tech2AndPhoneError           = 1419
	Tech2AndFaxError             = 1420
	Tech2AndNotifyMailError      = 1421
	EtcCertIDError               = 1501
	EtcPasswordError             = 1502
)

var statusText = map[int]string{
	// 追加
	IPv4Register: "IPv4登録",
	IPv4Edit:     "IPv4変更",
	IPv6Register: "IPv6登録",
	IPv6Edit:     "IPv6変更",

	Infra:         "インフラ",
	User:          "ユーザ",
	Reassignment1: "再割り振り",
	Reassignment2: "再割り当て",

	// エラーコードやステータスコード
	IPRegistrySystemError:                "IPレジストリシステム内異常",
	NoCertificateOrUnableToGetMentorCode: "証明書無し、メンテナーコード取得不可",
	AuthorityError:                       "権限エラー",
	ApplicationProcessingError:           "申請処理エラー",
	InadequateParameters:                 "パラメータ不備",

	MissingRequiredFieldsError:                    "必須項目欠落",
	ExceedsStringError:                            "文字列超過",
	ViolationOfTypeError:                          "文字列種別違反",
	InadequateContentFormatError:                  "内容不備_フォーマットエラー",
	InadequateContentExistenceError:               "内容不備_存在エラー",
	InadequateContentMultipleNetworksNotSpecified: "内容不備_複数ネットワーク特定不可",
	InadequateContentEtc:                          "内容不備_その他",

	NetWorkAndKindIDError:        "ネットワーク情報(業務区分)",
	NetworkAndIPAddressError:     "ネットワーク情報(IPネットワークアドレス)",
	NetworkAndNetworkNameError:   "ネットワーク情報(ネットワーク名)",
	NetworkAndInfraUserKindError: "ネットワーク情報(インフラ・ユーザ区分)",
	NetworkAndOrgJP1Error:        "ネットワーク情報(組織名1(JP))",
	NetworkAndOrgJP2Error:        "ネットワーク情報(組織名2(JP))",
	NetworkAndOrgJP3Error:        "ネットワーク情報(組織名3(JP))",
	NetworkAndOrg1Error:          "ネットワーク情報(組織名1(EN))",
	NetworkAndOrg2Error:          "ネットワーク情報(組織名2(EN))",
	NetworkAndOrg3Error:          "ネットワーク情報(組織名3(EN))",
	NetworkAndZipCodeError:       "ネットワーク情報(郵便番号)",
	NetworkAndAddrJP1Error:       "ネットワーク情報(住所1(JP))",
	NetworkAndAddrJP2Error:       "ネットワーク情報(住所2(JP))",
	NetworkAndAddrJP3Error:       "ネットワーク情報(住所3(JP))",
	NetworkAndAddr1Error:         "ネットワーク情報(住所1(EN))",
	NetworkAndAddr2Error:         "ネットワーク情報(住所2(EN))",
	NetworkAndAddr3Error:         "ネットワーク情報(住所3(EN))",
	NetworkAndAbuseError:         "ネットワーク情報(Abuse)",
	NetworkAndRyakusyoError:      "ネットワーク情報(会員略称)",
	NetworkAndNameServerError:    "ネットワーク情報(NameServer)",
	NetworkAndNotifyEmailError:   "ネットワーク情報(通知アドレス)",
	NetworkAndPlanError:          "ネットワーク情報(Plan)",
	NetworkAndDeliNoError:        "ネットワーク情報(審議番号)",
	NetworkAndReturnDateError:    "ネットワーク情報(返却日)",
	AdminAndJPNICHandleError:     "管理者連絡窓口(JPNICハンドル名)",
	AdminAndNameJPError:          "管理者連絡窓口(名前(JP))",
	AdminAndNameError:            "管理者連絡窓口(名前(EN))",
	AdminAndEmailError:           "管理者連絡窓口(メールアドレス)",
	AdminAndOrgJP1Error:          "管理者連絡窓口(Org1(JP))",
	AdminAndOrgJP2Error:          "管理者連絡窓口(Org2(JP))",
	AdminAndOrgJP3Error:          "管理者連絡窓口(Org3(JP))",
	AdminAndOrg1Error:            "管理者連絡窓口(Org1(EN))",
	AdminAndOrg2Error:            "管理者連絡窓口(Org2(EN))",
	AdminAndOrg3Error:            "管理者連絡窓口(Org3(EN))",
	AdminAndZipCodeError:         "管理者連絡窓口(郵便番号)",
	AdminAndAddrJP1Error:         "管理者連絡窓口(住所1(JP))",
	AdminAndAddrJP2Error:         "管理者連絡窓口(住所2(JP))",
	AdminAndAddrJP3Error:         "管理者連絡窓口(住所3(JP))",
	AdminAndAddr1Error:           "管理者連絡窓口(住所1(EN))",
	AdminAndAddr2Error:           "管理者連絡窓口(住所2(EN))",
	AdminAndAddr3Error:           "管理者連絡窓口(住所3(EN))",
	AdminAndDivisionJPError:      "管理者連絡窓口(部署(JP))",
	AdminAndDivisionError:        "管理者連絡窓口(部署(EN))",
	AdminAndPhoneError:           "管理者連絡窓口(電話番号)",
	AdminAndFaxError:             "管理者連絡窓口(Fax)",
	AdminAndNotifyMailError:      "管理者連絡窓口(通知アドレス)",
	Tech1AndJPNICHandleError:     "技術連絡担当者1(JPNICハンドル名)",
	Tech1AndNameJPError:          "技術連絡担当者1(名前(JP))",
	Tech1AndNameError:            "技術連絡担当者1(名前(EN))",
	Tech1AndEmailError:           "技術連絡担当者1(メールアドレス)",
	Tech1AndOrgJP1Error:          "技術連絡担当者1(Org1(JP))",
	Tech1AndOrgJP2Error:          "技術連絡担当者1(Org2(JP))",
	Tech1AndOrgJP3Error:          "技術連絡担当者1(Org3(JP))",
	Tech1AndOrg1Error:            "技術連絡担当者1(Org1(EN))",
	Tech1AndOrg2Error:            "技術連絡担当者1(Org2(EN))",
	Tech1AndOrg3Error:            "技術連絡担当者1(Org3(EN))",
	Tech1AndZipCodeError:         "技術連絡担当者1(郵便番号)",
	Tech1AndAddrJP1Error:         "技術連絡担当者1(住所1(JP))",
	Tech1AndAddrJP2Error:         "技術連絡担当者1(住所2(JP))",
	Tech1AndAddrJP3Error:         "技術連絡担当者1(住所3(JP))",
	Tech1AndAddr1Error:           "技術連絡担当者1(住所1(EN))",
	Tech1AndAddr2Error:           "技術連絡担当者1(住所2(EN))",
	Tech1AndAddr3Error:           "技術連絡担当者1(住所3(EN))",
	Tech1AndDivisionJPError:      "技術連絡担当者1(部署(JP))",
	Tech1AndDivisionError:        "技術連絡担当者1(部署(EN))",
	Tech1AndPhoneError:           "技術連絡担当者1(電話番号)",
	Tech1AndFaxError:             "技術連絡担当者1(Fax)",
	Tech1AndNotifyMailError:      "技術連絡担当者1(通知アドレス)",
	Tech2AndJPNICHandleError:     "技術連絡担当者2(JPNICハンドル名)",
	Tech2AndNameJPError:          "技術連絡担当者2(名前(JP))",
	Tech2AndNameError:            "技術連絡担当者2(名前(EN))",
	Tech2AndEmailError:           "技術連絡担当者2(メールアドレス)",
	Tech2AndOrgJP1Error:          "技術連絡担当者2(Org1(JP))",
	Tech2AndOrgJP2Error:          "技術連絡担当者2(Org2(JP))",
	Tech2AndOrgJP3Error:          "技術連絡担当者2(Org3(JP))",
	Tech2AndOrg1Error:            "技術連絡担当者2(Org1(EN))",
	Tech2AndOrg2Error:            "技術連絡担当者2(Org2(EN))",
	Tech2AndOrg3Error:            "技術連絡担当者2(Org3(EN))",
	Tech2AndZipCodeError:         "技術連絡担当者2(郵便番号)",
	Tech2AndAddrJP1Error:         "技術連絡担当者2(住所1(JP))",
	Tech2AndAddrJP2Error:         "技術連絡担当者2(住所2(JP))",
	Tech2AndAddrJP3Error:         "技術連絡担当者2(住所3(JP))",
	Tech2AndAddr1Error:           "技術連絡担当者2(住所1(EN))",
	Tech2AndAddr2Error:           "技術連絡担当者2(住所2(EN))",
	Tech2AndAddr3Error:           "技術連絡担当者2(住所3(EN))",
	Tech2AndDivisionJPError:      "技術連絡担当者2(部署(JP))",
	Tech2AndDivisionError:        "技術連絡担当者2(部署(EN))",
	Tech2AndPhoneError:           "技術連絡担当者2(電話番号)",
	Tech2AndFaxError:             "技術連絡担当者2(Fax)",
	Tech2AndNotifyMailError:      "技術連絡担当者2(通知アドレス)",
	EtcCertIDError:               "認証ID",
	EtcPasswordError:             "パスワード",
}

// ErrorStatusの場合はcodeを自動で+1000する
func ErrorStatusText(code int) string {
	code += 1000
	return statusText[code]
}

func StatusText(code int) string {
	return statusText[code]
}
