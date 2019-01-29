//
// Created by Nicholas Robison on 2019-01-28.
//

#ifndef CMSAUTHZCLI_CALLBACKHANDLER_HPP
#define CMSAUTHZCLI_CALLBACKHANDLER_HPP


#include <QAbstractOAuthReplyHandler>

class CallbackHandler : public QAbstractOAuthReplyHandler {

public:
    explicit CallbackHandler();
};


#endif //CMSAUTHZCLI_CALLBACKHANDLER_HPP
