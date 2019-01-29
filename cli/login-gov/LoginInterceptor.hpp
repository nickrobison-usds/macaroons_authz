//
// Created by Nicholas Robison on 2019-01-28.
//

#ifndef CMSAUTHZCLI_LOGINHANDLER_HPP
#define CMSAUTHZCLI_LOGINHANDLER_HPP


#include <QtCore>
#include <httpbakery/interceptor.hpp>

class LoginInterceptor : public Interceptor {

public:
    LoginInterceptor(std::string root_location);
    http_request intercept(http_request &request, const std::string &location) override;

private:
    const QString root_location;
};


#endif //CMSAUTHZCLI_LOGINHANDLER_HPP
