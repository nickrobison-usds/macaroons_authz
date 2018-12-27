java_binary(
    name = "JavaClient",
    srcs = glob(["src/main/java/**/*.java"]),
    main_class = "com.nickrobison.cmsauthz.JavaClientApplication",
    deps = [
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_databind",
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_core",
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_annotations",
        "//3rdparty/jvm/javax/ws/rs:javax_ws_rs_api",
        "//3rdparty/jvm/javax/validation:validation_api",
        "//3rdparty/jvm/io/dropwizard:dropwizard_core",
        "//3rdparty/jvm/io/dropwizard:dropwizard_client",
        "//3rdparty/jvm/io/dropwizard:dropwizard_jersey",
        "//3rdparty/jvm/io/dropwizard:dropwizard_jackson",
        "//3rdparty/jvm/io/dropwizard:dropwizard_metrics",
        "//3rdparty/jvm/com/github/nitram509:jmacaroons",
        "//3rdparty/jvm/eu/neilalexander:jnacl"
    ]
)