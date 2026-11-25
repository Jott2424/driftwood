# Workflow

How changes get into `main`.

- `main` always holds working code.
- One short-lived branch per milestone, named like `m08-auth-register`.
- Each branch merges into `main` through a pull request, using **squash merge**,
  so every milestone is one commit in the history.
- Commit messages start with the milestone, like `M8: register endpoint`.
- Once CI exists (milestone 5), turn on branch protection so a pull request
  cannot merge until the checks pass.
- Tag a service when it is deployed to a VM, like `auth-v0.1.0`, so it is clear
  what is running where.
- The planning changes made before milestone 1 go straight to `main`. Everything
  from milestone 1 onward uses a branch.
