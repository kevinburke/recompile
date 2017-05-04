# recompile

This invokes a command individually and concurrently on file arguments passed.
Each command invocation should write the new file output to stdout. That file
output gets written to disk according to the `--extension` and `--out-dir`
arguments.

Example:

```
recompile --command='babel --plugins=transform-react-jsx' --extension=js \
    --out-dir=public path/to/file1.jsx path/to/file2.jsx path/to/file3.jsx
```

That will run three instances of Babel concurrently and print the results to
public/file1.js, public/file2.js, and public/file3.js respectively.

### Why

Babel has a `--out-dir` flag, but it doesn't execute commands concurrently, and
it runs on every file in a directory. In our case we want to run concurrently,
and we only want to run babel on changed files; this isn't really possible with
Babel's current arguments.

You can run babel in `-w` mode, but that makes it trickier to integrate with
a larger build pipeline, for example if you want to restart the server each
time a file changes.

### Example Usage

Say you want to compile `.jsx` targets in the "components" directory to `.js`
files in the "public" directory:

```
JSX_TARGETS := $(shell find ./components -name '*.jsx' -depth 1)
JS_TARGETS := $(shell find ./public -name '*.js' -depth 1)
RECOMPILE := $(shell command -v recompile)

$(JS_TARGETS): $(JSX_TARGETS)
ifndef RECOMPILE
	go get -u github.com/kevinburke/recompile
endif
	recompile --command='./node_modules/.bin/babel --plugins transform-react-jsx' --extension=js --out-dir=public "$?"

recompile: $(JS_TARGETS)
```

That will only recompile the newer (changed) files.
