# spectacle
An attempt at creating a way to generate API specifications Go packages. As most go packages are documented using `go doc` there doesn't seem to be any standard way of generating programmatically consumable API specifications. The end goal with this repo is to do just that, API specs in `yaml` or `json` for go packages.

## Run

```
  go run main.go
```
It will print all consts, types, variables and functions in the **decls** package which is the test set so far.
It will also create a debug log under log.

## How it's being build

Using the builtin compiler packages in go, such as `go/` - `ast`, `token` and `parser` I am trying to build something that would be a general solution. I've been using parts of the `go doc` implementation as examples but this project differs in some significant ways. We want to explicitly state the types of all variables for example.

## Remaining

What remains to be done.

  1. ~~Lookup assigned values for consts (and variables).~~
  2. Handle imports (for lookups) -- seems to be pretty hard.
  3. Handle concurrency errors, probably by first finding decls and then resolving.
  4. Filter unexported consts, variables, funcs and types.
  5. Sort all that is exported in alphabetical order (in groups).
  6. Group methods with their receivers.
  7. Specify API specification format (some architect please).
  8. Choose and implement yaml or json specification generation.

Lastly, the code is a bit messy and it needs to be cleaner and we need solid test cases.
Anything that can be handled by the go compiler should preferably be handled by this packages as well, but that might be a dream.
