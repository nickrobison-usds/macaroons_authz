//
// Created by Nicholas Robison on 2019-01-18.
//

#ifndef CMSAUTHZCLI_SIMPLELOGGER_HPP
#define CMSAUTHZCLI_SIMPLELOGGER_HPP


class SimpleLogger {
public:
    void info(const std::string &value) const {
        write_string(value);
    }

    void debug(const std::string &value) const {
        write_string(value);
    }

    void warn(const std::string &value) const {
        write_string(value);
    }

    void error(const std::string &value) const {
        write_string(value);
    }

    void trace(const std::string &value) const {
        write_string(value);
    }

private:
    void write_string(const std::string &value) const {
        std::cout << value << std::endl;
    }
};

#endif //CMSAUTHZCLI_SIMPLELOGGER_HPP
