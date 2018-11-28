cmake_minimum_required(VERSION 3.12)
project(CMSAuthZCLI
        VERSION 0.1
        DESCRIPTION "CLI client for CMS AuthZ demo"
        LANGUAGES C CXX)

set(CMAKE_CXX_STANDARD 17)

include(cmake/Findlibmacaroons.cmake)

find_package(cpprestsdk REQUIRED)
#find_package(libmacaroons REQUIRED)
# We need to manually link Boost::System and thread. For some reason it's not being pulled in via REST SDK
find_package(Boost REQUIRED COMPONENTS system thread)

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

if(NOT EXISTS "${PROJECT_SOURCE_DIR}/extern/termcolor/CMakeLists.txt")
    message(FATAL_ERROR "The submodules were not downloaded! GIT_SUBMODULE was turned off or failed. Please update submodules and try again.")
endif()

add_subdirectory(extern/cli11)
add_subdirectory(extern/termcolor)
add_subdirectory(extern/cppcodec)
#add_subdirectory(extern/cliutils)

#set(CLI_INCLUDE_DIR ${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include CACHE PATH "CLIUtils include dir")
#
#install(DIRECTORY ${CLI_INCLUDE_DIR}/CLI DESTINATION install)

# Add the dependencies
#if (NOT TARGET CLI11::CLI11)
#add_subdirectory(extern/cliutils)
#endif()

# Should I do this?
#file(GLOB CLI11_headers ${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include/CLI/*)

add_executable(cli main.cpp)
target_link_libraries(cli PRIVATE
        cli11
        termcolor::termcolor
        cpprestsdk::cpprest
        Boost::system
        Boost::thread
        libmacaroons::libmacaroons
        cppcodec)
target_include_directories(cli PRIVATE ${LIBMACAROONS_INCLUDE_DIRS})
message("${LIBMACAROONS_INCLUDE_DIRS}")
set(CMAKE_INCLUDE_PATH /usr/local/include)
#target_include_directories(cli PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/extern/cli11/include)
#target_include_directories(cli PUBLIC ${CLI11_headers})
#target_include_directories(cli PRIVATE ${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include)
#target_include_directories(cli PUBLIC
#        $<BUILD_INTERFACE:${CLI_INCLUDE_DIR}>
#        $<INSTALL_INTERFACE:install>)
#target_include_directories(cli PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/extern/cliutils/include)
