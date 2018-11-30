cmake_minimum_required(VERSION 3.12)
project(CMSAuthZCLI
        VERSION 0.1
        DESCRIPTION "CLI client for CMS AuthZ demo"
        LANGUAGES C CXX)

set(CMAKE_CXX_STANDARD 17)
# Use this for importing the custom cmake finders
set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH}} "${CMAKE_SOURCE_DIR}/cmake")
set(CMAKE_PREFIX_PATH ${CMAKE_MODULE_PATH} "${CMAKE_BINARY_DIR}/extern")

find_package(cpprestsdk REQUIRED)
find_package(libmacaroons 0.4)
# We need to manually link Boost::System and thread. For some reason it's not being pulled in via REST SDK
find_package(Boost REQUIRED COMPONENTS system thread)

# Handle the submodules
find_package(Git QUIET)
if (GIT_FOUND AND EXISTS "${PROJECT_SOURCE_DIR}/../.git")
    option(GIT_SUBMODULE "Checkout submodules during build." ON)
    if (GIT_SUBMODULE)
        message(STATUS "Initializing submodules")
        execute_process(COMMAND ${GIT_EXECUTABLE} submodule update --init --recursive
                WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
                RESULT_VARIABLE GIT_SUBMOD_RESULT)
        if (NOT GIT_SUBMOD_RESULT EQUAL "0")
            message(FATAL_ERROR "git submodule update --init failed with ${GIT_SUBMOD_RESULT}, please checkout submodules")
        endif ()
    endif ()
endif ()

if (NOT EXISTS "${PROJECT_SOURCE_DIR}/extern/termcolor/CMakeLists.txt")
    message(FATAL_ERROR "The submodules were not downloaded! GIT_SUBMODULE was turned off or failed. Please update submodules and try again.")
endif ()

# Static build libmacaroons, if not found
if (NOT libmacaroons_FOUND)
    message("Did not find system libmacaroons, building static")
    include(ExternalProject)
    ExternalProject_Add(project_libmacaroons
            PREFIX ${CMAKE_CURRENT_BINARY_DIR}/extern
            GIT_REPOSITORY https://github.com/rescrv/libmacaroons.git
            CONFIGURE_COMMAND autoreconf --install > /dev/null 2>&1 || autoreconf && ./configure --prefix=${CMAKE_CURRENT_BINARY_DIR}/extern --enable-shared=false
            BUILD_COMMAND make -j
            BUILD_IN_SOURCE 1
            INSTALL_COMMAND make -j install
            PREFIX=${CMAKE_CURRENT_BINARY_DIR}/extern)

    ExternalProject_Get_Property(project_libmacaroons install_dir)
    add_library(libmacaroons STATIC IMPORTED)
    set_property(TARGET libmacaroons PROPERTY IMPORTED_LOCATION ${install_dir}/lib/libmacaroons.a)
    add_dependencies(libmacaroons project_libmacaroons)
    target_include_directories(libmacaroons INTERFACE ${install_dir}/include)
    add_executable(test_cmake main.cpp)
endif ()

add_subdirectory(extern/cli11)
add_subdirectory(extern/termcolor)
add_subdirectory(extern/cppcodec)

add_executable(cli main.cpp Macaroon.cpp Macaroon.hpp)
target_link_libraries(cli PRIVATE
        cli11
        termcolor::termcolor
        cpprestsdk::cpprest
        Boost::system
        Boost::thread
        libmacaroons
        cppcodec)
message(${libmacaroons_INCLUDE_DIR})
message(${libmacaroons_LIBRARY})
target_include_directories(cli PRIVATE ${libmacaroons_INCLUDE_DIR})