load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "platforms",
    sha256 = "5eda539c841265031c2f82d8ae7a3a6490bd62176e0c038fc469eabf91f6149b",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/platforms/releases/download/0.0.9/platforms-0.0.9.tar.gz",
        "https://github.com/bazelbuild/platforms/releases/download/0.0.9/platforms-0.0.9.tar.gz",
    ],
)

load("@platforms//host:extension.bzl", "host_platform_repo")

host_platform_repo(name = "host_platform")

http_archive(
    name = "bazel_skylib",
    sha256 = "66ffd9315665bfaafc96b52278f57c7e2dd09f5ede279ea6d39b2be471e7e3aa",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.4.2/bazel-skylib-1.4.2.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.4.2/bazel-skylib-1.4.2.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

http_archive(
    name = "aspect_bazel_lib",
    sha256 = "6d758a8f646ecee7a3e294fbe4386daafbe0e5966723009c290d493f227c390b",
    strip_prefix = "bazel-lib-2.7.7",
    url = "https://github.com/aspect-build/bazel-lib/releases/download/v2.7.7/bazel-lib-v2.7.7.tar.gz",
)

http_archive(
    name = "aspect_rules_js",
    sha256 = "3dfccf2713288e0518c0485b65574ca66426c6e06495299abe6f5c64e3bc6314",
    strip_prefix = "rules_js-2.0.0-rc3",
    url = "https://github.com/aspect-build/rules_js/releases/download/v2.0.0-rc3/rules_js-v2.0.0-rc3.tar.gz",
)

http_archive(
    name = "aspect_rules_ts",
    sha256 = "3ea5cdb825d5dbffe286b3d9c5197a2648cf04b5e6bd8b913a45823cdf0ae960",
    strip_prefix = "rules_ts-3.0.0-rc0",
    url = "https://github.com/aspect-build/rules_ts/releases/download/v3.0.0-rc0/rules_ts-v3.0.0-rc0.tar.gz",
)

http_archive(
    name = "aspect_rules_swc",
    sha256 = "c085647585c3d01bee3966eb9ba433a1efbb0ee79bb1b8c67882a81d82a9b37f",
    strip_prefix = "rules_swc-2.0.0-rc0",
    url = "https://github.com/aspect-build/rules_swc/releases/download/v2.0.0-rc0/rules_swc-v2.0.0-rc0.tar.gz",
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "af47f30e9cbd70ae34e49866e201b3f77069abb111183f2c0297e7e74ba6bbc0",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.47.0/rules_go-v0.47.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.47.0/rules_go-v0.47.0.zip",
    ],
)

http_archive(
    name = "rules_proto",
    sha256 = "dc3fb206a2cb3441b485eb1e423165b231235a1ea9b031b4433cf7bc1fa460dd",
    strip_prefix = "rules_proto-5.3.0-21.7",
    urls = [
        "https://github.com/bazelbuild/rules_proto/archive/refs/tags/5.3.0-21.7.tar.gz",
    ],
)

http_archive(
    name = "rules_proto_grpc",
    sha256 = "9ba7299c5eb6ec45b6b9a0ceb9916d0ab96789ac8218269322f0124c0c0d24e2",
    strip_prefix = "rules_proto_grpc-4.5.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/releases/download/4.5.0/rules_proto_grpc-4.5.0.tar.gz"],
)

http_archive(
    name = "rules_buf",
    sha256 = "bc2488ee497c3fbf2efee19ce21dceed89310a08b5a9366cc133dd0eb2118498",
    strip_prefix = "rules_buf-0.2.0",
    urls = [
        "https://github.com/bufbuild/rules_buf/archive/refs/tags/v0.2.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    integrity = "sha256-MpOL2hbmcABjA1R5Bj2dJMYO2o15/Uc5Vj9Q0zHLMgk=",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
    ],
)

http_archive(
    name = "rules_rust",
    integrity = "sha256-ZQGWDD5NoySV0eEAfe0HaaU0yxlcMN6jaqVPnYo/A2E=",
    urls = ["https://github.com/bazelbuild/rules_rust/releases/download/0.38.0/rules_rust-v0.38.0.tar.gz"],
)

