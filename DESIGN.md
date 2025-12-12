

## Requirements
A. Execution & Strategy Features

1. "Fill All" Mechanism (Simultaneous Execution):
   - Trigger: Add a global "Fill All" button.
   - Logic: Execute the entire Stock Basket and the corresponding Futures Hedge in a single atomic action.
   - Supported Directions:
     - Buy Stock Basket + Short Futures.
     - Sell Stock Basket + Long Futures.
   - Purpose: Minimize slippage by executing basket and hedge simultaneously in one action.slippage.
2. Fixed-Quantity Basket Support:
   - New Mode: Support creating baskets based on specific integer quantities (e.g., `VPB: 1500`) instead of weights.
   - Purpose: To strictly replicate the **E1VFVN30** ETF structure.

B. Access Control & Risk Management

1. Role-Based Access Control (RBAC):
   - Admin: Full access. Can set limits for Traders.
   - Trader: Restricted access. Can only trade within assigned limits.

2. Capital Limit (Risk Check):
   - Admin Action: Admin sets a `capital_limit` (VND) for each Trader.
   - Pre-Trade Check: Before accepting an order: $(\text{Current Position} + \text{Pending} + \text{New Order}) \leq \text{Capital Limit}$
   - Enforcement: Reject order immediately if limit is exceeded.

C. Tracking

1. Manual Order Tracking:
   - Identity: The system must record **Who** placed the manual order (User ID).
   - Timestamp: The system must record **When** the order was placed (Exact Timestamp).
## Database
## Database Schema

### login_credentials
| Field | Type | Description |
|-------|------|-------------|
| 🔑 id | PK | Primary Key |
| user_id | string | User identifier |
| display_name | string | Display name |
| account_type | string | Account type |
| init_cash | decimal | Initial cash |
| updated | timestamp | Last updated |
| balance | decimal | Current balance |
| role | string | User role |

### user_orders
| Field | Type | Description |
|-------|------|-------------|
| 🔑 id | PK | Primary Key |
| 🔗 credential_id | FK | Foreign Key → login_credentials.id |
| user_id | string | User identifier |
| symbol | string | Trading symbol |
| symbol_type | string | Symbol type |
| side | string | Buy/Sell |
| order_type | string | Order type |
| order_price | decimal | Order price |
| quantity | integer | Order quantity |
| filled_qty | integer | Filled quantity |
| remaining_qty | integer | Remaining quantity |
| status | string | Order status |
| created_at | timestamp | Creation time |
| updated_at | timestamp | Last update |

**Relationship:** login_credentials (1) ─── (N) user_orders
## APIs

**prefix**: `/paper-trading/v1/`


### Authentication
All APIs require headers:
- `Authorization`: Bearer token / user identifier (get user_id trong access token)

---

### 1. Order Management

#### Create Order
```
POST /orders

Request:
{
  "credentialId": "xxx",
  "orderType": "L",           // M, L, ATO, ATC
  "price": 100,
  "quantity": 10,
  "side": "B",                // B, S
  "symbol": "SHB",
  "symbolType": "VnStock"     // VnStock, VnFuture
}

{
  "message": "Order created successfully",
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "Pending"         // "Pending" or "Rejected"
}
```
Validation Logic (Trader Role Only):
1. Limit Order (orderType = "L"):
   - Check: quantity × price < capital_limit
   - Example: 1000 × 24.5 = 24,500,000
     → Status = "Pending" if capital_limit >= 24,500,000
     → Status = "Rejected" if capital_limit < 24,500,000

2. Market Order (orderType = "M"):
   - Simulate matching against current orderbook
   - Calculate total_cost from matched prices
   - Check: total_cost < capital_limit
   - Example (Buy 400 shares):
     Orderbook: [100@24.5, 200@24.6, 500@24.7]
     Match: 100×24.5 + 200×24.6 + 100×24.7 = 9,790,000
   - Check: total_cost < capital_limit


3. ATO/ATC Orders (orderType = "ATO" or "ATC"):
   - ATO: Use open_price from snapshot
   - ATC: Use reference_price (floor price) from snapshot
   - Calculate: total_cost = quantity × price
   - Check: total_cost < capital_limit

4. Admin Role:
   - No capital_limit validation
   - Status always = "Pending"

Note: Orders with Status = "Rejected" will not be processed by match engine.
```
{
  "message": "Order created successfully",
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "Pending"         // "Pending" hoặc "Rejected"
}
Error Response (400):
{
  "error": "exceeds capital limit",
  "required_capital": 24500000,
  "capital_limit": 20000000,
  "available": 20000000
}
```
#### Get Order by ID

```
GET /orders/{orderId}?credentialId={credentialId}

Response (200):
{
  "order_id": "550e8400-...",
  "symbol": "SHB",
  "side": "B",
  "order_type": "L",
  "order_price": 100,
  "quantity": 10,
  "filled_qty": 5,
  "remaining_qty": 5,
  "status": "PartiallyFilled",
  "created_at": "2025-12-11T10:00:00Z",
  "updated_at": "2025-12-11T10:05:00Z"
}
```

#### Update Order
```
PUT /orders

Request:
{
  "orderId": "550e8400-...",
  "credentialId": "xxx",
  "quantity": 100
}

Response (200):
{
  "message": "Order updated successfully",
  "old_order_id": "550e8400-...",
  "new_order_id": "660e8400-..."
}
```

#### Cancel Order
```
DELETE /orders

Request:
{
  "OrderID": "550e8400-...",
  "credentialId": "xxx"
}

Response (200):
{
  "message": "Order canceled successfully",
  "order_id": "550e8400-..."
}
```

#### Get Pending Orders
```
GET /pending-orders?credentialId={credentialId}

Response (200):
[
  {
    "order_id": "550e8400-...",
    "symbol": "SHB",
    "side": "B",
    "order_type": "L",
    "order_price": 100,
    "quantity": 10,
    "filled_qty": 0,
    "remaining_qty": 10,
    "status": "Pending",
    "created_at": "2025-12-11T10:00:00Z"
  }
]
```
