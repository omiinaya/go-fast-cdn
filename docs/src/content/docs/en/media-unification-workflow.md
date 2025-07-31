# Media Unification Development Workflow

This document provides guidelines for the development workflow during the media unification project. It covers how to set up and work with the development branch, commit changes, sync with the main branch, and eventually merge the changes back.

## Table of Contents

1. [Setting Up the Development Branch](#setting-up-the-development-branch)
2. [Working with the Development Branch](#working-with-the-development-branch)
3. [Committing Changes](#committing-changes)
4. [Syncing with the Main Branch](#syncing-with-the-main-branch)
5. [Merging Changes Back to Main](#merging-changes-back-to-main)
6. [Troubleshooting](#troubleshooting)

## Setting Up the Development Branch

### Using the Setup Script

The easiest way to set up the development branch is to use the provided setup script:

#### For Linux/macOS:

```bash
./scripts/setup-media-unification-branch.sh
```

#### For Windows:

```cmd
scripts\setup-media-unification-branch.bat
```

The script will:
- Create a new branch named `feature/media-unification` from the main branch
- Switch to the new branch
- Optionally push the branch to the remote repository

### Manual Setup

If you prefer to set up the branch manually:

1. Ensure you're on the main branch:
   ```bash
   git checkout main
   ```

2. Pull the latest changes:
   ```bash
   git pull origin main
   ```

3. Create and switch to a new branch:
   ```bash
   git checkout -b feature/media-unification
   ```

4. Push the new branch to the remote repository:
   ```bash
   git push -u origin feature/media-unification
   ```

### Script Options

The setup script accepts the following options:

- `-h, --help`: Show help message
- `-b, --branch NAME`: Specify a custom branch name (default: feature/media-unification)
- `-s, --source NAME`: Specify source branch (default: main)

Examples:
```bash
# Create with default settings
./scripts/setup-media-unification-branch.sh

# Create with custom branch name
./scripts/setup-media-unification-branch.sh -b my-media-feature

# Create from a different source branch
./scripts/setup-media-unification-branch.sh -b my-media-feature -s develop
```

## Working with the Development Branch

### Switching Between Branches

To switch to the development branch:
```bash
git checkout feature/media-unification
```

To switch back to the main branch:
```bash
git checkout main
```

To see all available branches:
```bash
git branch -a
```

### Checking Your Current Branch

To see which branch you're currently on:
```bash
git branch
```

The current branch will be marked with an asterisk (*).

## Committing Changes

### Commit Guidelines

During the media unification work, follow these guidelines for commits:

1. **Atomic Commits**: Each commit should represent a single logical change. Don't mix unrelated changes in the same commit.

2. **Clear Commit Messages**: Write descriptive commit messages that explain what was changed and why.

   Format:
   ```
   Type(scope): Short description

   Detailed explanation of the change, if necessary.
   ```

   Types:
   - `feat`: A new feature
   - `fix`: A bug fix
   - `docs`: Documentation changes
   - `style`: Code style changes (formatting, etc.)
   - `refactor`: Code refactoring
   - `test`: Adding or modifying tests
   - `chore`: Maintenance tasks

   Example:
   ```
   feat(media): implement unified media repository

   Add a new unified media repository that handles both images and documents.
   This replaces the separate image and document repositories with a single
   interface that can handle multiple media types.
   ```

3. **Commit Frequently**: Commit small, incremental changes rather than large, infrequent ones.

4. **Test Before Committing**: Ensure all tests pass before committing your changes.

### Commit Workflow

1. Stage your changes:
   ```bash
   git add .
   ```
   Or stage specific files:
   ```bash
   git add path/to/file1 path/to/file2
   ```

2. Commit your changes with a descriptive message:
   ```bash
   git commit -m "feat(media): add unified media upload functionality"
   ```

3. Push your changes to the remote repository:
   ```bash
   git push origin feature/media-unification
   ```

### Viewing Commit History

To see the commit history for the current branch:
```bash
git log --oneline
```

To see a more detailed history:
```bash
git log
```

To see the differences between your branch and the main branch:
```bash
git log main..feature/media-unification --oneline
```

## Syncing with the Main Branch

### Why Sync?

It's important to periodically sync your development branch with the main branch to:
- Incorporate the latest changes and bug fixes from the main branch
- Reduce the risk of merge conflicts when it's time to merge your changes back
- Ensure your changes are compatible with the latest codebase

### Sync Workflow

1. Switch to the main branch:
   ```bash
   git checkout main
   ```

2. Pull the latest changes:
   ```bash
   git pull origin main
   ```

3. Switch back to your development branch:
   ```bash
   git checkout feature/media-unification
   ```

4. Merge the main branch into your development branch:
   ```bash
   git merge main
   ```

5. Resolve any merge conflicts if they occur (see [Troubleshooting](#troubleshooting))

6. Push the updated branch to the remote repository:
   ```bash
   git push origin feature/media-unification
   ```

### Alternative: Rebase

An alternative to merging is to rebase your development branch on top of the main branch:

1. Switch to your development branch:
   ```bash
   git checkout feature/media-unification
   ```

2. Fetch the latest changes:
   ```bash
   git fetch origin
   ```

3. Rebase your branch on top of the main branch:
   ```bash
   git rebase origin/main
   ```

4. Resolve any conflicts if they occur (see [Troubleshooting](#troubleshooting))

5. Force push your rebased branch to the remote repository:
   ```bash
   git push --force-with-lease origin feature/media-unification
   ```

**Note**: Use rebase with caution, especially if you're collaborating with others on the same branch. Rebase rewrites the history of your branch, which can cause issues for other developers who have based their work on your branch.

## Merging Changes Back to Main

### Pre-Merge Checklist

Before merging your changes back to the main branch, ensure that:

1. All features are complete and working as expected
2. All tests pass
3. The code has been reviewed (if applicable)
4. Documentation is updated
5. Your branch is up to date with the latest main branch

### Merge Workflow

1. Ensure your branch is up to date with the main branch (see [Syncing with the Main Branch](#syncing-with-the-main-branch))

2. Switch to the main branch:
   ```bash
   git checkout main
   ```

3. Pull the latest changes:
   ```bash
   git pull origin main
   ```

4. Merge your development branch into main:
   ```bash
   git merge feature/media-unification
   ```

5. Resolve any merge conflicts if they occur (see [Troubleshooting](#troubleshooting))

6. Push the updated main branch to the remote repository:
   ```bash
   git push origin main
   ```

### Alternative: Pull Request

If your team uses pull requests (or merge requests), follow these steps:

1. Push your final changes to your development branch:
   ```bash
   git push origin feature/media-unification
   ```

2. Create a pull request from your development branch to the main branch through your Git hosting platform (GitHub, GitLab, etc.)

3. Address any review comments and make necessary changes

4. Once the pull request is approved, merge it through the platform

### Post-Merge Cleanup

After your changes have been successfully merged into the main branch:

1. You can delete the development branch locally:
   ```bash
   git branch -d feature/media-unification
   ```

2. Delete the development branch from the remote repository:
   ```bash
   git push origin --delete feature/media-unification
   ```

## Troubleshooting

### Merge Conflicts

Merge conflicts occur when Git can't automatically reconcile differences between branches. To resolve them:

1. When a conflict occurs during a merge or rebase, Git will mark the conflicted files in your working directory.

2. Open the conflicted files and look for the conflict markers:
   ```
   <<<<<<< HEAD
   Changes from the current branch (main)
   =======
   Changes from the other branch (feature/media-unification)
   >>>>>>> feature/media-unification
   ```

3. Edit the files to resolve the conflicts by:
   - Choosing which version to keep
   - Combining both versions
   - Writing new code that incorporates both sets of changes

4. Remove the conflict markers.

5. Stage the resolved files:
   ```bash
   git add path/to/resolved/file
   ```

6. Continue the merge or rebase:
   - For merge: `git commit`
   - For rebase: `git rebase --continue`

7. If you need to abort the merge or rebase:
   - For merge: `git merge --abort`
   - For rebase: `git rebase --abort`

### Accidentally Committing to the Wrong Branch

If you accidentally commit changes to the wrong branch:

1. Reset the branch to remove the unwanted commits:
   ```bash
   git reset --hard HEAD~1  # Remove the last commit
   ```

2. Switch to the correct branch:
   ```bash
   git checkout feature/media-unification
   ```

3. Cherry-pick the commit you removed:
   ```bash
   git cherry-pick main  # Replace with the commit hash if needed
   ```

### Lost Commits

If you think you've lost commits:

1. Check the reflog to see a history of all actions:
   ```bash
   git reflog
   ```

2. Identify the commit you want to restore and reset to it:
   ```bash
   git reset --hard commit_hash
   ```

### Push Rejected

If your push is rejected because the remote branch has diverged:

1. Pull the latest changes with rebase:
   ```bash
   git pull --rebase origin feature/media-unification
   ```

2. Resolve any conflicts if they occur

3. Push your changes again:
   ```bash
   git push origin feature/media-unification
   ```

## Conclusion

This workflow provides a structured approach to developing the media unification feature while maintaining code quality and minimizing conflicts. By following these guidelines, you can ensure a smooth development process and a successful merge of your changes into the main branch.