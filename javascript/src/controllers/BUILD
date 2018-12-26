package(default_visibility = ["//visibility:public"])

load("@build_bazel_rules_typescript//:defs.bzl", "ts_library")

ts_library(
  name = "controllers",
  srcs = ["auth.ts", "token.ts"],
  deps = [
    "@npm//@types",
    "@npm//express",
    "@npm//macaroon",
    "@npm//pg",
    "@npm//retryer",
    "@npm//varint",
    "//javascript/typings",
  ],
)