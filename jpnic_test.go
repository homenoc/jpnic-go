package jpnic

import (
	"testing"
)

//func TestResultProcess(t *testing.T) {
//	var result Result
//
//	str := "<html>\n<body>\nRET=00\nRET_CODE=00000000\nRECEP_NO=001\nRECEP_HSNO=0\nADM_JPNIC_HDL=JP00\nADM_GNAME_JP=Y-Net\nTECH1_JPNIC_HDL=JP00\nTECH1_GNAME_JP=Y-Net\nTECH2_JPNIC_HDL=\nTECH2_GNAME_JP=\nCOUNT=0\n</body>\n</html>"
//	reader := strings.NewReader(str)
//	scanner := bufio.NewScanner(reader)
//
//	success := false
//	var retCode []string
//	ret := "00"
//
//	for scanner.Scan() {
//		// RET
//		if strings.Contains(scanner.Text(), "RET=") {
//			t.Logf("Error: %s", scanner.Text()[4:])
//			ret = scanner.Text()[4:]
//			if scanner.Text()[4:] == "00" {
//				success = true
//			}
//			ret = scanner.Text()[4:]
//		}
//
//		// RET_CODE
//		if strings.Contains(scanner.Text(), "RET_CODE=") {
//			t.Logf("RET_CODE Error: %s", scanner.Text()[9:])
//			retCode = append(retCode, scanner.Text()[9:])
//		}
//
//		// RECEP_NO
//		if strings.Contains(scanner.Text(), "RECEP_NO=") {
//			t.Logf("RECEP_NO=%s", scanner.Text()[9:])
//			result.RecepNo = scanner.Text()[9:]
//		}
//
//		// Admin
//		if strings.Contains(scanner.Text(), "ADM_JPNIC_HDL=") {
//			t.Logf("ADM_JPNIC_HDL=%s", scanner.Text()[14:])
//			result.AdmJPNICHdl = scanner.Text()[14:]
//		}
//
//		// Tech1
//		if strings.Contains(scanner.Text(), "TECH1_JPNIC_HDL=") {
//			t.Logf("TECH1_JPNIC_HDL=%s", scanner.Text()[16:])
//			result.Tech1JPNICHdl = scanner.Text()[16:]
//		}
//
//		// Tech2
//		if strings.Contains(scanner.Text(), "TECH2_JPNIC_HDL=") {
//			t.Logf("TECH2_JPNIC_HDL=%s", scanner.Text()[16:])
//			result.Tech2JPNICHdl = scanner.Text()[16:]
//		}
//
//		t.Log(scanner.Text())
//	}
//
//	// RET
//	if ret != "00" {
//		code, _ := strconv.Atoi(ret)
//		ErrorStatusText(code)
//	}
//
//	// RET_CODE
//	for _, code := range retCode {
//		var errStr string
//		t.Logf("%s", code[4:7])
//
//		// interface
//		if code[4:7] == "000" {
//			code, _ := strconv.Atoi(code[4:7])
//			errStr += ErrorStatusText(code)
//
//		}
//
//		// error genre
//		if code[7:] != "0" {
//			code, _ := strconv.Atoi(code[7:])
//			errStr += "_" + ErrorStatusText(code)
//		}
//	}
//	t.Log(success)
//}

func TestSend(t *testing.T) {
	con := Config{
		URL:          "",
		CertFilePath: "/home/yonedayuto/Documents/HomeNOC/cert/v4-cert.pem",
		KeyFilePath:  "/home/yonedayuto/Documents/HomeNOC/cert/v4-prvkey.pem",
		CAFilePath:   "/home/yonedayuto/Documents/HomeNOC/cert/rootcacert_r3.cer",
	}
	input := WebTransaction{}

	err := con.Send(input)
	if err.Err != nil {
		t.Log(err)
	}
}
