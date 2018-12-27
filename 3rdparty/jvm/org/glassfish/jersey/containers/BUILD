java_library(
    name = "jersey_container_servlet",
    exports = [
        "//external:jar/org/glassfish/jersey/containers/jersey_container_servlet"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/ws/rs:javax_ws_rs_api",
        "//3rdparty/jvm/org/glassfish/jersey/core:jersey_common",
        "//3rdparty/jvm/org/glassfish/jersey/core:jersey_server",
        ":jersey_container_servlet_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jersey_container_servlet_core",
    exports = [
        "//external:jar/org/glassfish/jersey/containers/jersey_container_servlet_core"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/ws/rs:javax_ws_rs_api",
        "//3rdparty/jvm/org/glassfish/hk2/external:javax_inject",
        "//3rdparty/jvm/org/glassfish/jersey/core:jersey_common",
        "//3rdparty/jvm/org/glassfish/jersey/core:jersey_server"
    ],
    visibility = [
        "//visibility:public"
    ]
)


