//
// Created by Nicholas Robison on 2019-01-28.
//

#include <QtGui>
#include "LoginInterceptor.hpp"
#include "LoginComponent.hpp"

LoginInterceptor::LoginInterceptor(const std::string root_location) : root_location(QString::fromStdString(root_location)) {

}

http_request LoginInterceptor::intercept(http_request &request, const std::string &location) {

    if (location.rfind("http://localhost:5000") == 0) {
        int i = 0;
        char *argv[] = {};
        QGuiApplication app(i, argv);
        LoginComponent comp;
        comp.login();
        QString token;
        QObject::connect(&comp, &LoginComponent::token, [&comp, &token, this](const QString & token_resp) {
            qDebug() << "has token " << token_resp;
            token = token_resp;
            QGuiApplication::quit();
        });

        QGuiApplication::exec();
        uri_builder builder(request.absolute_uri());
        builder.append_query("token", token.toStdString());
        request.set_request_uri(builder.to_uri());
    }
    return request;
}
