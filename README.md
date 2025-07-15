# GoFr Payment Service

This is a **self-hosted payment processor** built using the [GoFr framework](https://gofr.dev).

---

## 📌 Problem

Modern payment processing requires:
- High security for sensitive payment data.
- High reliability for transaction handling.
- Full audit trails for compliance and reconciliation.

---

## ✅ Solution

This project aims to provide an **open-source, self-hosted payment processor**, designed to integrate with open payment gateways.

Planned features:
- 🔹 **Payment creation & callback endpoints**
- 🔹 **Refund and reconciliation APIs**
- 🔹 **Webhook support for payment updates**
- 🔹 **Integration with open-source payment gateways only**

---

## 🗂️ Project Structure

```bash
.
├── go.mod
├── main.go
├── handler/             # API handlers
├── service/             # Business logic
├── store/               # Data persistence logic
├── migrations/          # DB migrations
├── static/openapi.json  # Swagger docs
├── configs/.env         # Environment configs
├── .golangci.yml        # Linter config
└── README.md
