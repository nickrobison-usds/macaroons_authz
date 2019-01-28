//
// Created by Nicholas Robison on 2019-01-28.
//

#ifndef CMSAUTHZCLI_LOGINHANDLER_HPP
#define CMSAUTHZCLI_LOGINHANDLER_HPP


#include <QtCore>
#include <QtNetworkAuth>
#include <httpbakery/interceptor.hpp>

class LoginInterceptor : public QObject, public Interceptor {
    Q_OBJECT

public:
    LoginInterceptor(std::string root_location, QObject *parent = nullptr);
    http_request intercept(http_request &request, const std::string &location) override;

private:
    QOAuth2AuthorizationCodeFlow oauth;
    const QString root_location;
};


#endif //CMSAUTHZCLI_LOGINHANDLER_HPP
