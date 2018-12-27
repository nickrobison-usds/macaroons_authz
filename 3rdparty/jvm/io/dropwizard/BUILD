java_library(
    name = "dropwizard_client",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_client"
    ],
    runtime_deps = [
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_httpclient",
        "//3rdparty/jvm/org/apache/httpcomponents:httpclient",
        "//3rdparty/jvm/org/glassfish/jersey/connectors:jersey_apache_connector",
        "//3rdparty/jvm/org/glassfish/jersey/core:jersey_client",
        "//3rdparty/jvm/org/glassfish/jersey/ext/rx:jersey_rx_client",
        ":dropwizard_core"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_configuration",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_configuration"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/fasterxml/jackson/dataformat:jackson_dataformat_yaml",
        "//3rdparty/jvm/org/apache/commons:commons_lang3",
        "//3rdparty/jvm/org/apache/commons:commons_text",
        ":dropwizard_jackson",
        ":dropwizard_validation"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_core",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_core"
    ],
    runtime_deps = [
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_core",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_healthchecks",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_jmx",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_jvm",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_servlets",
        "//3rdparty/jvm/net/sourceforge/argparse4j:argparse4j",
        "//3rdparty/jvm/org/eclipse/jetty/toolchain/setuid:jetty_setuid_java",
        ":dropwizard_configuration",
        ":dropwizard_jackson",
        ":dropwizard_jersey",
        ":dropwizard_jetty",
        ":dropwizard_lifecycle",
        ":dropwizard_logging",
        ":dropwizard_metrics",
        ":dropwizard_request_logging",
        ":dropwizard_servlets",
        ":dropwizard_util",
        ":dropwizard_validation"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_jackson",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_jackson"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_annotations",
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_core",
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_databind",
        "//3rdparty/jvm/com/fasterxml/jackson/datatype:jackson_datatype_guava",
        "//3rdparty/jvm/com/fasterxml/jackson/datatype:jackson_datatype_jdk8",
        "//3rdparty/jvm/com/fasterxml/jackson/datatype:jackson_datatype_joda",
        "//3rdparty/jvm/com/fasterxml/jackson/datatype:jackson_datatype_jsr310",
        "//3rdparty/jvm/com/fasterxml/jackson/module:jackson_module_afterburner",
        "//3rdparty/jvm/com/fasterxml/jackson/module:jackson_module_parameter_names",
        "//3rdparty/jvm/com/google/guava:guava",
        "//3rdparty/jvm/org/slf4j:slf4j_api",
        ":dropwizard_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_jersey",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_jersey"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/fasterxml/jackson/jaxrs:jackson_jaxrs_json_provider",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_jersey2",
        "//3rdparty/jvm/org/apache/commons:commons_lang3",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_continuation",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_server",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_webapp",
        "//3rdparty/jvm/org/glassfish/jersey/containers:jersey_container_servlet",
        "//3rdparty/jvm/org/glassfish/jersey/core:jersey_server",
        "//3rdparty/jvm/org/glassfish/jersey/ext:jersey_bean_validation",
        "//3rdparty/jvm/org/glassfish/jersey/ext:jersey_metainf_services",
        ":dropwizard_jackson",
        ":dropwizard_logging",
        ":dropwizard_validation"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_jetty",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_jetty"
    ],
    runtime_deps = [
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_jetty9",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_http",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_server",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_servlet",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_servlets",
        ":dropwizard_logging"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_lifecycle",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_lifecycle"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/google/guava:guava",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_server",
        "//3rdparty/jvm/org/slf4j:slf4j_api",
        ":dropwizard_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_logging",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_logging"
    ],
    runtime_deps = [
        "//3rdparty/jvm/ch/qos/logback:logback_classic",
        "//3rdparty/jvm/ch/qos/logback:logback_core",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_logback",
        "//3rdparty/jvm/org/eclipse/jetty:jetty_util",
        "//3rdparty/jvm/org/slf4j:jcl_over_slf4j",
        "//3rdparty/jvm/org/slf4j:jul_to_slf4j",
        "//3rdparty/jvm/org/slf4j:log4j_over_slf4j",
        "//3rdparty/jvm/org/slf4j:slf4j_api",
        ":dropwizard_jackson",
        ":dropwizard_validation"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_metrics",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_metrics"
    ],
    runtime_deps = [
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_core",
        "//3rdparty/jvm/org/slf4j:slf4j_api",
        ":dropwizard_jackson",
        ":dropwizard_lifecycle",
        ":dropwizard_validation"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_request_logging",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_request_logging"
    ],
    runtime_deps = [
        "//3rdparty/jvm/ch/qos/logback:logback_access",
        ":dropwizard_jetty",
        ":dropwizard_logging"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_servlets",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_servlets"
    ],
    runtime_deps = [
        "//3rdparty/jvm/ch/qos/logback:logback_classic",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_annotation",
        "//3rdparty/jvm/io/dropwizard/metrics:metrics_core",
        "//3rdparty/jvm/org/slf4j:slf4j_api",
        ":dropwizard_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_util",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_util"
    ],
    runtime_deps = [
        "//3rdparty/jvm/com/fasterxml/jackson/core:jackson_annotations",
        "//3rdparty/jvm/com/google/code/findbugs:jsr305",
        "//3rdparty/jvm/com/google/guava:guava",
        "//3rdparty/jvm/joda_time:joda_time"
    ],
    visibility = [
        "//visibility:public"
    ]
)



java_library(
    name = "dropwizard_validation",
    exports = [
        "//external:jar/io/dropwizard/dropwizard_validation"
    ],
    runtime_deps = [
        "//3rdparty/jvm/org/glassfish:javax_el",
        "//3rdparty/jvm/org/hibernate:hibernate_validator",
        "//3rdparty/jvm/org/javassist:javassist",
        "//3rdparty/jvm/org/slf4j:slf4j_api",
        ":dropwizard_util"
    ],
    visibility = [
        "//visibility:public"
    ]
)


