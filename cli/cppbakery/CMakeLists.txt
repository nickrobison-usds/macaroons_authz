cmake_minimum_required(VERSION 3.12)
project(CPPBakery
        VERSION 0.1
        DESCRIPTION "C++ port of go-bakery library"
        LANGUAGES CXX)

# Handle the submodules
find_package(Git QUIET)
if (GIT_FOUND AND EXISTS "${PROJECT_SOURCE_DIR}/../../.git")
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

if (NOT EXISTS "${PROJECT_SOURCE_DIR}/extern/cppcodec/CMakeLists.txt")
    message(FATAL_ERROR "The submodules were not downloaded! GIT_SUBMODULE was turned off or failed. Please update submodules and try again.")
endif ()

add_subdirectory(bakery)
add_subdirectory(httpbakery)