#include <utility>

//
// Created by Nicholas Robison on 2019-01-29.
//

#ifndef CMSAUTHZCLI_LOG_HPP
#define CMSAUTHZCLI_LOG_HPP

static const char *const memory_name = "Macaroons";

#include <boost/interprocess/managed_shared_memory.hpp>
#include <boost/interprocess/sync/interprocess_semaphore.hpp>
#include <httpbakery/interceptor.hpp>
#include <login-gov/SharedMemory.hpp>

namespace Shared {
    namespace bip = boost::interprocess;
    using Segment = bip::managed_shared_memory;
    using Manager = Segment::segment_manager;
    template<typename T>
    using Alloc = bip::allocator<T, Manager>;
}

class LoginInterceptor : public Interceptor {
public:
    explicit LoginInterceptor(std::string location) : m_location(std::move(location)) {


    }

    http_request intercept(http_request &request, const std::string &location) override {

        if (location.rfind("http://localhost:5000") == 0) {

            //    Connect to shared memory
            Shared::bip::managed_shared_memory memory(Shared::bip::open_only, memory_name);

            using A = Shared::Alloc<char>;
            A alloc(memory.get_segment_manager());

            auto *data = memory.find_or_construct<SharedMemory<A>>("data")(1024, memory.get_segment_manager());

            std::cout << "Connected to Macaroons shared" << std::endl;
            std::cout << "Intercepting" << std::endl;
            data->mutex.post();
            std::cout << "Waiting" << std::endl;
            data->done.wait();
            const auto token_string = std::string(data->mData.data(), data->mData.size());
            std::cout << token_string << std::endl;

            uri_builder builder(request.absolute_uri());
            builder.append_query("token", token_string);
            request.set_request_uri(builder.to_uri());
            
        }


        return request;
    }

private:
    const std::string m_location;
};


#endif //CMSAUTHZCLI_LOG_HPP
