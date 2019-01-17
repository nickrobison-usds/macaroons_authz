# We need this, for now
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

if (NOT EXISTS "${PROJECT_SOURCE_DIR}/extern/cppcodec/CMakeLists.txt")
    message(FATAL_ERROR "The submodules were not downloaded! GIT_SUBMODULE was turned off or failed. Please update submodules and try again.")
endif ()

add_subdirectory(extern/libmacaroons)

add_library(bakery
        Macaroon.cpp)
target_compile_features(bakery
        PRIVATE
        cxx_lambdas
        cxx_auto_type)
target_link_libraries(bakery
        PUBLIC
        macaroons
        cppcodec
        PRIVATE
        cpprestsdk::cpprest
        Boost::system
        Boost::thread)
target_include_directories(bakery
        INTERFACE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/..>
        $<INSTALL_INTERFACE:${CMAKE_INSTALL_INCLUDEDIR}>)

add_dependencies(bakery cppcodec)

add_library(Bakery::bakery ALIAS bakery)