cmake_minimum_required(VERSION 3.12)
project(LoginGovAuth
        VERSION 0.1
        DESCRIPTION "QT based library for interacting with Login.gov"
        LANGUAGES CXX)

find_package(cpprestsdk REQUIRED)
# Add QT
find_package(Qt5 5.12 COMPONENTS Network NetworkAuth Gui REQUIRED)
set(CMAKE_AUTOMOC ON)

add_library(logingov LoginInterceptor.cpp)
target_link_libraries(logingov
        PUBLIC
        Qt5::NetworkAuth
        Qt5::Network
        cpprestsdk::cpprest
        PRIVATE
        Bakery::httpbakery
        Boost::system
        Boost::thread
        Qt5::Gui)
target_include_directories(logingov PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/..)