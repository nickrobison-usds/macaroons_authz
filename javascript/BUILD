package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_typescript//:defs.bzl", "ts_devserver")

# We export the tsconfig file, so that we can alias it in the package root. This means we can use the ts_library rules without manually specifying the tsconfig dependency.
exports_files([
  "tsconfig.json",
])

ts_devserver(
  name = "devserver",
  port = 3002,
  deps = ["//javascript/src:endpoint", "@npm//requirejs"],
  serving_path = "/bundle.js",
)

filegroup(
    name = "node_modules",
    srcs = glob([
        "node_modules/**",
        "node_modules/react-dom/**",
    ]),
)