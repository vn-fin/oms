# Account / Order / Signal Management and Trading

## Features
### **Current Features**

#### **1. Basket-Based Trading**
* Users can create custom stock baskets with predefined target positions or weights per stock.
* Support executing a basket **multiple times** (repeatable basket trades).
* Allow **Buy**/**Sell** using **Bid**/**Ask**/**Mid** price selection.
* Support **Fill All** (convert all legs into market orders) or **Cancel All** actions.
* Display real-time futures on the same dashboard and allow selecting future prices for hedging.

#### **2. Basket Hedging**
* Automated hedging using futures (e.g., **VN30F1M**, **VN30F2M**).
* Hedge quantity auto-calculated based on basket notional and multiplier.

#### **3. Account Details & Monitoring**
* Display positions across all symbols.
* Show account balance, fees, PnL, buying power,...
* Display hedge ratio and capital usage.

#### **4. Visualization & Analytics**
* Performance metrics and portfolio summary.
* Index/Futures charts with basis visualization.
* Basket-level PnL and per-symbol contribution.


### **Recommended Enhancements**

#### **1. Multi-Account Trading**
* Support trading across **multiple accounts** with:
  * Independent balances
  * Shared balance pools
  * Account grouping (e.g., Trading, Hedging, Scalping)

#### **2. Strategy & Signal Integration**
* Integrate system-generated strategy signals to determine:
  * Target cash allocation
  * Target exposure
  * Rebalancing actions
* Unified engine to combine signals using StrategyWeight, Exposure, and RiskBudget.

#### **3. Per-Ticker “Fill Now”**
* Allow “**Fill Now**” for **individual tickers** inside a basket instead of only “Fill All”.
* Helps complete partially filled baskets more flexibly.

#### **4. Cash-Based Basket Creation**
* Create baskets using **cash allocation** instead of fixed quantities.
* System auto-converts cash → weight → quantity based on real-time prices.

#### **5. Manual Trades as Offsets**
* All manual trades are treated as **offsets** to strategy-generated target positions.
* Keeps final target position consistent even when trader intervenes manually.

#### **6. Portfolio Booksize Adjustment**
* Allow creating and managing portfolios where the user can:
  * **Increase** booksize (scale in)
  * **Reduce** booksize (scale out)
* Automatically adjusts all positions proportionally.

## Core Concepts

* **Account Credentials**
  * Multiple trading accounts/sub-accounts.
  * Supports shared balance pools and permission controls.

* **Account Groups**
  * Logical grouping of accounts 

* **MaxBookSize**
  * Maximum allowable cash or notional value allocated to a portfolio/book.
  * Used to cap exposure and prevent over-leverage.

* **Exposure**
  * Percentage of book size allocated to a position (e.g., 10%, 70%).
  * Used in risk and position sizing.

* **TradingType**: Defines the execution characteristics of the book/strategy
  * **Daily**: low frequency
  * **Intraday**: medium frequency
  * **HFT**: high frequency

* **TradingWeight**: Defines how the total MaxBookSize is split across trading types, default values:
  * **Daily**: 100%
  * **Intraday**: 0%
  * **HFT**: 0%

* **StrategyWeight**
  * Weight used when combining multiple strategy signals to derive a unified target position.
  * Higher weight means higher influence on final decision.

* **Offset**

  * Adjustment to the target position caused by manual trades.
  * Offset can be considered as **Intraday trading type**.
  * Final position is calculated as:

```
FinalPosition = StrategyCombinedBookSize + Offset
```
where the **StrategyCombinedBookSize** is the sum across trading types:

```
StrategyCombinedBookSize = 
  (TradingWeight_Daily   * BookSize_Daily   * Exposure_Daily   * StrategyWeight_Daily) +
  (TradingWeight_Intraday * BookSize_Intraday * Exposure_Intraday * StrategyWeight_Intraday) +
  (TradingWeight_HFT      * BookSize_HFT      * Exposure_HFT      * StrategyWeight_HFT)
```

## Examples
### Fundamental Trading
#### 1. **MaxBookSize**: 50.000.000.000

And we allocate for 4 stocks: **ACB**, **HPG**, **FPT** and **VIC** with the following weights:

| Stock     | Weight | Allocated Amount |
|-----------| ------ | ---------------- |
| **ACB**   | 20%    | 10,000,000,000   |
| **HPG**   | 20%    | 10,000,000,000   |
| **FPT**   | 20%    | 10,000,000,000   |
| **VIC**   | 40%    | 20,000,000,000   |
| **Total** | 100%   | 50,000,000,000   |

We now illustrate the calculation for **ACB**.

#### 2. **TradingType**, **Exposure** and **TradingWeight**: 

| TradingType | Exposure | TradingWeight <br/>(Total 100%) | BookSize * Exposure * TradingWeight        |
| ----------- |----------|----------------------------|--------------------------------------------|
| Daily       | 70%      | 80%                        | 10,000,000,000 * 0.7 * 0.8 = 5,600,000,000 |
| Intraday    | 100%     | 20%                        | 10,000,000,000 * 1.0 * 0.2 = 2,000,000,000 |

#### 3. Strategy Signal Allocation (after combining multiple strategies):

| TradingType | Signal | Allocated Amount                    |
| ----------- |--------|-------------------------------------|
| Daily       | 60%    | 5,600,000,000 * 0.6 = 3,360,000,000 |
| Intraday    | 70%    | 2,000,000,000 * 0.7 = 1,400,000,000 |

#### 4. Manual Trades (Offset)

* Offset: 1,000,000,000

then, the current holding position: 

* Current: 3,360,000,000 + 1,400,000,000 + 1,000,000,000 = 5,760,000,000

#### 5. Final Position
- The current price of **ACB** is 24,150 VND (as of 2025-12-11 11:20:00), so the target position is: 6,520,000,000 / 24,150 ~ **238,500** shares.

### Future trades:
- Similar calculation applies, but may include HFT trading type and allow short positions (negative weights).