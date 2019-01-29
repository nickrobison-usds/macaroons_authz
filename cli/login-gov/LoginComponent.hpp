//
// Created by Nicholas Robison on 2019-01-28.
//

#ifndef CMSAUTHZCLI_LOGINCOMPONENT_HPP
#define CMSAUTHZCLI_LOGINCOMPONENT_HPP


#include <QObject>
#include <QOAuth2AuthorizationCodeFlow>
#include <cpprest/http_client.h>

using namespace web::http;          // HTTP client features

class LoginComponent : public QObject {
    Q_OBJECT

public:
    explicit LoginComponent(QObject *parent = nullptr);

    void login();
    QString getToken() const;

signals:
    void token(QString token);

private:
    QString m_token;
    QOAuth2AuthorizationCodeFlow oauth;
    QString m_code;
    const QByteArray code_verifier;
};


#endif //CMSAUTHZCLI_LOGINCOMPONENT_HPP
