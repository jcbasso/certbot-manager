# Git Branching Model

## Overview

This model uses dedicated branches for new features, bug fixes, and major release preparations to keep the `main` branch stable and deployment-ready. All changes are integrated back into `main` through Pull Requests.

## Branches

1.  **`main`**
    *   **Purpose:** Represents the latest stable, production-ready code.
    *   **Rule:** **Do not commit directly to `main`.** All changes must come from merged Pull Requests.

2.  **`feature/`**
    *   **Purpose:** Develop new features (adds functionality).
    *   **Workflow:** Create from `main` -> Develop -> Open PR to `main`.
    *   **Naming:** `feature/<short-description>` (e.g., `feature/user-login`)
    *   **Versioning:** Merging results in a **minor** version increase (e.g., 1.1.0 -> 1.2.0).

```
                       0.1.0-rc2
                0.1.0-rc   │    
                   ╵       ╵    
feature    ● ───── ● ───── ●    
          ╱                 ╲   
main   ─ ● ───────────────── ● ─
                             ╷  
                           0.1.0
```

3.  **`fix/` (or `bugfix/`, `hotfix/`)**
    *   **Purpose:** Fix bugs in the existing codebase.
    *   **Workflow:** Create from `main` -> Develop fix -> Open PR to `main`.
    *   **Naming:** `fix/<short-description>` (e.g., `fix/typo-on-homepage`)
    *   **Versioning:** Merging results in a **patch** version increase (e.g., 1.2.2 -> 1.2.3).

```ascii
                       0.0.1-rc2
                0.0.1-rc   │    
                   ╵       ╵    
fix        ● ───── ● ───── ●    
          ╱                 ╲   
main   ─ ● ───────────────── ● ─
                             ╷  
                           0.0.1
```

4.  **`release/`**
    *   **Purpose:** Prepare for a **major** version release, often involving breaking changes or significant updates.
    *   **Workflow:** Create from `main` -> Develop changes -> Open PR to `main`.
    *   **Naming:** `release/<version-or-description>` (e.g., `release/2.0.0`, `release/api-overhaul`)
    *   **Versioning:** Merging results in a **major** version increase (e.g., 1.2.3 -> 2.0.0).

```ascii
                       1.0.0-rc2
                1.0.0-rc   │    
                   ╵       ╵    
release    ● ───── ● ───── ●    
          ╱                 ╲   
main   ─ ● ───────────────── ● ─
                             ╷  
                           1.0.0
```

## Basic Workflow

1.  **Create Branch:** From the latest `main`, create a branch following the naming conventions (`feature/`, `fix/`, `release/`).
    ```bash
    git checkout main
    git pull origin main
    git checkout -b feature/my-new-feature
    ```
2.  **Develop:** Make your code changes and commit them to your branch.
3.  **Push Branch:** Push your branch to the remote repository.
    ```bash
    git push -u origin feature/my-new-feature
    ```
4.  **Open Pull Request (PR):** Create a Pull Request on GitHub, targeting the `main` branch. Triggering CI/CD pipelines and generating the corresponding tags as release candidates.
5.  **Review & Merge:** After code review and checks pass, the PR is merged into `main`.
6.  **Release (Automated/Manual):** Merging into `main` typically triggers the release process (tagging, release notes based on the merged branch type).

This model ensures `main` always reflects a stable state, while development happens in isolated, purpose-driven branches integrated through reviewed Pull Requests.