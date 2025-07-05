# git-scythe

A fast command-line tool built in Go that safely removes all local branches that have been merged into your repository's default branch (e.g., `main` or `master`). Keep your branch tree clean and organized by automating the deletion of stale branches.

## Description

`git-scythe` is designed to streamline your Git workflow by identifying and deleting all local branches that have already been merged into your default branch. It provides a final confirmation prompt before performing the bulk deletion.

This tool is particularly useful for development teams that create many feature branches and want to maintain a clean and manageable repository without manually pruning each branch.

## Usage

To run `git-scythe`, simply execute the command in your terminal:

```bash
git-scythe
```

The tool will display a list of all local branches that have been merged into your current `HEAD`. It will then prompt you to confirm whether you want to delete them.

### Options

- `--version`, `-v`: Print the version of the `git-scythe` tool.

## Examples

### Displaying the Version

To check the installed version of `git-scythe`, use the `--version` or `-v` flag.

```bash
$ git-scythe --version
git-scythe v0.1.0
```

### Deleting Merged Branches

When you run `git-scythe` without any flags, it will scan for local branches that have been merged into your default branch, list them, and ask for confirmation before deleting them.

```bash
$ git-scythe
The following merged branches will be deleted:
  > feature/add-login
  > fix/issue-123

Do you want to continue? (Y/n) y

Deleting 2 branches...
Deleted 2 branches.
```

### No Branches to Delete

If there are no merged branches to delete, the tool will inform you and exit gracefully.

```bash
$ git-scythe
No merged branches to display.
```

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any feedback or suggestions.
