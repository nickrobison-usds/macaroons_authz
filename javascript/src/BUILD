package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_typescript//:defs.bzl", "ts_library")
load("@build_bazel_rules_nodejs//:defs.bzl", "nodejs_binary", "rollup_bundle")

ts_library(
  name = "endpoint",
  srcs = ["app.ts"],
  deps = [
       "@npm//@types",
       "@npm//express",
       "@npm//cookie-parser",
       "//javascript/src/controllers"
  ],
)

nodejs_binary(
  name = "server",
  node_modules = "//javascript:node_modules",
  entry_point = "macaroons_authz_demo/javascript/src/app.js",
  data = [
       ":endpoint",
    ]
)

rollup_bundle(
  name = "bundle",
  node_modules = "//javascript:node_modules",
  entry_point = "macaroons_authz_demo/javascript/src/app.js",
  deps = [
    ":endpoint",
    ]
)
