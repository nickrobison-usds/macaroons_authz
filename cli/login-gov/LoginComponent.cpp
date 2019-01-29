//
// Created by Nicholas Robison on 2019-01-28.
//

#include <QByteArray>
#include <QOAuthHttpServerReplyHandler>
#include <QDesktopServices>
#include <QRandomGenerator>
#include <QDebug>
#include <QEventLoop>
#include <QCryptographicHash>
#include "LoginComponent.hpp"


LoginComponent::LoginComponent(QObject *parent) : QObject(parent), m_token(""), code_verifier("test-state-minimum-of-32-which-is-very-very-long") {

    auto hash = QCryptographicHash(QCryptographicHash::Sha256);
    hash.addData(code_verifier);
    const auto code_challenge = hash.result().toBase64(QByteArray::Base64UrlEncoding).toStdString();

    auto replyHandler = new QOAuthHttpServerReplyHandler(1337, this);
    oauth.setScope("openid email");
    oauth.setAuthorizationUrl(QUrl("http://localhost:3000/openid_connect/authorize"));
    oauth.setAccessTokenUrl(QUrl("http://localhost:3000/api/openid_connect/token"));
    oauth.setProperty("acr_values", QString("http://idmanagement.gov/ns/assurance/loa/1"));
    oauth.setClientIdentifier(QString("urn:gov:gsa:openidconnect.profiles:sp:sso:mad:mac_dev"));
    oauth.setReplyHandler(replyHandler);
    oauth.setState(QString("test-state-minimum-of-32-which-is-very-very-long"));

    // This sets the custom parameters for Login.gov
    oauth.setModifyParametersFunction([this, code_challenge](QAbstractOAuth::Stage state, QVariantMap *map) {
        if (state == QAbstractOAuth::Stage::RequestingAuthorization) {
            map->insert("nonce", QString("test-state-minimum-of-32-which-is-very-very-long"));
            map->insert("acr_values", QString("http://idmanagement.gov/ns/assurance/loa/1"));
            map->insert("code_challenge", QString::fromStdString(code_challenge));
            map->insert("code_challenge_method", "S256");

            qDebug() << "printing query params";
            qDebug() << map->keys();
            std::for_each(map->begin(), map->end(), [](const QVariant &value) {
                qDebug() << value;
            });
        }
        if (state == QAbstractOAuth::Stage::RequestingAccessToken) {
            qDebug() << "Requesting token";
            map->insert("code", this->m_code);
            map->insert("code_verifier", this->code_verifier);
//            map->insert("client_assertion", "Test assertion");
//            map->insert("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer");

            qDebug() << "keys " << map->keys();
            std::for_each(map->begin(), map->end(), [](const QVariant &val) {
                qDebug() << "things" << val;
//                qDebug() << key << val;
            });
        }
    });
    // Reply handler?
    const auto handler = oauth.replyHandler();
    connect(handler, &QAbstractOAuthReplyHandler::callbackReceived, [handler, this](const QVariantMap &values) {
        if (!values.empty()) {
            qDebug() << "called back for values" << values.keys();
            auto const val = values.value("code", "");
            qDebug() << "code " << val;
            this->m_code = val.toString();
        }
    });
    // Setup the handlers
    connect(&oauth, &QOAuth2AuthorizationCodeFlow::authorizeWithBrowser, &QDesktopServices::openUrl);
}

void LoginComponent::login() {
    oauth.grant();
    connect(&oauth, &QOAuth2AuthorizationCodeFlow::statusChanged, [this](QAbstractOAuth::Status status) {
        if (status == QAbstractOAuth::Status::Granted) {
            qDebug() << "granted";
            qDebug() << this->oauth.token();
            this->m_token = this->oauth.token();
            emit token(this->oauth.token());
        }
    });
}

QString LoginComponent::getToken() const {
    return m_token;
}
