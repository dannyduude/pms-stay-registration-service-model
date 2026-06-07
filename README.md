pms-stay-registration-service

Objective

pms-stay-registration-service is the domain backend for the PMS Stay Registration module.

This service owns the core business logic and persistence for:

* reservations / stays
* persons attached to a reservation
* identity documents such as DNI, NIE and passports
* OCR attempts and parsed document data
* guest intake tokens
* future police TXT export jobs

This service is the only component in this module that connects directly to MongoDB.

Service Type

Domain backend.

Tech Stack

* Go
* MongoDB
* OpenAPI
* Docker
* Tesseract OCR, planned/optional for local OCR execution

Template

This service is based on:

go-printos-backend-quickstart

The implementation should preserve the template conventions around:

* api/openapi.yml
* api/controllers
* api/models/logical
* internal/business
* internal/entities
* internal/repositories
* app/server
* config
* Dockerfile

Responsibilities

This service is responsible for:

* creating and managing reservations
* calculating reservation nights from check-in and check-out dates
* validating reservation fields
* creating and managing persons for a reservation
* storing identity document metadata
* receiving document uploads
* running OCR or OCR stubs
* parsing raw text from DNI, NIE and passport inputs
* creating guest intake tokens
* validating guest intake tokens
* exposing a future police TXT export endpoint as TODO/not implemented

Non-Responsibilities

This service must not:

* render frontend screens
* expose UI-specific view models
* act as a frontend BFF/view-service
* contain React code
* manage finance, payments, cleaning, pricing or availability
* call MongoDB databases owned by other services
* allow frontend code to access MongoDB directly

MongoDB Ownership

MongoDB physical instance: shared platform MongoDB.

Logical database owned by this service:

pms_stay_registration

Collections:

* reservations
* persons
* identity_documents
* ocr_attempts
* intake_tokens
* export_jobs

Only this service should read/write these collections.

Main API Areas

Health:

* GET /health

Reservations:

* POST /api/reservations
* GET /api/reservations
* GET /api/reservations/{reservationId}
* PUT /api/reservations/{reservationId}
* DELETE /api/reservations/{reservationId}

Persons:

* GET /api/reservations/{reservationId}/persons
* POST /api/reservations/{reservationId}/persons
* GET /api/persons/{personId}
* PUT /api/persons/{personId}
* DELETE /api/persons/{personId}

Documents and OCR:

* POST /api/reservations/{reservationId}/documents
* POST /api/reservations/{reservationId}/ocr
* POST /api/reservations/{reservationId}/parse-text

Guest Intake:

* POST /api/reservations/{reservationId}/intake-token
* GET /api/intake-tokens/{token}/reservation

Exports:

* POST /api/exports/police-txt

The police TXT export endpoint is intentionally a placeholder for now.

Business Rules

* checkOut must be after checkIn.
* nights is calculated by the backend.
* A person must always belong to a reservationId.
* OCR must not automatically create or update a person.
* OCR only returns parsed fields.
* The user or admin must review, correct and confirm data before saving a person.
* Police TXT export is not implemented yet and must return a clear TODO/not implemented response.
* Only this service connects to MongoDB.

Environment Variables

* PORT=8081
* MONGO_URI=mongodb://root:rootpassword@mongo:27017/pms_stay_registration?authSource=admin
* MONGO_DATABASE=pms_stay_registration
* UPLOAD_DIR=/app/uploads
* OCR_ENGINE=tesseract

Local Development

Run MongoDB through pms-platform-infra, then start this service locally.

Example:

PORT=8081 MONGO_URI=“mongodb://root:rootpassword@localhost:27017/pms_stay_registration?authSource=admin” MONGO_DATABASE=“pms_stay_registration” go run ./app/server

Docker

This service must provide its own Dockerfile.

The image should include:

* Go service binary
* service config
* optional OCR runtime dependencies
* upload directory support

Architecture Rule

This is a domain backend. It owns business logic and MongoDB persistence for the stay registration domain.
