cmake_minimum_required(VERSION 3.12)
project(LoginGovAuth
        VERSION 0.1
        DESCRIPTION "QT based library for interacting with Login.gov"
        LANGUAGES CXX)

find_package(cpprestsdk REQUIRED)
# Add QT
find_package(Qt5 5.12 COMPONENTS Network NetworkAuth Gui REQUIRED)
set(CMAKE_AUTOMOC ON)

add_library(liblogingov LoginComponent.cpp SharedMemory.hpp)

add_executable(logingov LoginGuiHandler.cpp)
target_link_libraries(liblogingov
        PUBLIC
        Qt5::NetworkAuth
        Qt5::Network
        cpprestsdk::cpprest
        Qt5::Gui
        PRIVATE
        Bakery::httpbakery
        Boost::system
        Boost::thread)
target_link_libraries(logingov liblogingov)
target_include_directories(liblogingov PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/..)