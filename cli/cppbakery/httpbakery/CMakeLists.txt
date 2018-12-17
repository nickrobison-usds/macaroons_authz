find_package(cpprestsdk REQUIRED)
# We need to manually link Boost::System and thread. For some reason it's not being pulled in via REST SDK
find_package(Boost REQUIRED COMPONENTS system thread)

# We have to manually handle OpenSSL, for some reason.
find_package(OpenSSL REQUIRED)

add_library(httpbakery include/httpbakery/client.hpp client.cpp)
target_link_libraries(httpbakery
        PUBLIC
        Bakery::bakery
        PRIVATE
        cpprestsdk::cpprest
        Boost::system
        Boost::thread
        cppcodec)
target_include_directories(httpbakery
        INTERFACE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
        $<INSTALL_INTERFACE:${CMAKE_INSTALL_INCLUDEDIR}>)

add_library(Bakery::httpbakery ALIAS httpbakery)