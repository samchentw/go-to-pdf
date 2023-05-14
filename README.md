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
   ![image](https://github.com/samchentw/go-to-pdf/assets/89454932/471d6112-412d-41b8-876d-f2d4bce7dab9)
2. 將網址及API送出，將回傳PDF網址
   ![image](https://github.com/samchentw/go-to-pdf/assets/89454932/29b5fadd-ea4d-4f1a-854d-0643a35ecf43)
3. 執行回傳結果網址
   ![image](https://github.com/samchentw/go-to-pdf/assets/89454932/aad90d89-e55c-423d-8c62-572c6f443f07)

## 測試

```
go test
```

## 依賴套件

https://github.com/chromedp/chromedp

https://github.com/gin-gonic/gin

https://github.com/google/uuid
