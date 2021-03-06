# We need this, for now
find_package(cpprestsdk REQUIRED)
# We need to manually link Boost::System and thread. For some reason it's not being pulled in via REST SDK
find_package(Boost REQUIRED COMPONENTS system thread)

# We have to manually handle OpenSSL, for some reason.
find_package(OpenSSL REQUIRED)

add_subdirectory(extern/libmacaroons)

if (NOT TARGET fmt)
    add_subdirectory(${CMAKE_CURRENT_SOURCE_DIR}/../../extern/fmt ${CMAKE_CURRENT_BINARY_DIR}/fmt)
endif()

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
        Boost::thread
        fmt)
target_include_directories(bakery
        INTERFACE
        $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/..>
        $<INSTALL_INTERFACE:${CMAKE_INSTALL_INCLUDEDIR}>)

add_dependencies(bakery cppcodec)

add_library(Bakery::bakery ALIAS bakery)
