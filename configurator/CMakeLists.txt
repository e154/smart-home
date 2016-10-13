cmake_minimum_required(VERSION 3.2)

set(GOPATH "${CMAKE_CURRENT_BINARY_DIR}/vendor")

project(configurator)

SET(SOURCE_LIB

        src/main.go
        )


include(common.cmake)

add_subdirectory(src)
