//
// Created by Nicholas Robison on 2019-01-28.
//

#include <QtGui>
#include <QtNetworkAuth>

#include "LoginInterceptor.hpp"

LoginInterceptor::LoginInterceptor(const std::string root_location, QObject *parent) : QObject(parent), root_location(QString::fromStdString(root_location)) {
    auto replyHandler = new QOAuthHttpServerReplyHandler(1337, this);
    oauth.setScope("openid");
    oauth.setAuthorizationUrl(QUrl("http://localhost:3000/openid_connect/authorized"));
    oauth.setToken("http://localhost:3000/openid_connect/token");
    oauth.setProperty("acr_values", QString("http://idmanagement.gov/ns/assurance/loa/1"));
    oauth.setReplyHandler(replyHandler);
    connect(&oauth, &QOAuth2AuthorizationCodeFlow::authorizeWithBrowser, &QDesktopServices::openUrl);
}

http_request LoginInterceptor::intercept(http_request &request, const std::string &location) {
    // Authorize and grant the token
    oauth.grant();
    const auto token = oauth.token();

    uri_builder builder(request.absolute_uri());
    builder.append_query("token", token.toStdString());
    request.set_request_uri(builder.to_uri());
    return request;
}
