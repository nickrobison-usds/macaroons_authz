java_library(
    name = "jersey_client",
    exports = [
        "//external:jar/org/glassfish/jersey/core/jersey_client"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/ws/rs:javax_ws_rs_api",
        "//3rdparty/jvm/org/glassfish/hk2/external:javax_inject",
        "//3rdparty/jvm/org/glassfish/hk2:hk2_api",
        "//3rdparty/jvm/org/glassfish/hk2:hk2_locator",
        ":jersey_common"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jersey_common",
    exports = [
        "//external:jar/org/glassfish/jersey/core/jersey_common"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/annotation:javax_annotation_api",
        "//3rdparty/jvm/javax/ws/rs:javax_ws_rs_api",
        "//3rdparty/jvm/org/glassfish/hk2/external:javax_inject",
        "//3rdparty/jvm/org/glassfish/hk2:hk2_api",
        "//3rdparty/jvm/org/glassfish/hk2:hk2_locator",
        "//3rdparty/jvm/org/glassfish/hk2:osgi_resource_locator",
        "//3rdparty/jvm/org/glassfish/jersey/bundles/repackaged:jersey_guava"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jersey_server",
    exports = [
        "//external:jar/org/glassfish/jersey/core/jersey_server"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/annotation:javax_annotation_api",
        "//3rdparty/jvm/javax/validation:validation_api",
        "//3rdparty/jvm/javax/ws/rs:javax_ws_rs_api",
        "//3rdparty/jvm/org/glassfish/hk2/external:javax_inject",
        "//3rdparty/jvm/org/glassfish/hk2:hk2_api",
        "//3rdparty/jvm/org/glassfish/hk2:hk2_locator",
        "//3rdparty/jvm/org/glassfish/jersey/media:jersey_media_jaxb",
        ":jersey_client",
        ":jersey_common"
    ],
    visibility = [
        "//visibility:public"
    ]
)