# Container rules
http_archive(
    name = "rules_oci",
    sha256 = "647f4c6fd092dc7a86a7f79892d4b1b7f1de288bdb4829ca38f74fd430fcd2fe",
    strip_prefix = "rules_oci-1.7.6",
    url = "https://github.com/bazel-contrib/rules_oci/releases/download/v1.7.6/rules_oci-v1.7.6.tar.gz",
)

http_archive(
    name = "rules_pkg",
    sha256 = "8c20f74bca25d2d442b327ae26768c02cf3c99e93fad0381f32be9aab1967675",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_pkg/releases/download/0.8.1/rules_pkg-0.8.1.tar.gz",
        "https://github.com/bazelbuild/rules_pkg/releases/download/0.8.1/rules_pkg-0.8.1.tar.gz",
    ],
)

http_archive(
    name = "container_structure_test",
    sha256 = "42edb647b51710cb917b5850380cc18a6c925ad195986f16e3b716887267a2d7",
    strip_prefix = "container-structure-test-104a53ede5f78fff72172639781ac52df9f5b18f",
    urls = ["https://github.com/GoogleContainerTools/container-structure-test/archive/104a53ede5f78fff72172639781ac52df9f5b18f.zip"],
)

http_archive(
    name = "buildifier_prebuilt",
    sha256 = "8ada9d88e51ebf5a1fdff37d75ed41d51f5e677cdbeafb0a22dda54747d6e07e",
    strip_prefix = "buildifier-prebuilt-6.4.0",
    urls = ["https://github.com/keith/buildifier-prebuilt/archive/6.4.0.tar.gz"],
)

http_archive(
    name = "aspect_cli",
    repo_mapping = {
        "@com_github_smacker_go_tree_sitter": "@aspectcli-com_github_smacker_go_tree_sitter",
    },
    sha256 = "045f0186edb25706dfe77d9c4916eec630a2b2736f9abb59e37eaac122d4b771",
    strip_prefix = "aspect-cli-5.8.20",
    url = "https://github.com/aspect-build/aspect-cli/archive/5.8.20.tar.gz",
)

load("@aspect_bazel_lib//lib:repositories.bzl", "aspect_bazel_lib_dependencies", "aspect_bazel_lib_register_toolchains")

aspect_bazel_lib_dependencies()

aspect_bazel_lib_register_toolchains()

http_archive(
    name = "rules_apko",
    patch_args = ["-p1"],
    patches = [
        # required due to https://github.com/chainguard-dev/apko/issues/1052
        "//third_party/rules_apko:repository_label_strip.patch",
        # required until a release contains https://github.com/chainguard-dev/rules_apko/pull/53
        "//third_party/rules_apko:apko_run_runfiles_path.patch",
        # symlinking the lockfile appears to be problematic in CI https://github.com/sourcegraph/sourcegraph/pull/61877
        "//third_party/rules_apko:copy_dont_symlink_lockfile.patch",
    ],
    sha256 = "f176171f95ee2b6eef1572c6da796d627940a1e898a32d476a2d7a9a99332960",
    strip_prefix = "rules_apko-1.2.2",
    url = "https://github.com/chainguard-dev/rules_apko/releases/download/v1.2.2/rules_apko-v1.2.2.tar.gz",
)

# hermetic_cc_toolchain setup ================================
HERMETIC_CC_TOOLCHAIN_VERSION = "v2.2.1"

