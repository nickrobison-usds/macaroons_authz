//
// Created by usds on 2018-12-16.
//

#ifndef CMSAUTHZCLI_CLIENT_HPP
#define CMSAUTHZCLI_CLIENT_HPP

#include <bakery/Macaroon.hpp>

class Client {

public:
    Client();

    const std::string dischargeMacaroon(const Macaroon m, const macaroon_format format = MACAROON_V1) const;

private:
    pplx::task<Macaroon> dischargeCaveat(const MacaroonCaveat &cav) const;
};

#endif //CMSAUTHZCLI_CLIENT_HPP
