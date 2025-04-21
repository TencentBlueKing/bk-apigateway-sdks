### Function description

Update product attributes (supports partial field updates)

### Request parameters

#### Parameter description

| Field | Type | Required | Parameter location | Description |
|--------------|--------------------|-------|----------|-----------------------------------|
| product_id | integer | Yes | Path parameter | Product ID to be updated |
| description | string (nullable) | No | Body parameter | Product description (keep the original value if not passed) |
| stock | integer (nullable) | No | Body parameter | Stock quantity (keep the original value if not passed) |
| type | string (nullable) | No | Body parameter | Product type (keep the original value if not passed) |

#### Request example

```json
{
"description": "New generation of smart devices",
"stock": 150,
"type": "electronics"
}
