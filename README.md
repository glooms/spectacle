# spectacle
An attempt at creating a way to generate API specifications Go packages. As most go packages are documented using `go doc` there doesn't seem to be any standard way of generating programmatically consumable API specifications. The end goal with this repo is to do just that, API specs in `yaml` or `json` for go packages.

## Run

```
  go run main.go
  go run main.go <go-package-path>
```
Running without arguments will be like running it on the directory `.`.
It will print any packages found and log the specifications to a **log** directory and
any debug information to a **debug** directory. Normally these directories are created in the
root of the git repo, i.e. `./log` but it's relative to where you are running it from.

Added so that the absolute path to the log and debug directories are printed when the program
is run.

## Remaining

What remains to be done.

  1. ~~Lookup assigned values for consts (and variables).~~
  2. ~~Handle imports (for lookups) -- seems to be pretty hard.~~
  3. ~~Handle concurrency errors, probably by first finding decls and then resolving.~~
  4. ~~Filter unexported consts, variables, funcs and types.~~
  5. ~~Sort all that is exported in alphabetical order (in groups).~~
  6. Group methods with their receivers.
  7. Specify API specification format (some architect please).
  8. Choose and implement yaml or json specification generation.
  9. Include parsed comments?

After some research I found the **go/types** package which provides almost all utilities need to create
a general and robust implementation of a golang specificatios.

The only drawback is that I've not yet found how to include the parsed documentation from the source files.

What remains now is to get some clarity as to what the specification should look like. Research concluded.
