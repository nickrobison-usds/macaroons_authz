//
// Created by usds on 2018-11-30.
//

#ifndef CMSAUTHZCLI_MACAROON_HPP
#define CMSAUTHZCLI_MACAROON_HPP


#include <string>
#include <macaroons.h>

struct MacaroonCaveat {
    std::string location;
    std::string identifier;

    MacaroonCaveat(std::string loc, std::string id): location(loc), identifier(id) {
//        Not used
    }
};

class Macaroon {

private:
    const struct macaroon* m;

    explicit Macaroon(const macaroon *m);


public:
    void inspect();
    std::vector<MacaroonCaveat> get_third_party_caveats();
    const std::string location();
    const macaroon* M();

    const static Macaroon importMacaroons(const std::string &string);
};


#endif //CMSAUTHZCLI_MACAROON_HPP
