## Build, Lint, and Test

- **Build**: `go build`
- **Test**: `go test ./...`
- **Run a single test**: `go test -run ^TestMyFunction$`
- **Lint**: `golangci-lint run` (if installed) or `go vet ./...`

## Code Style

- **Formatting**: Use `gofmt` or `goimports` to format code.
- **Imports**: Group standard library, third-party, and internal packages separately.
- **Types**: Use structs for complex data. Keep interfaces small and focused.
- **Naming**: Use `camelCase` for variables and functions. Use `PascalCase` for exported symbols.
- **Error Handling**: Use `if err != nil` for error handling. Add context to errors when possible.
- **Comments**: Add comments to explain complex logic or public APIs.
- **Dependencies**: Use `go mod tidy` to manage dependencies.
- **Packages**: Keep packages focused on a single responsibility.
- **Concurrency**: Use goroutines and channels for concurrent operations.
- **Testing**: Write unit tests for all new features and bug fixes.

## Commit Messages

- **Format**: Use [Conventional Commits](https://www.conventionalcommits.org/).
- **Type**: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`.
- **Scope**: Optional, e.g., `feat(parser): add new feature`. If not scope it's recognized, ignore it.
- **Description**: 
    * Direct, 80 characters or less, 
    * Capitalize the first letter of the first word, the rest should be lowercase.
    * Always add a period (.) on the end of the description (IMPORTANT).
    * e.g `feat(app): Enhances list display and user interaction.`.
- **Body**: Explain the "what" and "why". No bullet points, no bullet lists, just plain text. Ignore the body if the changes are small or self-explanatory, eg. `Introduces i18n support for list items.`. Do not list the files to be added in the body. If present, it should be a paragraph 100 characters or less.
- **Footer**: Only add if breaking changes.

The commit message should be structured as follows:

```
<type>(optional scope): <description>

[optional body]

[optional footer(s)]
```

### Examples:

`refactor(utils): Moves dependency to correct package declaration.`

```
style: Updates layout colors.

The colors are now more vibrant and the layout is more modern.
```

```
feat(deps): Upgrades `foo` dependency.

BREAKING CHANGE: the bump of the dependency is not compatible with the previous version.
```

## Commiting Changes
- **Commit**: NEVER COMMIT WITHOUT PREVIEWING. Only commit when the prompt strictly is `commit changes`
- **Preview**: Preview generated commit message and added/staged/unstaged files before actual commiting. If the message is not correct, do not commit. The prompt `generate commit` or similars are not valid for commiting changes, just to enable the preview. Do not list the files to be added in the commit message. Do not add/stage files to the commit, only the `commit` command is valid for commiting changes and adding/staging files.
- **Files**: If we have staged/added files, generate the commit and proceed to evaluate all the committing flow based on those files. Only, and only if we don't have added/staged files, work with all the unstaged changes.

## Branching
- **Main**: Main branch, always green
- **Dev**: Development branch, always yellow
- **Feature**: Feature branch, always blue
- **Hotfix**: Hotfix branch, always red
