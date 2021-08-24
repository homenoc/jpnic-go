# JPNIC Web Transaction

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