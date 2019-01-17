//
// Created by usds on 2018-12-16.
//

#ifndef CMSAUTHZCLI_CLIENT_HPP
#define CMSAUTHZCLI_CLIENT_HPP

#include "bakery/Macaroon.hpp"
#include <cpprest/http_client.h>
#include "interceptor.hpp"

using namespace web::http;                  // Common HTTP functionality

class Client {

private:
    pplx::task<std::vector<Macaroon>> dischargeCaveat(const MacaroonCaveat &cav) const;

    std::vector<const Interceptor*> interceptors;

public:
    Client() = default;

    void addInterceptor(const std::string &location, const Interceptor *interceptor);
    const std::string dischargeMacaroon(Macaroon m, macaroon_format format = MACAROON_V2J) const;
};

#endif //CMSAUTHZCLI_CLIENT_HPP
