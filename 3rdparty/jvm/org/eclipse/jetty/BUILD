java_library(
    name = "jetty_continuation",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_continuation"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_http",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_http"
    ],
    runtime_deps = [
        ":jetty_io",
        ":jetty_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_io",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_io"
    ],
    runtime_deps = [
        ":jetty_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_security",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_security"
    ],
    runtime_deps = [
        ":jetty_server"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_server",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_server"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/servlet:javax_servlet_api",
        ":jetty_http",
        ":jetty_io"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_servlet",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_servlet"
    ],
    runtime_deps = [
        ":jetty_security"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_servlets",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_servlets"
    ],
    runtime_deps = [
        ":jetty_continuation",
        ":jetty_http",
        ":jetty_io",
        ":jetty_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_util",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_webapp",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_webapp"
    ],
    runtime_deps = [
        ":jetty_servlet",
        ":jetty_xml"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "jetty_xml",
    exports = [
        "//external:jar/org/eclipse/jetty/jetty_xml"
    ],
    runtime_deps = [
        ":jetty_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)


