🇷🇺 [Читать на русском языке](README_RU.md)

# 🏥 sdmedik Backend Engine

A specialized e-commerce backend built with Go and Fiber, designed to automate complex price management for medical rehabilitation products (TRU) across Russian regions.

---

## 📌 Table of Contents

1.  [🎯 The Core Problem: Manual Labor](#-the-core-problem-manual-labor)
2.  [⚙️ The Solution: Automated Scraper Engine](#️-the-solution-automated-scraper-engine)
3.  [🏗 Architecture & Data Flow](#-architecture--data-flow)
4.  [🛠 Tech Stack](#-tech-stack)
5.  [💬 Usage Examples](#-usage-examples)
6.  [🚀 Local Setup](#-local-setup)
7.  [⚠️ Handling Unstable Systems & Edge Cases](#-handling-unstable-systems--edge-cases)
8.  [📁 Project Structure](#-project-structure)

---

## 🎯 The Core Problem: Manual Labor

Before this project, **sdmedik** employees spent dozens of hours manually updating prices in the PayKeeper system. The data was sourced from the [KTSN SFR](https://ktsr.sfr.gov.ru/) portal, which lists product categories and certificate-based prices for Technical Rehabilitation Units (TRU) across all 89 regions of Russia.

**Challenges:**
*   **Massive Data Volume**: Thousands of price points (Product × Region × Category).
*   **Human Error**: Manual entry led to critical pricing discrepancies.
*   **Time Consumption**: Prices change frequently, making manual updates impossible to keep in sync.

---

## ⚙️ The Solution: Automated Scraper Engine

The project implements a sophisticated scraping and synchronization system that eliminates manual work entirely.

### 1. ESNSI Integration (TRU Codes)
The system fetches the latest valid TRU codes from the official [ESNSI (GosUslugi)](https://esnsi.gosuslugi.ru/rest/ext/v1/classifiers/10616/data) registry. This ensures that only officially recognized medical equipment codes are processed.

### 2. SFR Scraper (Regional Pricing)
An optimized parser that crawls `ktsr.sfr.gov.ru` to:
*   Extract product categories and individual items.
*   Map prices to specific Russian regions.
*   Sync the gathered data directly with the **PayKeeper** payment gateway.

### 3. "Gentle" Scraping Optimization
The KTSN SFR website is notorious for its instability and frequent crashes. Our engine was specifically re-engineered to prevent "dropping" the source site by:
*   **Strict Rate Limiting**: Intelligent delays between requests.
*   **Sequential Processing**: Avoiding aggressive parallelization that triggers server-side protection or crashes.
*   **State Persistence**: Resuming from where it left off in case of external server failures.

---

## 🏗 Architecture & Data Flow

The project follows a **Modular Monolith** pattern using **Clean Architecture** principles.

### 🔄 Data Synchronization Flow
1.  **Trigger**: Scheduled cron job or manual admin request.
2.  **Fetch Codes**: Connect to ESNSI API to get the current TRU dictionary.
3.  **Crawl SFR**: The Scraper module traverses the SFR catalog using optimized, non-intrusive HTTP requests.
4.  **Transformation**: Data is normalized and mapped to internal product entities.
5.  **Integration**: Updated prices are pushed to **PayKeeper** and the local **PostgreSQL** database.
6.  **Caching**: Final results are cached in **Redis** for high-speed API access.

---

## 🛠 Tech Stack
*   **Backend**: Go (Fiber framework)
*   **Scraping**: Go-Query, specialized HTTP clients with custom retry logic.
*   **Database**: PostgreSQL (GORM)
*   **Caching**: Redis
*   **Auth/Permissions**: Casbin (RBAC) & JWT
*   **External Integrations**: PayKeeper API, SFR Portal, ESNSI API.

---

## 💬 Usage Examples

### 1. Triggering a Price Sync (Admin only)
`POST /api/scraper/sync-prices`
```json
{
  "status": "processing",
  "message": "Synchronization started. Target: SFR Portal + ESNSI. Mode: Optimized (Low Load)."
}
```

### 2. Fetching Regional TRU Price
`GET /api/tru/prices?region_id=77&code=12.1.3`
```json
{
  "region": "Moscow",
  "tru_code": "12.1.3",
  "price_certificate": 45000.00,
  "last_updated": "2024-03-15T10:00:00Z"
}
```

---

## 🚀 Local Setup

### Installation
1.  **Clone & Prepare**:
    ```bash
    git clone https://github.com/Fi44er/sdmed.git && cd sdmed
    cp example.env .env
    ```
2.  **Configure PayKeeper & SFR**: Ensure your `.env` contains valid PayKeeper credentials for synchronization to work.
3.  **Run with Docker**:
    ```bash
    docker-compose up -d
    ```

---

## ⚠️ Handling Unstable Systems & Edge Cases

| Scenario | Our Solution |
| :--- | :--- |
| **SFR Site Downtime** | Exponential backoff retry strategy. The scraper waits for the site to recover without flooding it. |
| **Data Format Change** | Robust HTML parsing with `goquery` and validation of incoming TRU codes against ESNSI. |
| **PayKeeper Timeout** | Implementation of the **Unit of Work** pattern to ensure data consistency between local DB and payment gateway. |
| **Regional Mappings** | Internal mapping table to resolve differences in regional naming between SFR and GosUslugi. |

---

## 📁 Project Structure

*   `internal/module/scraper`: The heart of the automation. Contains logic for parsing and site-friendly HTTP communication.
*   `internal/module/tru`: Logic for handling Technical Rehabilitation Units and regional data.
*   `internal/module/order`: Integration with PayKeeper for final price application.
*   `pkg/process_manager`: Manages long-running background scraping tasks.
