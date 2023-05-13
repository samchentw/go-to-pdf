# html 轉 pdf服務

提供輸入網址後轉PDF檔API


## 開發

```
go run main.go
```
預設開啟8080 Port


## 佈署

```
go build
```

## API
1. 提供要轉成PDF的網址
![image](https://github.com/samchentw/go-to-pdf/assets/89454932/562d0c4d-b700-4bb6-acc2-2d311be07799)

2. 將網址及API送出，將回傳PDF網址
![image](https://github.com/samchentw/go-to-pdf/assets/89454932/0abd81d0-5b24-4c47-99eb-94b0e5c5325f)

3. 執行回傳結果網址
![image](https://github.com/samchentw/go-to-pdf/assets/89454932/e4762037-b9b1-4efd-b8fc-ae7abe7a2cce)



## 依賴套件

github.com/chromedp/cdproto

github.com/gin-gonic/gin

github.com/google/uuid
