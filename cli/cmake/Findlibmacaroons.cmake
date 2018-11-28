# Try to find libmacaroons

find_package(PkgConfig)

message("Finding pkg")

pkg_check_modules(PC_LIBMACAROONS QUIET libmacaroons)
set(LIBMACAROONS_DEFINITIONS ${PC_LIBMACAROONS_CFLAGS_OTHER})

message("Finding path")

find_path(LIBMACAROONS_INCLUDE_DIR
        NAMES macaroons.h
        HINTS ${PC_LIBMACAROONS_INCLUDE_DIRS} ${PC_LIBMACAROONS_INCLUDEDIR})

find_library(LIBMACAROONS_LIBRARY NAMES macaroons
        HINTS ${PC_LIBMACAROONS_LIBDIR} ${PC_LIBMACAROONS_LIBRARY_DIRS})

include(FindPackageHandleStandardArgs)
find_package_handle_standard_args(libmacaroons DEFAULT_MSG
        LIBMACAROONS_INCLUDE_DIR LIBMACAROONS_LIBRARY)

mark_as_advanced(LIBMACAROONS_INCLUDE_DIR LIBMACAROONS_LIBRARY)

if (LIBMACAROONS_FOUND AND NOT TARGET libmacaroons::libmacaroons)
    add_library(libmacaroons::libmacaroons UNKNOWN IMPORTED)
    set_target_properties(libmacaroons::libmacaroons PROPERTIES
            IMPORTED_LOCATION "${LIBMACAROONS_LIBRARY}"
            INTERFACE_INCLUDE_DIRS "${LIBMACAROONS_INCLUDE_DIR}")
endif ()


set(LIBMACAROONS_LIBRARIES ${LIBMACAROONS_LIBRARY})
set(LIBMACAROONS_INCLUDE_DIRS ${LIBMACAROONS_INCLUDE_DIR})