# Please note that we only use hermetic-cc for local development purpose and Nix, at it eases the path to cross-compile
# so we can produce container images locally on Mac laptops.
#
# @jhchabran See https://github.com/sourcegraph/sourcegraph/pull/55969, there is an ongoing issue with UBSAN
# and treesitter, that breaks the compilation of syntax-highlighter. Since we only use
# hermetic_cc for local development purposes, while it's a bit heavy handed for a --copt, it's acceptable
# at this point. Passing --copt=-fno-sanitize=undefined sadly doesn't fix the problem, which is why
# we have to patch to inject the flag.
http_archive(
    name = "hermetic_cc_toolchain",
    patch_args = ["-p1"],
    patches = [
        "//third_party/hermetic_cc:disable_ubsan.patch",
    ],
    sha256 = "3b8107de0d017fe32e6434086a9568f97c60a111b49dc34fc7001e139c30fdea",
    urls = [
        "https://mirror.bazel.build/github.com/uber/hermetic_cc_toolchain/releases/download/{0}/hermetic_cc_toolchain-{0}.tar.gz".format(HERMETIC_CC_TOOLCHAIN_VERSION),
        "https://github.com/uber/hermetic_cc_toolchain/releases/download/{0}/hermetic_cc_toolchain-{0}.tar.gz".format(HERMETIC_CC_TOOLCHAIN_VERSION),
    ],
)

# rules_js setup ================================
load("@aspect_rules_js//js:repositories.bzl", "rules_js_dependencies")

rules_js_dependencies()

load("@aspect_rules_js//js:toolchains.bzl", "rules_js_register_toolchains")

rules_js_register_toolchains(node_version = "20.8.0")

# rules_js npm setup ============================
load("@aspect_rules_js//npm:repositories.bzl", "npm_translate_lock")

npm_translate_lock(
    name = "npm",
    npm_package_target_name = "{dirname}_pkg",
    npmrc = "//:.npmrc",
    pnpm_lock = "//:pnpm-lock.yaml",
    # Required for ESLint test targets.
    # See https://github.com/aspect-build/rules_js/issues/239
    # See `public-hoist-pattern[]=*eslint*` in the `.npmrc` of this monorepo.
    public_hoist_packages = {
        "@typescript-eslint/eslint-plugin": [""],
        "@typescript-eslint/parser@5.56.0_qxbo2xm47qt6fxnlmgbosp4hva": [""],
        "eslint-config-prettier": [""],
        "eslint-plugin-ban": [""],
        "eslint-plugin-etc": [""],
        "eslint-plugin-import": [""],
        "eslint-plugin-jest-dom": [""],
        "eslint-plugin-jsdoc": [""],
        "eslint-plugin-jsx-a11y": [""],
        "eslint-plugin-react@7.32.1_eslint_8.34.0": [""],
        "eslint-plugin-react-hooks": [""],
        "eslint-plugin-rxjs": [""],
        "eslint-plugin-unicorn": [""],
        "eslint-plugin-unused-imports": [""],
        "eslint-import-resolver-node": [""],
    },
    verify_node_modules_ignored = "//:.bazelignore",
)

# rules_ts npm setup ============================
load("@npm//:repositories.bzl", "npm_repositories")

npm_repositories()

load("@aspect_rules_ts//ts:repositories.bzl", "rules_ts_dependencies")

rules_ts_dependencies(ts_version = "4.9.5")

# rules_swc setup ==============================
load("@aspect_rules_swc//swc:dependencies.bzl", "rules_swc_dependencies")

rules_swc_dependencies()

load("@aspect_rules_swc//swc:repositories.bzl", "LATEST_SWC_VERSION", "swc_register_toolchains")

swc_register_toolchains(
    name = "swc",
    swc_version = LATEST_SWC_VERSION,
)

# rules_esbuild setup ===========================
http_archive(
    name = "aspect_rules_esbuild",
    patch_args = ["-p1"],
    patches = [
        # Includes https://github.com/aspect-build/rules_esbuild/pull/201 as well as a fix for
        # object-inspect being weird, see the comments in the patch for further links.
        "//third_party/rules_esbuild:sandbox-plugin-fixes.patch",
    ],
    sha256 = "ef7163a2e8e319f8a9a70560788dd899126aebf3538c76f8bc1f0b4b52ba4b56",
    strip_prefix = "rules_esbuild-0.21.0-rc1",
    url = "https://github.com/aspect-build/rules_esbuild/releases/download/v0.21.0-rc1/rules_esbuild-v0.21.0-rc1.tar.gz",
)

load("@aspect_rules_esbuild//esbuild:dependencies.bzl", "rules_esbuild_dependencies")

rules_esbuild_dependencies()

