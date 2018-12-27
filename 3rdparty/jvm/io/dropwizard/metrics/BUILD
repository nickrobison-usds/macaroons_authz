java_library(
    name = "metrics_annotation",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_annotation"
    ],
    runtime_deps = [
        "//3rdparty/jvm/org/slf4j:slf4j_api"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_core",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_core"
    ],
    runtime_deps = [
        "//3rdparty/jvm/org/slf4j:slf4j_api"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_healthchecks",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_healthchecks"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_httpclient",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_httpclient"
    ],
    runtime_deps = [
        ":metrics_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_jersey2",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_jersey2"
    ],
    runtime_deps = [
        ":metrics_annotation",
        ":metrics_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_jetty9",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_jetty9"
    ],
    runtime_deps = [
        ":metrics_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_jmx",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_jmx"
    ],
    runtime_deps = [
        ":metrics_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_json",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_json"
    ],
    runtime_deps = [
        ":metrics_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_jvm",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_jvm"
    ],
    runtime_deps = [
        ":metrics_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_logback",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_logback"
    ],
    runtime_deps = [
        ":metrics_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "metrics_servlets",
    exports = [
        "//external:jar/io/dropwizard/metrics/metrics_servlets"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/papertrail:profiler",
        ":metrics_core",
        ":metrics_healthchecks",
        ":metrics_json",
        ":metrics_jvm"
    ],
    visibility = [
        "//visibility:public"
    ]
)


