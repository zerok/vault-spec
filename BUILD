load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    prefix = "github.com/zerok/vault-spec",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/zerok/vault-spec",
    visibility = ["//visibility:private"],
    deps = ["//cmd:go_default_library"],
)

go_binary(
    name = "vault-spec",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