# Register a toolchain containing esbuild npm package and native bindings
load("@aspect_rules_esbuild//esbuild:repositories.bzl", "esbuild_register_toolchains")

esbuild_register_toolchains(
    name = "esbuild",
    # Note, this differs from the version noted in package.json, however we've been inadvertently building with this version for some time now so we'll stick with it and revisit.
    esbuild_version = "0.19.2",
)

# Go toolchain setup

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("//:deps.bzl", "go_dependencies")
load("//:linter_deps.bzl", "linter_dependencies")

go_repository(
    name = "com_github_aws_aws_sdk_go_v2_service_ssooidc",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/aws/aws-sdk-go-v2/service/ssooidc",
    sum = "h1:xLPZMyuZ4GuqRCIec/zWuIhRFPXh2UOJdLXBSi64ZWQ=",
    version = "v1.14.5",
)

# Overrides the default provided protobuf dep from rules_go by a more
# recent one.
go_repository(
    name = "org_golang_google_protobuf",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/protobuf",
    sum = "h1:uNO2rsAINq/JlFpSdYEKIZ0uKD/R9cpdv0T+yoGwGmI=",
    version = "v1.33.0",
)

# Pin protoc-gen-go-grpc to 1.3.0
# See also //:gen-go-grpc
go_repository(
    name = "org_golang_google_grpc_cmd_protoc_gen_go_grpc",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
    sum = "h1:rNBFJjBCOgVr9pWD7rs/knKL4FRTKgpZmsRfV214zcA=",
    version = "v1.3.0",
)  # keep

# Pin specific version for aspect-cli's gazelle rules, with versions
# that it requires but that our codebase doesnt support.
go_repository(
    name = "aspectcli-com_github_smacker_go_tree_sitter",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/smacker/go-tree-sitter",
    sum = "h1:DxgjlvWYsb80WEN2Zv3WqJFAg2DKjUQJO6URGdf1x6Y=",
    version = "v0.0.0-20230720070738-0d0a9f78d8f8",
)  # keep

load("@aspect_cli//:go.bzl", aspect_cli_deps = "deps")

aspect_cli_deps()

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains(
    nogo = "@//:sg_nogo",
    version = "1.22.1",
)

linter_dependencies()

gazelle_dependencies()

# rust toolchain setup
load("@rules_rust//rust:repositories.bzl", "rules_rust_dependencies", "rust_register_toolchains", "rust_repository_set")

rules_rust_dependencies()

rust_version = "1.78.0"

rust_register_toolchains(
    edition = "2021",
    # Keep in sync with docker-images/syntax-highlighter/Dockerfile
    # and docker-images/syntax-highlighter/rust-toolchain.toml
    versions = [
        rust_version,
    ],
)

# Needed for locally cross-compiling rust binaries to linux/amd64 on a Mac laptop, when seeking to
# create container images in local for testing purposes.
rust_repository_set(
    name = "macos_arm_64",
    edition = "2021",
    exec_triple = "aarch64-apple-darwin",
    extra_target_triples = ["x86_64-unknown-linux-gnu"],
    versions = [rust_version],
)

load("@rules_rust//crate_universe:defs.bzl", "crates_repository")

crates_repository(
    name = "crate_index",
    cargo_config = "//docker-images/syntax-highlighter:.cargo/config.toml",
    cargo_lockfile = "//docker-images/syntax-highlighter:Cargo.lock",
    # this file has to be manually created and it will be filled when
    # the target is ran.
    # To regenerate this file run: CARGO_BAZEL_REPIN=1 bazel sync --only=crate_index
    lockfile = "//docker-images/syntax-highlighter:Cargo.Bazel.lock",
    # glob doesn't work in WORKSPACE files: https://github.com/bazelbuild/bazel/issues/11935
    manifests = [
        "//docker-images/syntax-highlighter:Cargo.toml",
        "//docker-images/syntax-highlighter:crates/syntax-analysis/Cargo.toml",
        "//docker-images/syntax-highlighter:crates/tree-sitter-all-languages/Cargo.toml",
        "//docker-images/syntax-highlighter:crates/scip-syntax/Cargo.toml",
    ],
)

