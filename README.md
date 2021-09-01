# JPNIC Web Transaction

## 概要

- IPv4登録
- IPv4変更
- IPv6登録
- IPv6変更
- IPv4情報の閲覧
- IPv6情報の閲覧
- 割当済みIPv4返却申請
- 割当済みIPv6返却申請
- 担当者情報の追加/変更
  
また、詳しい仕様に関してはJPNIC側のトランザクション資料と照らし合わせながら使う必要があります。

## 使用方法(例)

```
	con := Config{
		URL:          "https://[URL]",
		CertFilePath: "/home/[UserName]/cert.pem",
		KeyFilePath:  "/home/[UserName]/key.pem",
		CAFilePath:   "/home/[UserName]/ca.pem",
	}
	
	input := WebTransaction{}
	
	err := con.Send(input)
	if err != nil {
	    log.Println(err)
	}
```

## 未実装機能

- Check機能が未実装
- ResponseのError内容の判別機能が未実装