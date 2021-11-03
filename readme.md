# 帳戶權限管理系統API開發

## 軟體需求

- MySQL

  ​	帳號明細清單儲存

- Golang

  使用gcache存token, 具有有效期

  token驗證通過延長有效期(overwrite用以更新有效期)
  
  部分DB數據也會存gcache, 更新DB時也要更新gcache


## 帳號模板權限表
> 做在前端, 只要把權限值傳給後端保存即可

| 角色 template | 權限 | 來源           |
| ------------- | ---- | -------------- |
| Admin         | 255  | 系統初始化自帶 |
| Manager       | 200  | 由Admin建      |
| Member        | 100  | 由Manager建    |

## 帳號清單表 Accounts

> 後續添加屬性欄位

| 名稱     | datatype | 備註              |
| -------- | -------- | ----------------- |
| id       | char     | PK / UUID  UNIQUE |
| username | char     | zxcvbn  UNIQUE    |
| password | char     |                   |
| email    | char     |                   |



## API清單

### SignUp 建立帳號

> 上層用此API建帳號給下層, 並非user申請自己帳號

`POST /create_account`

| Params     | Type   | Notes                                                        |
| ---------- | ------ | ------------------------------------------------------------ |
| `username` | string | Must be present and unique.                                  |
| `password` | string | Must meet minimum complexity scoring per [zxcvbn](https://blogs.dropbox.com/tech/2012/04/zxcvbn-realistic-password-strength-estimation/). 後端會檢查是否符合規格, 通過後存SHA到DB |
| email      | string |                                                              |
| auth       | int    | 0~255. 數字越大權限越高                                      |

`success`

```json
{
  "result": {
    "id_token": "..."
  }
}
```

`fail`

```json
{
  "result":"fail",
  "errors": "..."
}
```




todo

- [x] 新建帳號的權限必須小於自己
- [ ] 某個權限以上才能用

### Login登入

> 帳密驗證通過後會先將舊token刪掉(如果存在)再建新的, 即一個帳號只會有一組token

`POST /login`

| Params     | Type   | Notes                       |
| ---------- | ------ | --------------------------- |
| `username` | string | Must be present and unique. |
| `password` | string |                             |

`success`

```json
{
  "result": "ok",
  "token": "..."
}

```

`fail`

```json
{
  "result":"fail",
  "error": "..."
}
```

todo

- [x] 將token保存並設定時效, key為token, value應該包含id, auth
- [x] 先刪除同帳號的token

### Logout登出

> 立馬移除此token

  `GET /Logout`

Token驗證 : Header add

```json
{"Authorization":"your-token"}
```

`success`

```json
{"result":"ok"}
```

沒有失敗的情境, token不存在也不用回報錯誤





### Get All Account

> 取得帳號清單

`GET /all_accounts`

Token驗證 : Header add

```json
{"Authorization":"your-token"}
```

`success`

```json
{
  "result": {
    "xxx": "yyy"
  }
}
```

`fail`

```json
{
  "result":"fail",
  "errors": "unauthorized"
}
```

todo

- [ ] 驗證token是否存在及查詢權限
- [ ] 更新token時效

### Update

> 更新權限及其他個人內容

`post /update`

Token驗證 : Header add

```json
{"Authorization":"your-token"}
```

| Params       | Type   | Notes                       |
| ------------ | ------ | --------------------------- |
| `id`         | string | string                      |
| `username`   | string | Must be present and unique. |
| 要修正的欄位 | string | ex: phone:00555589          |

`success`

```json
{"result":"ok"}
```

`fail`

```json
{
  "result":"fail",
  "errors": [
    {"field": "username", "message": "MISSING"},
    {"field": "username", "message": "FORMAT_INVALID"},
    {"field": "username", "message": "TAKEN"},
    {"field": "password", "message": "MISSING"},
    {"field": "password", "message": "INSECURE"}
  ]
}
```



todo

- [ ] 驗證token是否存在及查詢id是否為本人
- [ ] 更新token時效
- [ ] 更新後的權限必須小於當下token

### Change Password

`POST /password`

Token驗證 : Header add

```json
{"Authorization":"your-token"}
```
| Params            | Type   | Notes                                                        |
| ----------------- | ------ | ------------------------------------------------------------ |
| `id`              | string | string                                                       |
| `newPassword`     | string | Must meet minimum complexity scoring per [zxcvbn](https://blogs.dropbox.com/tech/2012/04/zxcvbn-realistic-password-strength-estimation/).後端會檢查是否符合規格, 通過後存SHA到DB |
| `currentPassword` | string | Must exist when changing a password while logged in (not using token) |

`success`

```json
{"result":"ok"}
```

`fail`

```json
{
  "result":"fail",
  "errors": [
    {"field": "username", "message": "MISSING"},
    {"field": "username", "message": "FORMAT_INVALID"},
    {"field": "username", "message": "TAKEN"},
    {"field": "password", "message": "MISSING"},
    {"field": "password", "message": "INSECURE"}
  ]
}
```



todo

- [ ] 驗證token是否存在及查詢id是否為本人
- [ ] 更新token時效

### 初始化密碼(忘記密碼)

> 管理員專用,給重設一組任意密碼

`POST /initpassword`

Token驗證 : Header add

```json
{"Authorization":"your-token"}
```

必須驗證為admin權限

| Params        | Type   | Notes                 |
| ------------- | ------ | --------------------- |
| `id`          | string | string                |
| `newPassword` | string | 任意密碼,不做強度檢查 |

`success`

```json
{"result":"ok"}
```

`fail`

```json
{
  "result":"fail",
  "errors": "unauthorized"
}
```



### Token驗證

> 給其他API驗證token

`GET /tokenverify`

Token驗證 : Header add

```json
{"Authorization":"your-token"}
```

`success`

```json
{"result":"valid","auth":111}
```

`fail`

```json
{"result":"invalid"}
```

- [ ] 更新token時效

