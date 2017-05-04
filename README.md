# recompile

This invokes a command individually on a series of files. The command should
write the new output to stdout. This gets rewritten to disk according to the
`--extension` and `--out-dir` arguments.

### Why

Babel has a `--out-dir` flag, but it doesn't execute commands concurrently, and
it runs on every file in a directory. In our case we want to run concurrently,
and we only want to run babel on changed files; this isn't possible with Babel.

### Example Usage

Say you want to compile `.jsx` targets in the "components" directory to `.js`
files in the "public" directory:

```
JSX_TARGETS := $(shell find ./components -name '*.jsx' -depth 1)
JS_TARGETS := $(shell find ./public -name '*.js' -depth 1)

$(JS_TARGETS): $(JSX_TARGETS)
	recompile --command='./node_modules/.bin/babel --plugins transform-react-jsx' --extension=js --out-dir=public "$?"

recompile: $(JS_TARGETS)
```

That will only recompile the newer (changed) files.
