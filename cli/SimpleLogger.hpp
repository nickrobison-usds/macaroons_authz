//
// Created by Nicholas Robison on 2019-01-18.
//

#ifndef CMSAUTHZCLI_SIMPLELOGGER_HPP
#define CMSAUTHZCLI_SIMPLELOGGER_HPP

#include "extern/rang/include/rang.hpp"
#include <fmt/format.h>

using namespace rang;

/**
 * SimpleLogger wraps stdout and stderror and applies some pretty colors to the output.
 * It uses fmt to handle the formatting,
 */
class SimpleLogger {
public:

    explicit SimpleLogger(bool debugMode = false): enable_debug(debugMode) {}

    void info(const std::string &value) {
        write_string(value, fg::green);
    }
    template<typename... Args>
    void info(const std::string &value, const Args &... args) const {
        write_string(value, fg::green, args...);
    }

    void debug(const std::string &value) const {
        if (enable_debug)
            write_string(value, fg::cyan);
    }

    template<typename... Args>
    void debug(const std::string &value, const Args &... args) const {
        if (enable_debug)
            write_string(value, fg::cyan, args...);
    }

    void trace(const std::string &value) const {
        if (enable_debug)
            write_string(value, fg::blue);
    }

    template<typename... Args>
    void trace(const std::string &value, const Args &... args) const {
        if (enable_debug)
            write_string(value, fg::blue, args...);
    }

    void warn(const std::string &value) const {
        write_string(value, fg::magenta);
    }

    template<typename... Args>
    void warn(const std::string &value, const Args &... args) const {
        write_string(value, fg::magenta, args...);
    }

    void error(const std::string &value) const {
        write_string(std::cerr, value, fg::red);
    }

    template<typename... Args>
    void error(const std::string &value, const Args &... args) const {
        write_string(std::cerr, value, fg::red, args...);
    }

private:

    bool enable_debug;

    void write_string(const std::string &value, const fg color) const {
        write_string(std::cout, value, color);
    }

    template<typename... Args>
    void write_string(const std::string &value, fg color, const Args &... args) const {
        write_string(std::cout, value, color, args...);
    }

    template<typename... Args>
    void write_string(const std::ostream& stream, const std::string &value, fg color, const Args&... args) const {
        std::cout << color << fmt::format(value, args...) << fg::reset << std::endl;
    }

    void write_string(const std::ostream& stream, const std::string &value, const fg color) const {
        std::cout << color << value << fg::reset << std::endl;
    }
};

#endif //CMSAUTHZCLI_SIMPLELOGGER_HPP
