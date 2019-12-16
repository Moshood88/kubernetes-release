# gazelle:repository_macro repos.bzl%go_repositories
workspace(name = "io_k8s_release")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

################################################################################
# Go Build Definitions
################################################################################

git_repository(
    name = "io_k8s_repo_infra",
    commit = "4f75b5b6e1958c8be9bcbf95a1f3a4010d8687c0",
    remote = "https://github.com/kubernetes/repo-infra.git",
    shallow_since = "1576262829 -0800",
)

load("@io_k8s_repo_infra//:load.bzl", "repositories")

repositories()

load("@io_k8s_repo_infra//:repos.bzl", "configure")

configure()

load("//:repos.bzl", "go_repositories")

go_repositories()

http_file(
    name = "jq",
    downloaded_file_path = "jq",
    executable = True,
    sha256 = "af986793a515d500ab2d35f8d2aecd656e764504b789b66d7e1a0b727a124c44",
    urls = ["https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64"],
)
