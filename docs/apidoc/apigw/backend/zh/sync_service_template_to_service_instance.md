### 描述

同步服务模板信息到相应的服务实例上(版本：v3.12.3+，权限：服务实例的的创建、编辑、删除权限)

### 输入参数

| 参数名称                | 参数类型  | 必选 | 描述        |
|---------------------|-------|----|-----------|
| bk_biz_id           | int   | 是  | 业务ID      |
| service_template_id | int   | 是  | 服务模板ID    |
| bk_module_ids       | array | 是  | 待同步模块ID列表 |

### 调用示例

```json
{
  "bk_biz_id": 3,
  "service_template_id": 1,
  "bk_module_ids": [
    28
  ]
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |
