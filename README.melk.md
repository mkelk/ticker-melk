# Melk's Notes

Personal notes and observations for the Ticker project.

## Setup

<!-- Your local setup notes here -->

## Recipes

### Sync fork with upstream

Pull latest from the parent repo (pengelbrecht/ticker) and merge into your branch:

```bash
# 1. Fetch latest from upstream
git fetch upstream

# 2. Update your main branch
git checkout main
git merge upstream/main   # or: git reset --hard upstream/main (for clean sync)

# 3. Push updated main to your fork
git push origin main

# 4. Merge into your personal branch
git checkout main-melk
git merge main
```

## Things to Remember

<!-- Important things you've learned -->

## Ideas

<!-- Feature ideas, improvements, etc. -->

## Scratch

<!-- Temporary notes, debugging info, etc. -->
