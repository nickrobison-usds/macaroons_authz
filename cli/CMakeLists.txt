cmake_minimum_required(VERSION 3.12)

# Do this before project() to skip the compiler config process
if ("${CMAKE_BINARY_DIR}" STREQUAL "${CMAKE_SOURCE_DIR}") 
  message(FATAL_ERROR "In-source builds are disabled.
    Please create a subfolder and use `cmake ..` inside it.
    NOTE: cmake will now create CMakeCache.txt and CMakeFiles/*.
          You must delete them, or cmake will refuse to work.")
endif() # yes, end-markers and even else() need empty parens

project(CMSAuthZCLI
        VERSION 0.1
        DESCRIPTION "CLI client for CMS AuthZ demo"
        LANGUAGES C CXX)

set(CMAKE_CXX_STANDARD 17)

find_package(cpprestsdk REQUIRED)
# We need to manually link Boost::System and thread. For some reason it's not being pulled in via REST SDK
find_package(Boost REQUIRED COMPONENTS system thread)

# We have to manually handle OpenSSL, for some reason.
find_package(OpenSSL REQUIRED)

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

if (NOT EXISTS "${PROJECT_SOURCE_DIR}/extern/spdlog/CMakeLists.txt")
    message(FATAL_ERROR "The submodules were not downloaded! GIT_SUBMODULE was turned off or failed. Please update submodules and try again.")
endif ()

add_subdirectory(extern/cli11)

if (NOT TARGET fmt)
    add_subdirectory(extern/fmt)
endif()

add_subdirectory(extern/spdlog)
add_subdirectory(cppbakery)

add_executable(cli main.cpp SimpleLogger.hpp)
target_link_libraries(cli PRIVATE
        cli11
        Bakery::httpbakery
        cpprestsdk::cpprest
        Boost::system
        Boost::thread
        fmt
        spdlog::spdlog)
