package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_typescript//:defs.bzl", "ts_library")

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
