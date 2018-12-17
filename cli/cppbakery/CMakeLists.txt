cmake_minimum_required(VERSION 3.12)
project(CPPBakery
        VERSION 0.1
        DESCRIPTION "C++ port of go-bakery library"
        LANGUAGES C CXX)

add_subdirectory(bakery)
add_subdirectory(httpbakery)
add_subdirectory(extern/cppcodec)