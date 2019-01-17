//
// Created by Nicholas Robison on 2019-01-17.
//
/**
 * An interceptor modifies the base http_request by adding necessary parameters and paths.
 */
#ifndef CMSAUTHZCLI_INTERCEPTOR_HPP
#define CMSAUTHZCLI_INTERCEPTOR_HPP

#include <cpprest/http_client.h>

using namespace web::http;                  // Common HTTP functionality

struct Interceptor {
    virtual ~Interceptor() = default;
    virtual http_request intercept(http_request &request, const std::string &location) const = 0;
};

// Leaving this here to remind me how to do it.
//template<typename T>
//struct has_intercept_method {
//
//private:
//    template<typename U>
//    static auto test(int) -> decltype(std::declval<U>().intercept(std::declval<int>(), std::declval<std::string>()), std::true_type());
//
//
//
//    template<typename>
//    static std::false_type test(...);
//
//public:
//    static constexpr bool value = std::is_same<decltype(test<T>(0)), std::true_type>::value;
//};



#endif //CMSAUTHZCLI_INTERCEPTOR_HPP