load("@crate_index//:defs.bzl", "crate_repositories")

crate_repositories()

load("@hermetic_cc_toolchain//toolchain:defs.bzl", zig_toolchains = "toolchains")

zig_toolchains()

# containers steup       ===============================
load("@rules_oci//oci:dependencies.bzl", "rules_oci_dependencies")

rules_oci_dependencies()

load("@rules_oci//oci:repositories.bzl", "LATEST_CRANE_VERSION", "oci_register_toolchains")

oci_register_toolchains(
    name = "oci",
    crane_version = LATEST_CRANE_VERSION,
    # Uncommenting the zot toolchain will cause it to be used instead of crane for some tasks.
    # Note that it does not support docker-format images.
    # zot_version = LATEST_ZOT_VERSION,
)

# Optional, for oci_tarball rule
load("@rules_pkg//:deps.bzl", "rules_pkg_dependencies")

rules_pkg_dependencies()

load("//dev:oci_deps.bzl", "oci_deps")

oci_deps()

load("@container_structure_test//:repositories.bzl", "container_structure_test_register_toolchain")

container_structure_test_register_toolchain(name = "cst")

load("//dev:tool_deps.bzl", "tool_deps")

tool_deps()

# Buildifier
load("@buildifier_prebuilt//:deps.bzl", "buildifier_prebuilt_deps")

buildifier_prebuilt_deps()

load("@buildifier_prebuilt//:defs.bzl", "buildifier_prebuilt_register_toolchains")

buildifier_prebuilt_register_toolchains()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")
load("@rules_proto_grpc//doc:repositories.bzl", rules_proto_grpc_doc_repos = "doc_repos")
load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")

rules_proto_grpc_toolchains()

rules_proto_grpc_repos()

rules_proto_grpc_go_repos()

rules_proto_grpc_doc_repos()

load("@rules_buf//buf:repositories.bzl", "rules_buf_dependencies", "rules_buf_toolchains")

rules_buf_dependencies()

rules_buf_toolchains(
    sha256 = "f227f04f3f910a7611e8841d50172e3c0e9a94ad21760e6f8abbe3666d682ab5",
    version = "v1.31.0",
)

load("@rules_buf//gazelle/buf:repositories.bzl", "gazelle_buf_dependencies")

gazelle_buf_dependencies()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# keep revision up-to-date with client/browser/scripts/build-inline-extensions.js
http_archive(
    name = "sourcegraph_extensions_bundle",
    add_prefix = "bundle",
    build_file_content = """
package(default_visibility = ["//visibility:public"])

exports_files(["bundle"])

filegroup(
    name = "srcs",
    srcs = glob(["**"]),
)
    """,
    integrity = "sha256-Spx8LyM7k+dsGOlZ4TdAq+CNk5EzvYB/oxnY4zGpqPg=",
    strip_prefix = "sourcegraph-extensions-bundles-5.0.1",
    url = "https://github.com/sourcegraph/sourcegraph-extensions-bundles/archive/v5.0.1.zip",
)

load("//dev:schema_migrations.bzl", "schema_migrations")

schema_migrations(
    name = "schemas_migrations",
    updated_at = "2024-05-07 14:39",
)

# wolfi images setup ================================

load("@rules_apko//apko:repositories.bzl", "apko_register_toolchains", "rules_apko_dependencies")

rules_apko_dependencies()

# We don't register the default toolchains, and regsiter our own from a patched go_repository sourced
# go_binary target that contains some fixes that are not yet merged upstream.
# https://github.com/chainguard-dev/go-apk/pull/216
apko_register_toolchains(
    name = "apko",
    register = False,
)

register_toolchains("//:apko_linux_toolchain")

register_toolchains("//:apko_darwin_arm64_toolchain")

register_toolchains("//:apko_darwin_amd64_toolchain")

load("//wolfi-images:repo.bzl", "wolfi_lockfiles")

wolfi_lockfiles(name = "apko_lockfiles")

load("@apko_lockfiles//:translates.bzl", "apko_translate_locks")

apko_translate_locks()

load("@apko_lockfiles//:repositories.bzl", "apko_repositories")

apko_repositories()
