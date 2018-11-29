# Try to find libmacaroons

find_package(PkgConfig)

pkg_check_modules(PC_LIBMACAROONS QUIET libmacaroons)
set(LIBMACAROONS_DEFINITIONS ${PC_LIBMACAROONS_CFLAGS_OTHER})

find_path(libmacaroons_INCLUDE_DIR
        NAMES macaroons.h
        HINTS ${PC_LIBMACAROONS_INCLUDE_DIRS} ${PC_LIBMACAROONS_INCLUDEDIR})

find_library(libmacaroons_LIBRARY NAMES macaroons
        HINTS ${PC_LIBMACAROONS_LIBDIR} ${PC_LIBMACAROONS_LIBRARY_DIRS})


include(FindPackageHandleStandardArgs)
find_package_handle_standard_args(libmacaroons
        REQUIRED_VARS libmacaroons_INCLUDE_DIR libmacaroons_LIBRARY
        VERSION_VAR PC_LIBMACAROONS_VERSION)

# Set everything up
if (libmacaroons_FOUND AND NOT TARGET libmacaroons)
    add_library(libmacaroons UNKNOWN IMPORTED)
    set_target_properties(libmacaroons PROPERTIES
            IMPORTED_LOCATION "${libmacaroons_LIBRARY}"
            INTERFACE_INCLUDE_DIRS "${libmacaroons_INCLUDE_DIR}")
endif ()