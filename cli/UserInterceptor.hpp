//
// Created by Nicholas Robison on 2019-01-18.
//

#ifndef CMSAUTHZCLI_USER_INTERCEPTOR_HPP
#define CMSAUTHZCLI_USER_INTERCEPTOR_HPP

#include <httpbakery/interceptor.hpp>

// Forward declare the http_request, so we don't have to pull in the entire header.
namespace cppreset {
    namespace http_client {
        class http_request;
    }
}

class UserInterceptor : public Interceptor {

public:
    explicit UserInterceptor(const std::string &userID): userID(userID) {};

    http_request intercept(http_request &request, const std::string &location) override;

private:
    const std::string &userID;
};

#endif //CMSAUTHZCLI_USER_INTERCEPTOR_HPP
