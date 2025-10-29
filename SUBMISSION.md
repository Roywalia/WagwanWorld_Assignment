# Submission

**Candidate Name**: Pratishth 
**Date**: October 30, 2025  
**Email**: walia.pratishth@gmail.com

---

## Time Spent

**Total: ~7 hours**

- **Frontend (Svelte + UI/UX + Animations)**: 3.5 hours  
- **Backend (Go + PostgreSQL + API)**: 2.5 hours  
- **Bug fixes, Docker, Testing**: 1.0 hour  

---

## What I Built

### 1. Public RSVP Page (`/`)

- **Location**: `src/routes/+page.svelte`
- **Key Features**:
  - Clean, responsive RSVP form with event selection
  - **Real-time "Total RSVPs So Far" counter** (updates instantly after submit)
  - **Fancy animated inputs** with gradient focus glow, scale effect, and floating labels
  - Prevents duplicate RSVPs per email per event
  - Success/error feedback via **SweetAlert2** with branded styling
  - Human-readable event display: `Summer Gala – Downtown Hotel – 20 Aug 25`

- **Design Choices**:
  - Svelte 5 with `$state` and `$effect` for modern reactivity
  - Tailwind CSS for rapid, consistent styling
  - Mobile-first, accessible, and delightful UX

---

### 2. Admin Guest List Manager (`/admin`)

- **Location**: `src/routes/admin/+page.svelte`
- **Key Features**:
  - Full **CRUD** (list, add, delete)
  - **Real-time search** by name or email (backend-powered, case-insensitive)
  - **Status filter** (All / Pending / Attending / Declined)
  - Live stats panel: Total, Attending, Pending, Declined
  - Animated search input with **clear button**
  - Responsive table with status badges and hover effects

---

## Bugs Found & Fixed

| Bug | Problem | Location | Fix | Why |
|-----|-------|--------|-----|-----|
| **1. `phone` missing in type** | TypeScript error: `'phone' does not exist in type 'CreateGuestData'` | `src/lib/api.ts` | Added `phone: string` to `CreateGuestData` | Ensures type safety and allows phone field in API calls |
| **2. Duplicate `GuestHandler`** | Go compile error: `GuestHandler redeclared` | Multiple `guest_handler.go` files | Removed duplicates, kept one clean version | Eliminates build conflicts |
| **3. Filter mismatch** | Frontend sent `?filter=` but backend expected `?status=` | `src/lib/api.ts` + `GetGuests()` | Changed to `?status=` in both | Ensures filters work correctly |

---

## Challenges & Solutions

### Challenge 1: Real-Time Total RSVP Counter
**Problem**: Show total RSVPs across all events, update instantly  
**Solution**:
- Added `rsvps` count to `Event` via SQL `COUNT(g.id)`
- Used `reduce()` to sum in frontend
- Refreshed events after RSVP submit

### Challenge 2: Scalable Search in Admin Panel
**Problem**: Need fast, accurate search by name/email  
**Solution**:
- Backend SQL `LIKE` with `LOWER()` for case-insensitivity
- `URLSearchParams` for clean query building
- `$effect` triggers reload on every keystroke

---

## How to Test My Work

### Setup (Docker – One Command)

```bash
docker compose up --build