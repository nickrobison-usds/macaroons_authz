java_library(
    name = "jackson_jaxrs_base",
    exports = [
        "//external:jar/com/fasterxml/jackson/jaxrs/jackson_jaxrs_base"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_core",
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_databind"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jackson_jaxrs_json_provider",
    exports = [
        "//external:jar/com/fasterxml/jackson/jaxrs/jackson_jaxrs_json_provider"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/fasterxml/jackson/module:jackson_module_jaxb_annotations",
        ":jackson_jaxrs_base"
    ],
    visibility = [
        "//visibility:public"
    ]
)


