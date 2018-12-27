workspace(name = "macaroons_authz_demo")

#load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# Get the Node JS rules

git_repository(
    name = "build_bazel_rules_nodejs",
    remote = "https://github.com/bazelbuild/rules_nodejs.git",
    tag = "0.16.3",
)

git_repository(
    name = "build_bazel_rules_typescript",
    remote = "https://github.com/bazelbuild/rules_typescript.git",
    tag = "0.22.0",
)

load("@build_bazel_rules_typescript//:package.bzl", "rules_typescript_dependencies")

rules_typescript_dependencies()

load("@build_bazel_rules_nodejs//:package.bzl", "rules_nodejs_dependencies")

rules_nodejs_dependencies()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories", "yarn_install")

node_repositories(
    node_version = "10.13.0",
    yarn_version = "1.12.1",
)

yarn_install(
    name = "npm",
    package_json = "//javascript:package.json",
    yarn_lock = "//javascript:yarn.lock",
)


load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")


go_rules_dependencies()

go_register_toolchains()

load("//3rdparty:workspace.bzl", "maven_dependencies")
maven_dependencies()
