### 功能描述

更新产品属性（支持部分字段更新）

### 请求参数

#### 参数说明

| 字段         | 类型                | 必选  | 参数位置 | 描述                              |
|--------------|--------------------|-------|----------|-----------------------------------|
| product_id   | integer            | 是    | 路径参数 | 需要更新的产品ID                  |
| description  | string (nullable)  | 否    | body参数 | 产品描述（不传则保留原值）         |
| stock        | integer (nullable) | 否    | body参数 | 库存数量（不传则保留原值）         |
| type         | string (nullable)  | 否    | body参数 | 产品类型（不传则保留原值）         |

#### 请求示例

```json
{
    "description": "新一代智能设备",
    "stock": 150,
    "type": "electronics"
}
