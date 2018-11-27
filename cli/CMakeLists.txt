cmake_minimum_required(VERSION 3.12)
project(CMSAuthZCLI
        VERSION 0.1
        DESCRIPTION "CLI client for CMS AuthZ demo"
        LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)

# Handle the submodules
find_package(Git QUIET)
if(GIT_FOUND AND EXISTS "${PROJECT_SOURCE_DIR}/../.git")
    option(GIT_SUBMODULE "Checkout submodules during build." ON)
    if(GIT_SUBMODULE)
        message(STATUS "Initializing submodules")
        execute_process(COMMAND ${GIT_EXECUTABLE} submodule update --init --recursive
                WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
                RESULT_VARIABLE GIT_SUBMOD_RESULT)
        if(NOT GIT_SUBMOD_RESULT EQUAL "0")
            message(FATAL_ERROR "git submodule update --init failed with ${GIT_SUBMOD_RESULT}, please checkout submodules")
        endif()
    endif()
endif()

if(NOT EXISTS "${PROJECT_SOURCE_DIR}/extern/cliutils/CMakeLists.txt")
    message(FATAL_ERROR "The submodules were not downloaded! GIT_SUBMODULE was turned off or failed. Please update submodules and try again.")
endif()

# Add the dependencies
#if (NOT TARGET CLI11::CLI11)
add_subdirectory(extern/cliutils)
#endif()

# Should I do this?
#file(GLOB CLI11_headers ${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include/CLI/*)

add_executable(cli main.cpp)
target_link_libraries(cli INTERFACE CLI11::CLI11)
#target_include_directories(cli PUBLIC ${CLI11_headers})
#target_include_directories(cli PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include)
#target_include_directories(cli PRIVATE
#        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include/>
#        $<INSTALL_INTERFACE:include/CLI/>)
target_include_directories(cli PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include)
