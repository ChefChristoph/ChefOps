🍳 ChefOps – Git Workflow & Project Structure

This document describes how to work with the ChefOps codebase using Git. It is designed for fast iteration, safe experimentation, and the ability to always roll back to a known stable version.

ChefOps is a command-line tool written in Go, backed by SQLite, and meant to evolve rapidly — so a clean Git workflow is essential.

⸻

🗂 Project Structure

chefops/
│
├── cmd/chefops/        # main CLI commands
├── internal/           # database logic, migrations, helpers
├── db/                 # SQLite database + SQL views/migrations
├── ChefOpsProject.md   # main project plan
├── go.mod / go.sum     # Go modules
├── export.go           # export logic (markdown/json)
└── README.md           # (this file or primary project README)


⸻

🧰 Branching Strategy

ChefOps uses a “main = stable” model.

main
	•	always runnable
	•	always clean
	•	only merged when stable

feature branches

Every new feature, fix, or experiment is done in its own branch:

git checkout -b feature/<short-name>

Examples:

git checkout -b feature/export-fix
git checkout -b feature/import-csv
git checkout -b feature/nvim
git checkout -b feature/tui

When finished:

git checkout main
git merge feature/<short-name>

Tag the stable version:

git tag v0.3-stable
git push --tags


⸻

🧹 Typical ChefOps Development Workflow

1. Start from clean main

git checkout main
git pull

2. Create a feature branch

git checkout -b feature/<name>

3. Code & Test

go build -o chefops ./cmd/chefops
./chefops <command>

4. Commit in small steps

git add .
git commit -m "message"

5. Merge when stable

git checkout main
git merge feature/<name>

6. Tag a stable snapshot

git tag v0.x-stable
git push --tags

7. Start next feature

git checkout -b feature/<name2>


⸻

🔄 Reverting Mistakes (Safe Reset)

Reset to last stable tag:

git checkout main
git reset --hard v0.3-stable

Discard all working changes:

git reset --hard

Undo last commit but keep changes staged:

git reset --soft HEAD~1

Undo last commit entirely:

git reset --hard HEAD~1


⸻

🚀 Working with GitHub

First push:

git remote add origin git@github.com:ChefChristoph/chefops.git
git push -u origin main

Push new commits:

git push

Push tags:

git push --tags


⸻

🧪 Useful Git Commands

Show changed files:

git diff --name-only

Compare branches:

git diff main..feature/<name>

Show compact commit history:

git log --oneline

Delete a merged branch:

git branch -d feature/<name>


⸻

📓 Suggested .gitignore

chefops
db/*.db
db/*.backup
.idea/
.vscode/
*.swp
*.tmp
*.log


⸻

🎯 Summary

The Git workflow for ChefOps aims for:
	•	fast experimentation
	•	minimal risk
	•	reliable rollbacks
	•	clean stable versions
	•	easy integration with scripts, Neovim, and automation

Using small feature branches + frequent commits will make ChefOps more resilient as it grows.
