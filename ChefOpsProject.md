# ChefOps Project File  
_Operational Brain • Terminal-First Workflow • SQLite + Go_

This document tracks the full ChefOps CLI initiative — from design decisions to implementation notes — optimized for Neovim, terminal tools, and AI-assisted development.

_Last updated: {{DATE}}_

---

## 🧭 Project Overview

ChefOps is a **terminal-first food operations system** designed for chefs who work heavily on:
- Neovim  
- tmux  
- local LLMs (Ollama, Gemini CLI, OpenCode)  
- iPad/iPhone via iSH + Termius  
- GitHub and markdown files  

The goal: a **portable, dependency-free** CLI + SQLite database that calculates food cost, market lists, and recipe yields reliably across macOS, Linux, Windows, and iOS/iSH.

---

## 🚦 Project Status Snapshot

**Phase:** Planning / Early Implementation  
**Core Decisions:**  
- SQLite as the backbone for recipe + cost data ✔  
- Go (Golang) as the CLI language ✔  
- Neovim + terminal integration ✔  
- AI-assisted docs (Gemini/Ollama/OpenCode) ✔  
- Cross-platform binaries (macOS/Windows/Linux/iSH) planned  

**Next actionable step:**  
Define the CLI command structure + initial database schema.

---

## 📌 Key Decisions (Summary)

### 1️⃣ SQLite chosen as backend  
- Supports generated columns  
- Supports views (perfect for market lists)  
- Stable, fast, scalable  
- 100% terminal-native  
- Perfect for GitHub versioning

### 2️⃣ Go chosen for the CLI  
Reasons:  
- Single static binary  
- No dependency hell  
- Works in iSH on iPad  
- Very stable for long-term use  
- Better for portability than Python

### 3️⃣ Workflow focus  
- Markdown for all documentation  
- Neovim for editing  
- Tmux for multi-pane work  
- OpenCode & Gemini for AI helpers  
- Ollama for local reasoning + analysis  
- Use markdown tables for recipes & cost views  
- Export everything to md for blog or ChefOps docs

---

## 🏗 Project Architecture (Planned)

### Folder structure

chefops/
├── ChefsOpsProject.md
├── db/
│   ├── chefops.db             (auto-generated)
│   ├── schema.sql
│   ├── views.sql
│   └── seed.sql
├── cmd/
│   └── chefops/               (Go CLI entry)
├── internal/
│   ├── sqlite/                (SQL helpers)
│   ├── recipes/
│   ├── ingredients/
│   └── export/
└── out/
├── marketlist.md
└── recipe_costs.md

---

## 🧩 Database Model (High-Level)

### **Ingredients**

id | name | unit | cost_per_unit

### **Recipes**

id | name | yield_qty

### **Recipe Items**

recipe_id | ingredient_id | qty | total_cost (generated)

### **Views**
- `recipe_cost`
- `market_list`
- `recipe_items_expanded`

---

## 🎛 CLI Command Concept

### Ingredient management

chefops ingredient add “Frozen Lobster Meat” –unit kg –cost 150
chefops ingredient list
chefops ingredient update 12 –cost 145

### Recipe building

chefops recipe new “Lobster Mac”
chefops recipe add-item “Lobster Mac” “Frozen Lobster Meat” –qty 0.180
chefops recipe cost “Lobster Mac”

### Market list

chefops marketlist
chefops marketlist –export out/marketlist.md

### Export / AI integration

chefops export recipe –markdown out/recipes.md
chefops export marketlist –json | gemini -i -

---

## 📋 Development Roadmap

### **Phase 1 — Foundation (Current)**
- [ ] Define Go module layout  
- [ ] Write schema.sql + views.sql  
- [ ] Initialize SQLite with seed data  
- [ ] Build minimal CLI with `cobra`  
- [ ] Implement commands:
  - ingredients add/list
  - recipes add/list
  - recipe items add/list
  - basic cost calculator

### **Phase 2 — Market List Engine**
- [ ] Build market list view  
- [ ] Add CLI command for consolidated list  
- [ ] Export to markdown + json  
- [ ] AI workflows for “review this market list”

### **Phase 3 — Terminal UI (optional)**
Using BubbleTea:
- [ ] Search ingredients  
- [ ] Edit recipes interactively  
- [ ] Preview cost calculations live  

### **Phase 4 — Distribution**
- [ ] macOS ARM64 build  
- [ ] macOS Intel build  
- [ ] Linux AMD64 build  
- [ ] Windows build  
- [ ] iSH build (386)  
- [ ] GitHub Releases automation  

---

## 📚 AI Helper Prompts (for Gemini / Ollama / OpenCode)

### **Ask AI to extend a feature**

Extend the ChefOps CLI concept from ChefsOpsProject.md.
Propose improvements to the recipe cost engine and suggest optimized SQL views.

### **Ask AI to explain code**

Explain the following Go function in the context of the ChefOps CLI project.
Refer to the architectural decisions in ChefsOpsProject.md.

### **Ask AI to write code**

Write a Go function for the ChefOps CLI that inserts a new ingredient into the SQLite DB.
The schema is defined in ChefsOpsProject.md.
Return clean, idiomatic Go code.

### **Ask AI to generate docs**

Generate a Markdown guide for new chefs using the ChefOps CLI.
Use ChefsOpsProject.md as the technical reference.

---

## 🗂 Notes & Scratchpad

Keep quick thoughts here:

-  
-  
-  

---

## 🧲 Useful Commands (Neovim + Terminal)

Open this file in a split:
:vs ChefsOpsProject.md

Run AI on selected text:

:’<,’>w !ollama run gemma2:latest

Open SQLite shell:

sqlite3 db/chefops.db

Compile CLI:

go build -o chefops ./cmd/chefops

---

## ✔ Checklist for This Week

- [ ] Initial schema  
- [ ] Initial CLI skeleton  
- [ ] Basic ingredient logic  
- [ ] Basic recipe logic  
- [ ] Export commands  
- [ ] Markdown preview tests  

---

## 🏁 End of File
_This project file is meant to evolve with your workflow. Update as you go._


