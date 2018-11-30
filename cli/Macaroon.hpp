//
// Created by usds on 2018-11-30.
//

#ifndef CMSAUTHZCLI_MACAROON_HPP
#define CMSAUTHZCLI_MACAROON_HPP


#include <string>
#include <macaroons.h>

class Macaroon {

private:
    const struct macaroon* m;
    Macaroon(const macaroon *m);


public:
    void inspect();
    const std::string location();
    const macaroon* M();

    const static Macaroon importMacaroons(const std::string &string);
};


#endif //CMSAUTHZCLI_MACAROON_HPP
