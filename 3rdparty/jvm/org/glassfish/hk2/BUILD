java_library(
    name = "hk2_api",
    exports = [
        "//external:jar/org/glassfish/hk2/hk2_api"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/inject:javax_inject",
        "//3rdparty/jvm/org/glassfish/hk2/external:aopalliance_repackaged",
        ":hk2_utils"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "hk2_locator",
    exports = [
        "//external:jar/org/glassfish/hk2/hk2_locator"
    ],
    runtime_deps = [
        "//3rdparty/jvm/org/glassfish/hk2/external:aopalliance_repackaged",
        "//3rdparty/jvm/org/glassfish/hk2/external:javax_inject",
        "//3rdparty/jvm/org/javassist:javassist",
        ":hk2_api",
        ":hk2_utils"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "hk2_utils",
    exports = [
        "//external:jar/org/glassfish/hk2/hk2_utils"
    ],
    runtime_deps = [
        "//3rdparty/jvm/javax/inject:javax_inject"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "osgi_resource_locator",
    exports = [
        "//external:jar/org/glassfish/hk2/osgi_resource_locator"
    ],
    visibility = [
        "//visibility:public"
    ]
)


